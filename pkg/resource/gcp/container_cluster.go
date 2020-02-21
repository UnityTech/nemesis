package gcp

import (
	"encoding/json"
	"fmt"

	container "google.golang.org/api/container/v1"
)

const (
	loggingService    = "logging.googleapis.com"
	monitoringService = "monitoring.googleapis.com"
)

// ContainerClusterResource is a resource for testing information about a GKE Cluster's configuration
type ContainerClusterResource struct {
	c *container.Cluster
}

// NewContainerClusterResource returns a new ContainerClusterResource
func NewContainerClusterResource(c *container.Cluster) *ContainerClusterResource {
	r := new(ContainerClusterResource)
	r.c = c
	return r
}

// Marshal returns the underlying resource's JSON representation
func (r *ContainerClusterResource) Marshal() ([]byte, error) {
	return json.Marshal(&r.c)
}

// Name returns the name given to the container cluster
func (r *ContainerClusterResource) Name() string {
	return r.c.Name
}

// IsStackdriverLoggingEnabled indicates whether logging.googleapis.com is set as the logging service
func (r *ContainerClusterResource) IsStackdriverLoggingEnabled() bool {
	return r.c.LoggingService == loggingService
}

// IsStackdriverMonitoringEnabled indicates whether monitoring.googleapis.com is set as the monitoring service
func (r *ContainerClusterResource) IsStackdriverMonitoringEnabled() bool {
	return r.c.MonitoringService == monitoringService
}

// IsAliasIPEnabled indicates whether VPC Alias IPs are being used
func (r *ContainerClusterResource) IsAliasIPEnabled() bool {
	return r.c.IpAllocationPolicy.UseIpAliases
}

// IsPodSecurityPolicyControllerEnabled indicates whether PSP controller is enabled
// TODO - currently no way to implement this check by default
func (r *ContainerClusterResource) IsPodSecurityPolicyControllerEnabled() bool {
	return false
	// TODO - implement!
}

// IsDashboardAddonDisabled returns whether the GKE cluster has Kubernetes Dashboard add-on is enabled
func (r *ContainerClusterResource) IsDashboardAddonDisabled() bool {
	return r.c.AddonsConfig.KubernetesDashboard.Disabled
}

// IsMasterAuthorizedNetworksEnabled returns whether the GKE cluster is using master authorized networks
func (r *ContainerClusterResource) IsMasterAuthorizedNetworksEnabled() bool {
	return r.c.MasterAuthorizedNetworksConfig.Enabled
}

// IsAbacDisabled returns whether the GKE cluster is using (legacy) Atributed-Based Access Control
func (r *ContainerClusterResource) IsAbacDisabled() bool {
	if r.c.LegacyAbac == nil {
		return true
	}
	return !r.c.LegacyAbac.Enabled
}

// IsNetworkPolicyAddonEnabled returns whether the GKE cluster has Network Policy add-on enabled
func (r *ContainerClusterResource) IsNetworkPolicyAddonEnabled() bool {
	return !r.c.AddonsConfig.NetworkPolicyConfig.Disabled

}

// IsClientCertificateDisabled checks whether client certificates are disabled
func (r *ContainerClusterResource) IsClientCertificateDisabled() bool {
	if r.c.MasterAuth.ClientCertificateConfig == nil {
		return true
	}
	return !r.c.MasterAuth.ClientCertificateConfig.IssueClientCertificate
}

// IsMasterAuthPasswordDisabled returns whether the GKE cluster has username/password authentication enabled
func (r *ContainerClusterResource) IsMasterAuthPasswordDisabled() bool {
	return r.c.MasterAuth.Password == ""
}

// IsMasterPrivate returns whether the GKE cluster master is only accessible on private networks
func (r *ContainerClusterResource) IsMasterPrivate() bool {
	if r.c.PrivateClusterConfig == nil {
		return false
	}
	return r.c.PrivateClusterConfig.EnablePrivateEndpoint
}

// IsNodesPrivate returns whether the GKE cluster nodes are only accessible on private networks
func (r *ContainerClusterResource) IsNodesPrivate() bool {
	if r.c.PrivateClusterConfig == nil {
		return false
	}
	return r.c.PrivateClusterConfig.EnablePrivateNodes
}

// IsUsingDefaultServiceAccount returns whether the GKE cluster is using the default compute service account
func (r *ContainerClusterResource) IsUsingDefaultServiceAccount() bool {
	return r.c.NodeConfig.ServiceAccount == "default"
}

// IsUsingMinimalOAuthScopes returns whether the GKE cluster is using the defined minimal oauth scopes for a cluster
func (r *ContainerClusterResource) IsUsingMinimalOAuthScopes() (result bool, err error) {

	// Begin with the assumption that we are using minimal oauth scopes
	extraScopes := []string{}

	// Iterate over the cluster's OAuth scopes and determine if they are at most the minimal list provided.
	// If there are any scopes that are not in the whitelist, track them and report them as an error
	for _, scope := range r.c.NodeConfig.OauthScopes {

		found := false
		// Now check if the cluster's scope is in out oauth scopes
		for _, minimalScope := range minimalOAuthScopes {
			if minimalScope == scope {
				found = true
				break
			}
		}

		if !found {
			extraScopes = append(extraScopes, scope)
		}
	}

	result = len(extraScopes) == 0
	if !result {
		err = fmt.Errorf("Cluster is not using minimal scopes. The following scopes are not considered minimal: %v", extraScopes)
	}

	return
}
