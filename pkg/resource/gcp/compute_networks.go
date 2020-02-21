package gcp

import (
	"encoding/json"

	compute "google.golang.org/api/compute/v1"
)

// ComputeNetworkResource represents a Google Compute Engine network
type ComputeNetworkResource struct {
	n *compute.Network
}

// NewComputeNetworkResource returns a new ComputeNetworkResource
func NewComputeNetworkResource(n *compute.Network) *ComputeNetworkResource {
	r := new(ComputeNetworkResource)
	r.n = n
	return r
}

// Name returns the name of the Compute network
func (r *ComputeNetworkResource) Name() string {
	return r.n.Name
}

// Marshal returns the underlying resource's JSON representation
func (r *ComputeNetworkResource) Marshal() ([]byte, error) {
	return json.Marshal(&r.n)
}

// IsDefault tests whether the network's name is `default`, which usually comes with a project that
// just enabled it's Compute API
func (r *ComputeNetworkResource) IsDefault() bool {
	return r.n.Name == "default"
}

// IsLegacy tests whether the network is a legacy network
func (r *ComputeNetworkResource) IsLegacy() bool {

	// If IPv4Range is non-empty, then it is a legacy network
	return r.n.IPv4Range != ""
}

// NameEquals tests whether the network's name is equal to what is expected
func (r *ComputeNetworkResource) NameEquals(name string) (result bool, err error) {
	result = r.n.Name == name
	return
}
