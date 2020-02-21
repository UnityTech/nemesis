package gcp

import (
	"encoding/json"

	compute "google.golang.org/api/compute/v1"
)

// ComputeFirewallRuleResource is a resource describing a Google Compute Firewall Rule
type ComputeFirewallRuleResource struct {
	f *compute.Firewall
}

// NewComputeFirewallRuleResource returns a new ComputeFirewallRuleResource
func NewComputeFirewallRuleResource(f *compute.Firewall) *ComputeFirewallRuleResource {
	r := new(ComputeFirewallRuleResource)
	r.f = f
	return r
}

// Marshal returns the underlying resource's JSON representation
func (r *ComputeFirewallRuleResource) Marshal() ([]byte, error) {
	return json.Marshal(&r.f)
}

// Name returns the name of the firewall rule
func (r *ComputeFirewallRuleResource) Name() string {
	return r.f.Name
}

// Network returns the network the firewall rule resides within
func (r *ComputeFirewallRuleResource) Network() string {
	return r.f.Network
}

// AllowsSourceRange returns whether a given CIDR range is allowed by the firewall rule
func (r *ComputeFirewallRuleResource) AllowsSourceRange(sourceRange string) (result bool) {
	for _, s := range r.f.SourceRanges {
		if s == sourceRange {
			result = true
		}
	}
	return
}

// AllowsProtocolPort returns whether a given protocol:port combination is allowed by this firewall rule
func (r *ComputeFirewallRuleResource) AllowsProtocolPort(protocol string, port string) (result bool) {
	for _, allowRule := range r.f.Allowed {
		if allowRule.IPProtocol == protocol {
			for _, p := range allowRule.Ports {
				if port == p {
					result = true
				}
			}
		}
	}
	return
}
