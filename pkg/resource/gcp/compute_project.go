package gcp

import (
	"encoding/json"

	compute "google.golang.org/api/compute/v1"
)

// ComputeProjectResource represents a Google Compute Engine's project information.
type ComputeProjectResource struct {
	p *compute.Project
}

// NewComputeProjectResource returns a new ComputeProjectResource
func NewComputeProjectResource(p *compute.Project) *ComputeProjectResource {
	r := new(ComputeProjectResource)
	r.p = p
	return r
}

// Name is the Project's Name
func (r *ComputeProjectResource) Name() string {
	return r.p.Name
}

// Marshal returns the underlying resource's JSON representation
func (r *ComputeProjectResource) Marshal() ([]byte, error) {
	return json.Marshal(&r.p)
}

// IsXpnHost tests whether the project is configured as a Shared VPC (Xpn) host project
func (r *ComputeProjectResource) IsXpnHost() (result bool, err error) {
	result = r.p.XpnProjectStatus == "HOST"
	return
}
