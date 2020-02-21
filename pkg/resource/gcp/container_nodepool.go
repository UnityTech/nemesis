package gcp

import (
	"encoding/json"
	"errors"
	"fmt"

	container "google.golang.org/api/container/v1"
)

// ContainerNodePoolResource is a resource for testing information about a GKE Node Pool's configuration
type ContainerNodePoolResource struct {
	n *container.NodePool
}

// NewContainerNodePoolResource returns a new ContainerNodePoolResource
func NewContainerNodePoolResource(n *container.NodePool) *ContainerNodePoolResource {
	r := new(ContainerNodePoolResource)
	r.n = n
	return r
}

// Marshal returns the underlying resource's JSON representation
func (r *ContainerNodePoolResource) Marshal() ([]byte, error) {
	return json.Marshal(&r.n)
}

// Name returns the name given to the cluster nodepool
func (r *ContainerNodePoolResource) Name() string {
	return r.n.Name
}

// IsLegacyMetadataAPIDisabled returns whether the given Node Pool has legacy metadata APIs disabled
func (r *ContainerNodePoolResource) IsLegacyMetadataAPIDisabled() (result bool, err error) {
	var val string
	var ok bool
	if val, ok = r.n.Config.Metadata["disable-legacy-endpoints"]; !ok {
		err = errors.New("Could not find key 'disable-legacy-endpoints'")
	}
	if val != "true" {
		err = fmt.Errorf("Invalid value for `disable-legacy-endpoints`, got `%v'", val)
	}
	result = err == nil
	return
}

// IsAutoRepairEnabled returns whether a Node Pool is configured to automatically repair on error
func (r *ContainerNodePoolResource) IsAutoRepairEnabled() bool {
	return r.n.Management.AutoRepair
}

// IsAutoUpgradeEnabled returns whether a Node Pool is configured to automatically upgrade GKE versions
func (r *ContainerNodePoolResource) IsAutoUpgradeEnabled() bool {
	return r.n.Management.AutoUpgrade
}

// CheckDistributionTypeIs returns whether a Node Pool's OS distribution is the expected type
func (r *ContainerNodePoolResource) CheckDistributionTypeIs(expected string) (result bool, err error) {
	result = r.n.Config.ImageType == expected
	if !result {
		err = fmt.Errorf("Node pool is using %v, not %v", r.n.Config.ImageType, expected)
	}
	return
}
