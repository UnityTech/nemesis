package gcp

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	compute "google.golang.org/api/compute/v1"
)

// ComputeInstanceResource represents a Google Compute Engine instance
type ComputeInstanceResource struct {
	i *compute.Instance
}

// NewComputeInstanceResource returns a new ComputeInstanceResource
func NewComputeInstanceResource(i *compute.Instance) *ComputeInstanceResource {
	r := new(ComputeInstanceResource)
	r.i = i
	return r
}

// Name is the compute instance name
func (r *ComputeInstanceResource) Name() string {
	return r.i.Name
}

// Marshal returns the underlying resource's JSON representation
func (r *ComputeInstanceResource) Marshal() ([]byte, error) {
	return json.Marshal(&r.i)
}

// HasNatIP returns whether the instance has an external / NAT ip.
func (r *ComputeInstanceResource) HasNatIP() bool {
	return r.i.NetworkInterfaces[0].AccessConfigs != nil
}

// HasNumNetworkInterfaces returns whether the instance has the expected number of network interfaces
func (r *ComputeInstanceResource) HasNumNetworkInterfaces(num int) (result bool, err error) {
	actual := len(r.i.NetworkInterfaces)
	result = actual == num
	if !result {
		err = fmt.Errorf("Expected %v interfaces, found %v", num, actual)
	}
	return
}

// KeyValueEquals returns whether the metadata key equals a given value.
// Reports an error if they metadata key does not exist.
func (r *ComputeInstanceResource) KeyValueEquals(key string, value string) (result bool, err error) {

	// Loop over project metadata keys until we find the key.
	result = false
	found := false
	for _, item := range r.i.Metadata.Items {
		if item.Key == key {

			// If we found the key, we want to check its value.
			// If it is set correctly, all is well. Otherwise, report the
			found = true
			if strings.ToLower(*item.Value) == strings.ToLower(value) {
				result = true
			} else {
				err = fmt.Errorf("Instance metadata key '%v' is set to '%v'", key, *item.Value)
			}
			return
		}
	}

	// Report an error that the key did not exist
	if !found {
		err = fmt.Errorf("Could not find instance metadata key: %v", key)
	}
	return
}

// KeyAbsent returns whether the metadata key is absent
// Reports an error if they metadata key does not exist.
func (r *ComputeInstanceResource) KeyAbsent(key string) bool {

	// Loop over project metadata keys until we find the key.
	found := false
	for _, item := range r.i.Metadata.Items {
		if item.Key == key {

			// If we found the key, return
			found = true
			break
		}
	}

	return !found
}

// UsesDefaultServiceAccount returns whether the service account used to launch the instance
// is a default compute service account for any project
func (r *ComputeInstanceResource) UsesDefaultServiceAccount() bool {
	return strings.Contains(r.i.ServiceAccounts[0].Email, "-compute@developer.gserviceaccount.com")
}

// HasIPForwardingEnabled returns whether an instance can forward packets for different sources
func (r *ComputeInstanceResource) HasIPForwardingEnabled() bool {
	return r.i.CanIpForward
}

// UsesCustomerSuppliedEncryptionKeys returns whether the instance's disks are encrypted with a CSEK
func (r *ComputeInstanceResource) UsesCustomerSuppliedEncryptionKeys() (err error) {

	var errBuilder strings.Builder

	for _, d := range r.i.Disks {

		if d.DiskEncryptionKey == nil {
			errBuilder.WriteString(fmt.Sprintf("Disk does not use CSEK: %v", d.Source))
		}
	}

	errString := errBuilder.String()
	if errString != "" {
		err = errors.New(errString)
	} else {
		err = nil
	}

	return
}
