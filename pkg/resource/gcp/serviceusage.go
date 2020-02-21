package gcp

import (
	"encoding/json"

	serviceusage "google.golang.org/api/serviceusage/v1"
)

// ServiceAPIResource represents a Google Service API resource
type ServiceAPIResource struct {
	a *serviceusage.GoogleApiServiceusageV1Service
}

// NewServiceAPIResource returns a new ServiceAPIResource
func NewServiceAPIResource(a *serviceusage.GoogleApiServiceusageV1Service) *ServiceAPIResource {
	r := new(ServiceAPIResource)
	r.a = a
	return r
}

// Name returns the bucket's name
func (r *ServiceAPIResource) Name() string {
	return r.a.Name
}

// Marshal returns the underlying resource's JSON representation
func (r *ServiceAPIResource) Marshal() ([]byte, error) {
	return json.Marshal(&r.a)
}
