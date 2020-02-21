package client

import (
	"fmt"

	"github.com/Unity-Technologies/nemesis/pkg/report"
	"github.com/Unity-Technologies/nemesis/pkg/resource/gcp"
	"github.com/Unity-Technologies/nemesis/pkg/utils"
	"github.com/golang/glog"
	compute "google.golang.org/api/compute/v1"
)

func (c *Client) getZoneNames() ([]string, error) {

	// Get the zones from the first project. Should be the same for all projects
	zoneNames := make([]string, 0)
	var zones *compute.ZoneList
	var err error
	for i := 0; i < len(c.resourceprojects); i++ {

		zones, err = c.computeClient.Zones.List(c.resourceprojects[i].ProjectId).Do()
		if err == nil {
			// We got a valid list of zones, so skip
			break
		}
	}

	if zones == nil {
		err = fmt.Errorf("Error retrieving zones list from any project: %v", err)
		return nil, err
	}

	for _, z := range zones.Items {
		zoneNames = append(zoneNames, z.Name)
	}

	return zoneNames, nil
}

func (c *Client) getRegionNames() ([]string, error) {

	// Get the regions from the first project. Should be the same for all projects
	regionNames := make([]string, 0)
	var regions *compute.RegionList
	var err error
	for i := 0; i < len(c.resourceprojects); i++ {

		regions, err = c.computeClient.Regions.List(c.resourceprojects[i].ProjectId).Do()
		if err == nil {
			break
		}
	}

	if regions == nil {
		err = fmt.Errorf("Error retrieving regions list from any project: %v", err)
		return nil, err
	}

	for _, r := range regions.Items {
		regionNames = append(regionNames, r.Name)
	}

	return regionNames, nil
}

// GetComputeResources launches the process retrieving compute resources
func (c *Client) GetComputeResources() error {

	defer utils.Elapsed("GetComputeResources")()

	// Get list of all projects.
	projects := c.resourceprojects

	zoneNames, err := c.getZoneNames()
	if err != nil {
		glog.Fatalf("%v", err)
	}

	// Create a worker pool for querying zones a bit faster
	zoneWorker := func(projectID string, id int, zones <-chan string, results chan<- []*gcp.ComputeInstanceResource) {

		// For each zone passed to this worker
		for z := range zones {

			// Create the list of instance resources to be retrieved for this zone
			instanceResources := []*gcp.ComputeInstanceResource{}
			res, err := c.computeClient.Instances.List(projectID, z).Do()
			if err != nil {
				glog.Fatalf("Error retrieving project %v's instances in zone %v: %v", projectID, z, err)
			}

			// Create the resource
			for _, i := range res.Items {
				instanceResources = append(instanceResources, gcp.NewComputeInstanceResource(i))
			}

			// Pass the list of zones across the channel
			results <- instanceResources
		}
	}

	// For each project, collect information
	for _, p := range projects {

		// Only save the project ID
		projectID := p.ProjectId

		// Check that the compute API is enabled for the project. If not, then skip auditing compute resources for the project entirely
		if !c.isServiceEnabled(projectID, "compute.googleapis.com") {
			continue
		}

		// Get the compute API's version of the project
		project, err := c.computeClient.Projects.Get(projectID).Do()
		if err != nil {
			glog.Fatalf("Error retrieving project %v's metadata: %v", projectID, err)
		}

		// Store the project resource
		c.computeprojects = append(c.computeprojects, gcp.NewComputeProjectResource(project))

		// Store the project's compute metadata resource
		c.computeMetadatas[projectID] = gcp.NewComputeProjectMetadataResource(project.CommonInstanceMetadata)

		instances := []*gcp.ComputeInstanceResource{}
		jobs := make(chan string, len(zoneNames))
		results := make(chan []*gcp.ComputeInstanceResource, len(zoneNames))

		// Create the zone worker pool
		for w := 0; w < len(zoneNames); w++ {
			go zoneWorker(projectID, w, jobs, results)
		}

		// Feed the zone names
		for _, z := range zoneNames {
			jobs <- z
		}
		close(jobs)

		// Retrieve the full instances list
		for i := 0; i < len(zoneNames); i++ {
			instances = append(instances, <-results...)
		}

		// Store the project instance's resources
		c.instances[projectID] = instances

	}

	return nil
}

// GenerateComputeMetadataReports signals the client to process ComputeMetadataResource's for reports.
// If there are no metadata keys configured in the configuration, no reports will be created.
func (c *Client) GenerateComputeMetadataReports() (reports []report.Report, err error) {

	reports = []report.Report{}
	typ := "compute_metadata"

	// For each project compute metadata, generate one report
	for _, p := range c.computeprojects {
		projectID := p.Name()
		projectMetadata := c.computeMetadatas[projectID]
		r := report.NewReport(typ, fmt.Sprintf("Project %v Common Instance Metadata", projectID))

		// Always connect the data for the report with the source data
		if r.Data, err = projectMetadata.Marshal(); err != nil {
			glog.Fatalf("Failed to marshal project metadata: %v", err)
		}

		blockSSHKeys := report.NewCISControl(
			"4.2",
			"Project metadata should include 'block-project-ssh-keys' and be set to 'true'",
		)
		res, err := projectMetadata.KeyValueEquals("block-project-ssh-keys", "true")
		if err != nil {
			blockSSHKeys.Error = err.Error()
		} else {
			if res {
				blockSSHKeys.Passed()
			} else {
				glog.Fatalf("Could not determine the state of project %v's metadata, aborting...", projectID)
			}
		}

		osLogin := report.NewCISControl(
			"4.3",
			"Project metadata should include the key 'enable-oslogin' with value set to 'true'",
		)
		res, err = projectMetadata.KeyValueEquals("enable-oslogin", "true")
		if err != nil {
			osLogin.Error = err.Error()
		} else {
			if res {
				osLogin.Passed()
			} else {
				glog.Fatalf("Could not determine the state of project %v's metadata, aborting...", projectID)
			}
		}

		// Dynamic serial port access should be denied
		// serial-port-enable is a special case, where absence of the key is equivalent to disabling serial port access
		serialPortAccess := report.NewCISControl(
			"4.4",
			"Project metadata should include the key 'serial-port-enable' with value set to '0'",
		)
		if projectMetadata.KeyAbsent("serial-port-enable") {
			serialPortAccess.Passed()
		} else {
			res, err = projectMetadata.KeyValueEquals("serial-port-enable", "0")
			if err != nil {
				serialPortAccess.Error = err.Error()
			} else {
				if res {
					serialPortAccess.Passed()
				} else {
					glog.Fatalf("Could not determine the state of project %v's metadata, aborting...", projectID)
				}
			}
		}

		legacyMetadata := report.NewControl(
			"Ensure legacy metadata endpoints are not enabled for VM Instance",
			"Project metadata should include the key 'disable-legacy-endpoints' with value set to 'true'",
		)

		res, err = projectMetadata.KeyValueEquals("disable-legacy-endpoints", "true")
		if err != nil {
			legacyMetadata.Error = err.Error()
		} else {
			if res {
				legacyMetadata.Passed()
			} else {
				glog.Fatalf("Could not determine the state of project %v's metadata, aborting...", projectID)
			}
		}

		// Append the control to this resource's report
		r.AddControls(blockSSHKeys, osLogin, serialPortAccess, legacyMetadata)

		// Append the resource's report to our final list
		reports = append(reports, r)
		c.incrementMetrics(typ, projectID, r.Status(), projectID)
	}

	return
}

// GenerateComputeInstanceReports signals the client to process ComputeInstanceResource's for reports.
// If there are keys configured for instances in the configuration, no reports will be created.
func (c *Client) GenerateComputeInstanceReports() (reports []report.Report, err error) {

	reports = []report.Report{}
	typ := "compute_instance"

	for _, p := range c.computeprojects {

		projectID := p.Name()
		instanceResources := c.instances[projectID]
		metadata := c.computeMetadatas[projectID]

		for _, i := range instanceResources {
			r := report.NewReport(typ, fmt.Sprintf("Project %v Compute Instance %v", projectID, i.Name()))
			if r.Data, err = i.Marshal(); err != nil {
				glog.Fatalf("Failed to marshal compute instance: %v", err)
			}

			// Make sure the number of Network Interfaces matches what is expected
			numNicsControl := report.NewControl(
				fmt.Sprintf("numNetworkInterfaces=%v", *flagComputeInstanceNumInterfaces),
				fmt.Sprintf("Compute Instance should have a number of network interfaces equal to %v", *flagComputeInstanceNumInterfaces),
			)

			_, err := i.HasNumNetworkInterfaces(*flagComputeInstanceNumInterfaces)
			if err != nil {
				numNicsControl.Error = err.Error()
			} else {
				numNicsControl.Passed()
			}

			// Measure whether a NAT ip address is expected
			natIPControl := report.NewControl(
				fmt.Sprintf("hasNatIP=%v", *flagComputeInstanceAllowNat),
				fmt.Sprintf("Compute Instance should have a NAT ip configured: %v", *flagComputeInstanceAllowNat),
			)
			if i.HasNatIP() {
				if *flagComputeInstanceAllowNat {
					// External IP exists, and we want it to exist
					natIPControl.Passed()
				} else {
					// External IP exists, but we don't want it to exist
					natIPControl.Error = "Compute Instance has NAT IP address, but should not"
				}
			} else {
				// External IP does not exist, and we don't want it to exist
				if !*flagComputeInstanceAllowNat {
					natIPControl.Passed()
				} else {
					// It doesn't exist but we wanted it to exist
					natIPControl.Error = "Compute Instance does not have a NAT IP address, but it should"
				}
			}

			// Default compute service account should not be used to launch instances
			defaultSA := report.NewCISControl(
				"4.1",
				"Compute Instance should not use the project default compute service account",
			)

			if !i.UsesDefaultServiceAccount() {
				defaultSA.Passed()
			} else {
				defaultSA.Error = "Compute instance uses a default compute service account"
			}

			wrapMetadata := func(meta *gcp.ComputeProjectMetadataResource, i *gcp.ComputeInstanceResource, k string, v string) (bool, error) {
				result, err := meta.KeyValueEquals(k, v)
				if err != nil {
					result, err = i.KeyValueEquals(k, v)
				}
				return result, err
			}

			// Project-wide SSH keys should not be used to access instances
			blockSSHKeys := report.NewCISControl(
				"4.2",
				"Compute Instance metadata should include 'block-project-ssh-keys' and be set to 'true'",
			)
			res, err := wrapMetadata(metadata, i, "block-project-ssh-keys", "true")
			if err != nil {
				blockSSHKeys.Error = err.Error()
			} else {
				if res {
					blockSSHKeys.Passed()
				} else {
					glog.Fatalf("Could not determine the state of instance %v's metadata, aborting...", projectID)
				}
			}

			// Ensure os-login is enabled
			osLogin := report.NewCISControl(
				"4.3",
				"Compute Instance metadata should include the key 'enable-oslogin' with value set to 'true'",
			)
			res, err = wrapMetadata(metadata, i, "enable-oslogin", "true")
			if err != nil {
				osLogin.Error = err.Error()
			} else {
				if res {
					osLogin.Passed()
				} else {
					glog.Fatalf("Could not determine the state of instance %v's metadata, aborting...", projectID)
				}
			}

			// Dynamic serial port access should be denied
			// serial-port-enable is a special case, where absence of the key is equivalent to disabling serial port access
			serialPortAccess := report.NewCISControl(
				"4.4",
				"Compute Instance metadata should include the key 'serial-port-enable' with value set to '0'",
			)
			if metadata.KeyAbsent("serial-port-enable") && i.KeyAbsent("serial-port-enable") {
				serialPortAccess.Passed()
			} else {
				res, err = wrapMetadata(metadata, i, "serial-port-enable", "0")
				if err != nil {
					serialPortAccess.Error = err.Error()
				} else {
					if res {
						serialPortAccess.Passed()
					} else {
						glog.Fatalf("Could not determine the state of instance %v's metadata, aborting...", projectID)
					}
				}
			}

			// IP forwarding should not be enabled
			ipForwarding := report.NewCISControl(
				"4.5",
				"Compute Instance should not allow ip forwarding of packets",
			)
			if !i.HasIPForwardingEnabled() {
				ipForwarding.Passed()
			} else {
				if *flagComputeInstanceAllowIPForwarding {
					ipForwarding.Passed()
				} else {
					ipForwarding.Error = "Compute Instance allows IP Forwarding"
				}
			}

			// Disks should be using Customer Supplied Encryption Keys
			csekDisk := report.NewCISControl(
				"4.6",
				"Compute Instance should be encrypted with a CSEK",
			)
			if err := i.UsesCustomerSuppliedEncryptionKeys(); err != nil {
				csekDisk.Passed()
			} else {
				csekDisk.Error = "Compute Instance does not have CSEK encryption on disk"
			}

			r.AddControls(
				numNicsControl,
				natIPControl,
				defaultSA,
				blockSSHKeys,
				osLogin,
				serialPortAccess,
				ipForwarding,
			)

			// Add the instance resource report to the final list of reports
			reports = append(reports, r)
			totalResourcesCounter.Inc()
			c.incrementMetrics(typ, i.Name(), r.Status(), projectID)
		}
	}

	return
}
