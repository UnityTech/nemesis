package client

import (
	"fmt"

	"github.com/UnityTech/nemesis/pkg/report"
	"github.com/UnityTech/nemesis/pkg/resource/gcp"
	"github.com/UnityTech/nemesis/pkg/utils"
	"github.com/golang/glog"

	compute "google.golang.org/api/compute/v1"
)

// GetNetworkResources launches the process retrieving network resources
func (c *Client) GetNetworkResources() error {

	defer utils.Elapsed("GetNetworkResources")()

	regionNames, err := c.getRegionNames()
	if err != nil {
		glog.Fatalf("%v", err)
	}

	worker := func(projectIDs <-chan string, results chan<- networkCallResult) {
		id := <-projectIDs

		res := networkCallResult{
			ProjectID:   id,
			Networks:    []*gcp.ComputeNetworkResource{},
			Subnetworks: []*gcp.ComputeSubnetworkResource{},
			Firewalls:   []*gcp.ComputeFirewallRuleResource{},
			Addresses:   []*gcp.ComputeAddressResource{},
		}

		// Get all networks active in the project
		networks, err := c.computeClient.Networks.List(id).Do()
		if err != nil {
			glog.Fatalf("Error retrieving networks from project '%v': %v", id, err)
		}

		for _, n := range networks.Items {
			res.Networks = append(res.Networks, gcp.NewComputeNetworkResource(n))
		}

		// Get all subnetworks active in the project
		for _, region := range regionNames {
			var subnetworks *compute.SubnetworkList

			subnetworks, err = c.computeClient.Subnetworks.List(id, region).Do()
			if err != nil {
				glog.Fatalf("Error retrieving subnetworks from project '%v': %v", id, err)
			}

			for _, s := range subnetworks.Items {
				res.Subnetworks = append(res.Subnetworks, gcp.NewComputeSubnetworkResource(s))
			}

			for subnetworks.NextPageToken != "" {
				subnetworks, err := c.computeClient.Subnetworks.List(id, region).PageToken(subnetworks.NextPageToken).Do()
				if err != nil {
					glog.Fatalf("Error retrieving subnetworks from project '%v': %v", id, err)
				}

				for _, s := range subnetworks.Items {
					res.Subnetworks = append(res.Subnetworks, gcp.NewComputeSubnetworkResource(s))
				}
			}
		}

		// Get all firewall rules active in audited projects, for all networks in the projects

		var firewalls *compute.FirewallList

		firewalls, err = c.computeClient.Firewalls.List(id).Do()
		if err != nil {
			glog.Fatalf("Error retrieving firewall rules from project '%v': %v", id, err)
		}

		for _, f := range firewalls.Items {
			res.Firewalls = append(res.Firewalls, gcp.NewComputeFirewallRuleResource(f))
		}

		for firewalls.NextPageToken != "" {
			firewalls, err = c.computeClient.Firewalls.List(id).PageToken(firewalls.NextPageToken).Do()
			if err != nil {
				glog.Fatalf("Error retrieving firewall rules from project '%v': %v", id, err)
			}

			for _, f := range firewalls.Items {
				res.Firewalls = append(res.Firewalls, gcp.NewComputeFirewallRuleResource(f))
			}
		}

		// Get aggregated IPs for the projects
		aggregateAddressesList, err := c.computeClient.Addresses.AggregatedList(id).Do()
		if err != nil {
			glog.Fatalf("Error retrieving addresses from project '%v': %v", id, err)
		}

		scopedAddresses := aggregateAddressesList.Items
		for _, scope := range scopedAddresses {
			for _, a := range scope.Addresses {
				res.Addresses = append(res.Addresses, gcp.NewComputeAddressResource(a))
			}
		}

		results <- res
	}

	// Setup worker pool
	projectIDs := make(chan string, len(c.resourceprojects))
	results := make(chan networkCallResult, len(c.resourceprojects))
	numWorkers := len(c.resourceprojects)
	for w := 0; w < numWorkers; w++ {
		go worker(projectIDs, results)
	}

	// Feed the workers and collect the storage info
	for _, p := range c.resourceprojects {
		projectIDs <- p.ProjectId
	}

	// Collect the info
	for i := 0; i < numWorkers; i++ {
		res := <-results
		c.networks[res.ProjectID] = res.Networks
		c.subnetworks[res.ProjectID] = res.Subnetworks
		c.firewalls[res.ProjectID] = res.Firewalls
		c.addresses[res.ProjectID] = res.Addresses
	}

	return nil
}

type networkCallResult struct {
	ProjectID   string
	Networks    []*gcp.ComputeNetworkResource
	Subnetworks []*gcp.ComputeSubnetworkResource
	Firewalls   []*gcp.ComputeFirewallRuleResource
	Addresses   []*gcp.ComputeAddressResource
}

// GenerateComputeNetworkReports signals the client to process ComputeNetworkResource's for reports.
// If there are no networks found in the configuration, no reports will be created.
func (c *Client) GenerateComputeNetworkReports() (reports []report.Report, err error) {

	reports = []report.Report{}
	typ := "compute_network"

	for _, p := range c.computeprojects {
		projectID := p.Name()

		for _, n := range c.networks[p.Name()] {
			r := report.NewReport(
				typ,
				fmt.Sprintf("Network %v in Project %v", n.Name(), p.Name()),
			)
			r.Data, err = n.Marshal()
			if err != nil {
				glog.Fatalf("Failed to marshal network: %v", err)
			}

			// The default network should not be used in projects
			defaultNetworkControl := report.NewCISControl(
				"3.1",
				fmt.Sprintf("Project %v should not have a default network", p.Name()),
			)
			if n.IsDefault() {
				defaultNetworkControl.Error = fmt.Sprintf("Network %v is the default network", n.Name())
			} else {
				defaultNetworkControl.Passed()
			}

			// Legacy networks should not be used
			legacyNetworkControl := report.NewCISControl(
				"3.2",
				fmt.Sprintf("Project %v should not have legacy networks", p.Name()),
			)
			if n.IsLegacy() {
				legacyNetworkControl.Error = fmt.Sprintf("Network %v is a legacy network", n.Name())
			} else {
				legacyNetworkControl.Passed()
			}

			r.AddControls(defaultNetworkControl, legacyNetworkControl)

			reports = append(reports, r)
			c.incrementMetrics(typ, n.Name(), r.Status(), projectID)
		}
	}

	return
}

// GenerateComputeSubnetworkReports signals the client to process ComputeSubnetworkResource's for reports.
// If there are no subnetworks found in the configuration, no reports will be created.
func (c *Client) GenerateComputeSubnetworkReports() (reports []report.Report, err error) {

	reports = []report.Report{}
	typ := "compute_subnetwork"

	for _, p := range c.computeprojects {
		projectID := p.Name()

		for _, s := range c.subnetworks[p.Name()] {
			r := report.NewReport(
				typ,
				fmt.Sprintf("Subnetwork %v in region %v for Project %v", s.Name(), s.Region(), p.Name()),
			)
			r.Data, err = s.Marshal()
			if err != nil {
				glog.Fatalf("Failed to marshal subnetwork: %v", err)
			}

			privateAccessControl := report.NewCISControl(
				"3.8",
				fmt.Sprintf("Subnetwork %v should have Private Google Access enabled", s.Name()),
			)
			if s.IsPrivateGoogleAccessEnabled() {
				privateAccessControl.Passed()
			} else {
				privateAccessControl.Error = fmt.Sprintf("Subnetwork %v does not have Private Google Access enabled", s.Name())
			}

			flowLogsControl := report.NewCISControl(
				"3.9",
				fmt.Sprintf("Subnetwork %v should have VPC flow logs enabled", s.Name()),
			)
			if s.IsFlowLogsEnabled() {
				flowLogsControl.Passed()
			} else {
				flowLogsControl.Error = fmt.Sprintf("Subnetwork %v does not have VPC flow logs enabled", s.Name())
			}

			r.AddControls(privateAccessControl, flowLogsControl)

			reports = append(reports, r)
			c.incrementMetrics(typ, s.Name(), r.Status(), projectID)
		}
	}

	return
}

// GenerateComputeFirewallRuleReports signals the client to process ComputeFirewallRuleResource's for reports.
// If there are no network keys configured in the configuration, no reports will be created.
func (c *Client) GenerateComputeFirewallRuleReports() (reports []report.Report, err error) {

	reports = []report.Report{}
	typ := "compute_firewall_rule"

	for _, p := range c.computeprojects {
		projectID := p.Name()

		for _, f := range c.firewalls[p.Name()] {
			r := report.NewReport(
				typ,
				fmt.Sprintf("Network %v Firewall Rule %v", f.Network(), f.Name()),
			)
			r.Data, err = f.Marshal()
			if err != nil {
				glog.Fatalf("Failed to marshal firewall rule: %v", err)
			}

			// SSH access from the internet should not be allowed
			sshControl := report.NewCISControl(
				"3.6",
				"SSH should not be allowed from the internet",
			)
			if (f.AllowsProtocolPort("TCP", "22") || f.AllowsProtocolPort("UDP", "22")) && f.AllowsSourceRange("0.0.0.0/0") {
				sshControl.Error = fmt.Sprintf("%v allows SSH from the internet", f.Name())
			} else {
				sshControl.Passed()
			}

			// RDP access from the internet should not be allowed
			rdpControl := report.NewCISControl(
				"3.7",
				"RDP should not be allowed from the internet",
			)
			if (f.AllowsProtocolPort("TCP", "3389") || f.AllowsProtocolPort("UDP", "3389")) && f.AllowsSourceRange("0.0.0.0/0") {
				rdpControl.Error = fmt.Sprintf("%v allows RDP froom the internet", f.Name())
			} else {
				rdpControl.Passed()
			}

			r.AddControls(sshControl, rdpControl)
			reports = append(reports, r)
			c.incrementMetrics(typ, f.Name(), r.Status(), projectID)
		}
	}

	return
}

// GenerateComputeAddressReports signals the client to process ComputeAddressResource's for reports.
// If there are no network keys configured in the configuration, no reports will be created.
func (c *Client) GenerateComputeAddressReports() (reports []report.Report, err error) {

	reports = []report.Report{}
	typ := "compute_address"

	for _, p := range c.computeprojects {

		projectID := p.Name()
		for _, a := range c.addresses[projectID] {

			r := report.NewReport(
				typ,
				fmt.Sprintf("Compute Address %v", a.Name()),
			)
			r.Data, err = a.Marshal()
			if err != nil {
				glog.Fatalf("Failed to marshal compute address: %v", err)
			}

			reports = append(reports, r)
			c.incrementMetrics(typ, a.Name(), r.Status(), projectID)
		}
	}

	return
}
