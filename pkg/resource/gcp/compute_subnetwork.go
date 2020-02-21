package gcp

import (
	"encoding/json"

	compute "google.golang.org/api/compute/v1"
)

// ComputeSubnetworkResource represents a Google Compute Engine subnetwork
type ComputeSubnetworkResource struct {
	s *compute.Subnetwork
}

// NewComputeSubnetworkResource returns a new ComputeSubnetworkResource
func NewComputeSubnetworkResource(s *compute.Subnetwork) *ComputeSubnetworkResource {
	r := new(ComputeSubnetworkResource)
	r.s = s
	return r
}

// Name returns the name of the Compute subnetwork
func (r *ComputeSubnetworkResource) Name() string {
	return r.s.Name
}

// Region returns the GCP region of the Compute subnetwork
func (r *ComputeSubnetworkResource) Region() string {
	return r.s.Region
}

// Marshal returns the underlying resource's JSON representation
func (r *ComputeSubnetworkResource) Marshal() ([]byte, error) {
	return json.Marshal(&r.s)
}

// IsPrivateGoogleAccessEnabled returns whether private Google network access is enabled
func (r *ComputeSubnetworkResource) IsPrivateGoogleAccessEnabled() bool {
	return r.s.PrivateIpGoogleAccess
}

// IsFlowLogsEnabled returns whether the subnet has VPC flow logs enabled
func (r *ComputeSubnetworkResource) IsFlowLogsEnabled() bool {
	return r.s.EnableFlowLogs
}
