package client

import (
	"fmt"

	"github.com/UnityTech/nemesis/pkg/report"
	"github.com/UnityTech/nemesis/pkg/resource/gcp"
	"github.com/UnityTech/nemesis/pkg/utils"
	"github.com/golang/glog"
)

type projectParam struct {
	projectID string
}

func (p projectParam) Get() (key, value string) {
	return "projectId", p.projectID
}

// GetContainerResources launches the process retrieving container cluster and nodepool resources
func (c *Client) GetContainerResources() error {

	defer utils.Elapsed("GetContainerResources")()

	clustersService := c.containerClient.Projects.Locations.Clusters

	// Create a short-lived goroutine for retrieving a project's container clusters
	worker := func(projectIDs <-chan string, results chan<- containerCallResult) {

		id := <-projectIDs
		res := containerCallResult{ProjectID: id, Clusters: []*gcp.ContainerClusterResource{}}

		// Check that the container API is enabled. If not, don't audit container resources in the project
		if !c.isServiceEnabled(id, "container.googleapis.com") {
			results <- res
			return
		}

		// Perform the query
		location := fmt.Sprintf("projects/%v/locations/-", id)
		clusters, err := clustersService.List(location).Do()
		if err != nil {
			glog.Fatalf("Error retrieving container clusters in project %v: %v", id, err)
		}

		for _, cluster := range clusters.Clusters {
			res.Clusters = append(res.Clusters, gcp.NewContainerClusterResource(cluster))
		}

		results <- res
	}

	// Setup worker pool
	projectIDs := make(chan string, len(c.resourceprojects))
	results := make(chan containerCallResult, len(c.resourceprojects))
	numWorkers := len(c.resourceprojects)
	for w := 0; w < numWorkers; w++ {
		go worker(projectIDs, results)
	}

	// Feed the workers and collect the cluster info
	for _, p := range c.resourceprojects {
		projectIDs <- p.ProjectId
	}

	// Collect the info
	for i := 0; i < numWorkers; i++ {
		res := <-results
		c.clusters[res.ProjectID] = res.Clusters
	}

	return nil
}

type containerCallResult struct {
	ProjectID string
	Clusters  []*gcp.ContainerClusterResource
}

// GenerateContainerClusterReports signals the client to process ContainerClusterResource's for reports.
func (c *Client) GenerateContainerClusterReports() (reports []report.Report, err error) {

	reports = []report.Report{}
	typ := "container_cluster"

	for _, p := range c.computeprojects {
		projectID := p.Name()
		for _, cluster := range c.clusters[projectID] {
			r := report.NewReport(
				typ,
				fmt.Sprintf("Project %v Container Cluster %v", projectID, cluster.Name()),
			)
			if r.Data, err = cluster.Marshal(); err != nil {
				glog.Fatalf("Failed to marshal container cluster: %v", err)
			}

			// Clusters should have stackdriver logging enabled
			sdLogging := report.NewCISControl(
				"7.1",
				fmt.Sprintf("Cluster %v should have Stackdriver logging enabled", cluster.Name()),
			)
			if !cluster.IsStackdriverLoggingEnabled() {
				sdLogging.Error = "Stackdriver logging is not enabled"
			} else {
				sdLogging.Passed()
			}

			// Clusters should have stackdriver monitoring enabled
			sdMonitoring := report.NewCISControl(
				"7.2",
				fmt.Sprintf("Cluster %v should have Stackdriver monitoring enabled", cluster.Name()),
			)
			if !cluster.IsStackdriverMonitoringEnabled() {
				sdMonitoring.Error = "Stackdriver monitoring is not enabled"
			} else {
				sdMonitoring.Passed()
			}
			// Clusters should not enable Attribute-Based Access Control (ABAC)
			abac := report.NewCISControl(
				"7.3",
				fmt.Sprintf("Cluster %v should have Legacy ABAC disabled", cluster.Name()),
			)
			if !cluster.IsAbacDisabled() {
				abac.Error = "Cluster has Legacy ABAC enabled when it should not"
			} else {
				abac.Passed()
			}

			// Clusters should use Master authorized networks
			masterAuthNetworks := report.NewCISControl(
				"7.4",
				fmt.Sprintf("Cluster %v should have Master authorized networks enabled", cluster.Name()),
			)
			if !cluster.IsMasterAuthorizedNetworksEnabled() {
				masterAuthNetworks.Error = "Cluster does not have Master Authorized Networks enabled"
			} else {
				masterAuthNetworks.Passed()
			}

			// Clusters should not enable Kubernetes Dashboard
			dashboard := report.NewCISControl(
				"7.6",
				fmt.Sprintf("Cluster %v should have Kubernetes Dashboard disabled", cluster.Name()),
			)
			if !cluster.IsDashboardAddonDisabled() {
				dashboard.Error = "Cluster has Kubernetes Dashboard add-on enabled when it should not"
			} else {
				dashboard.Passed()
			}

			// Clusters should not allow authentication with username/password
			masterAuthPassword := report.NewCISControl(
				"7.10",
				fmt.Sprintf("Cluster %v should not have a password configured", cluster.Name()),
			)
			if !cluster.IsMasterAuthPasswordDisabled() {
				masterAuthPassword.Error = "Cluster has a password configured to allow basic auth when it should not"
			} else {
				masterAuthPassword.Passed()
			}

			// Clusters should enable network policies (pod-to-pod policy)
			networkPolicy := report.NewCISControl(
				"7.11",
				fmt.Sprintf("Cluster %v should have Network Policy addon enabled", cluster.Name()),
			)
			if !cluster.IsNetworkPolicyAddonEnabled() {
				networkPolicy.Error = "Cluster does not have Network Policy addon enabled when it should"
			} else {
				networkPolicy.Passed()
			}

			// Clusters should not allow authentication with client certificates
			clientCert := report.NewCISControl(
				"7.12",
				fmt.Sprintf("Cluster %v should not issue client certificates", cluster.Name()),
			)
			if !cluster.IsClientCertificateDisabled() {
				clientCert.Error = "Cluster has ABAC enabled when it should not"
			} else {
				clientCert.Passed()
			}

			// Clusters should be launched as VPC-native and use Pod Alias IP ranges
			aliasIps := report.NewCISControl(
				"7.13",
				fmt.Sprintf("Cluster %v should use VPC-native alias IP ranges", cluster.Name()),
			)
			if !cluster.IsAliasIPEnabled() {
				aliasIps.Error = "Cluster is not using VPC-native alias IP ranges"
			} else {
				aliasIps.Passed()
			}

			// Cluster master should not be accessible over public IP
			privateMaster := report.NewCISControl(
				"7.15",
				fmt.Sprintf("Cluster %v master should be private and not accessible over public IP", cluster.Name()),
			)
			if !cluster.IsMasterPrivate() {
				privateMaster.Error = "Cluster master is not private and is routeable on public internet"
			} else {
				privateMaster.Passed()
			}

			// Cluster nodes should not be accessible over public IP
			privateNodes := report.NewCISControl(
				"7.15",
				fmt.Sprintf("Cluster %v nodes should be private and not accessible over public IPs", cluster.Name()),
			)
			if !cluster.IsNodesPrivate() {
				privateNodes.Error = "Cluster nodes are not private and are routable on the public internet"
			} else {
				privateNodes.Passed()
			}

			// Cluster should not be launched using the default compute service account
			defaultSA := report.NewCISControl(
				"7.17",
				fmt.Sprintf("Cluster %v should not be using the default compute service account", cluster.Name()),
			)
			if cluster.IsUsingDefaultServiceAccount() {
				defaultSA.Error = "Cluster is using the default compute service account"
			} else {
				defaultSA.Passed()
			}

			// Cluster should be using minimal OAuth scopes
			oauthScopes := report.NewCISControl(
				"7.18",
				fmt.Sprintf("Cluster %v should be launched with minimal OAuth scopes", cluster.Name()),
			)
			if _, err := cluster.IsUsingMinimalOAuthScopes(); err != nil {
				oauthScopes.Error = err.Error()
			} else {
				oauthScopes.Passed()
			}

			r.AddControls(sdLogging, sdMonitoring, abac, masterAuthNetworks, dashboard, masterAuthPassword, networkPolicy, clientCert, aliasIps, privateMaster, privateNodes, defaultSA, oauthScopes)
			reports = append(reports, r)
			c.incrementMetrics(typ, cluster.Name(), r.Status(), projectID)
		}
	}

	return
}

// GenerateContainerNodePoolReports signals the client to process ContainerNodePoolResource's for reports.
func (c *Client) GenerateContainerNodePoolReports() (reports []report.Report, err error) {
	reports = []report.Report{}
	typ := "container_nodepool"

	for _, p := range c.computeprojects {
		projectID := p.Name()
		for _, nodepool := range c.nodepools[projectID] {
			r := report.NewReport(
				typ,
				fmt.Sprintf("Project %v Container Node Pool %v", projectID, nodepool.Name()),
			)
			if r.Data, err = nodepool.Marshal(); err != nil {
				glog.Fatalf("Failed to marshal container node pool: %v", err)
			}

			// Nodepools should not allow use of legacy metadata APIs
			legacyAPI := report.NewControl(
				"disableLegacyMetadataAPI",
				fmt.Sprintf("Node pool %v should have legacy metadata API disabled", nodepool.Name()),
			)
			if _, err := nodepool.IsLegacyMetadataAPIDisabled(); err != nil {
				legacyAPI.Error = err.Error()
			} else {
				legacyAPI.Passed()
			}

			// Node pools should be configured for automatic repairs
			repair := report.NewCISControl(
				"7.7",
				fmt.Sprintf("Node pool %v should have automatic repairs enabled", nodepool.Name()),
			)
			if !nodepool.IsAutoRepairEnabled() {
				repair.Error = "Automatic node repair is not enabled"
			} else {
				repair.Passed()
			}

			// Node pools should be configured for automatic upgrades
			upgrade := report.NewCISControl(
				"7.8",
				fmt.Sprintf("Node pool %v should have automatic upgrades enabled", nodepool.Name()),
			)
			if !nodepool.IsAutoUpgradeEnabled() {
				upgrade.Error = "Automatic node upgrade is not enabled"
			} else {
				upgrade.Passed()
			}

			// Node pools should be using COS (Google Container OS)
			cos := report.NewCISControl(
				"7.9",
				fmt.Sprintf("Node pool %v should be using COS", nodepool.Name()),
			)
			if _, err := nodepool.CheckDistributionTypeIs("COS"); err != nil {
				cos.Error = err.Error()
			} else {
				cos.Passed()
			}

			r.AddControls(legacyAPI, repair, upgrade, cos)
			reports = append(reports, r)
			c.incrementMetrics(typ, nodepool.Name(), r.Status(), projectID)
		}
	}

	return
}
