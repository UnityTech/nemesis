package gcp

import (
	"encoding/json"

	compute "google.golang.org/api/compute/v1"
)

// ComputeAddressResource is a resource describing a Google Compute Global Address, or Public IPv4 address
type ComputeAddressResource struct {
	a *compute.Address
}

// NewComputeAddressResource returns a new ComputeAddressResource
func NewComputeAddressResource(a *compute.Address) *ComputeAddressResource {
	r := new(ComputeAddressResource)
	r.a = a
	return r
}

// Marshal returns the underlying resource's JSON representation
func (r *ComputeAddressResource) Marshal() ([]byte, error) {
	return json.Marshal(&r.a)
}

// Name returns the name of the firewall rule
func (r *ComputeAddressResource) Name() string {
	return r.a.Name
}

// Network returns the network the firewall rule resides within
func (r *ComputeAddressResource) Network() string {
	return r.a.Network
}
