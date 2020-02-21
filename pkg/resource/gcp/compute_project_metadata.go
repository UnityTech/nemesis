package gcp

import (
	"encoding/json"
	"fmt"
	"strings"

	compute "google.golang.org/api/compute/v1"
)

// ComputeProjectMetadataResource is a resource for testing information about a Project's compute metadata configuration
type ComputeProjectMetadataResource struct {
	m *compute.Metadata
}

// NewComputeProjectMetadataResource returns a new ComputeProjectMetadataResource
func NewComputeProjectMetadataResource(m *compute.Metadata) *ComputeProjectMetadataResource {
	r := new(ComputeProjectMetadataResource)
	r.m = m
	return r
}

// Marshal returns the underlying resource's JSON representation
func (r *ComputeProjectMetadataResource) Marshal() ([]byte, error) {
	return json.Marshal(&r.m)
}

// Includes indicates whether the Metadata object contains the key specified
func (r *ComputeProjectMetadataResource) Includes(key string) (result bool, err error) {

	// Loop over all project metadata keys
	result = false
	for _, item := range r.m.Items {
		if item.Key == key {
			result = true
			break
		}
	}

	return
}

// KeyValueEquals returns whether the metadata key equals a given value.
// Reports an error if they metadata key does not exist.
func (r *ComputeProjectMetadataResource) KeyValueEquals(key string, value string) (result bool, err error) {

	// Loop over project metadata keys until we find the key.
	result = false
	found := false
	for _, item := range r.m.Items {
		if item.Key == key {

			// If we found the key, we want to check its value.
			// If it is set correctly, all is well. Otherwise, report the
			found = true
			if strings.ToLower(*item.Value) == strings.ToLower(value) {
				result = true
			} else {
				err = fmt.Errorf("Project metadata key '%v' is set to '%v'", key, *item.Value)
			}
			return
		}
	}

	// Report an error that the key did not exist
	if !found {
		err = fmt.Errorf("Could not find project metadata key: %v", key)
	}
	return
}

// KeyAbsent returns whether the metadata key is absent
// Reports an error if they metadata key does not exist.
func (r *ComputeProjectMetadataResource) KeyAbsent(key string) bool {

	// Loop over project metadata keys until we find the key.
	found := false
	for _, item := range r.m.Items {
		if item.Key == key {

			// If we found the key, return
			found = true
			break
		}
	}

	return !found
}
