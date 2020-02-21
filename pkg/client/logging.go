package client

import (
	"context"
	"fmt"

	"github.com/Unity-Technologies/nemesis/pkg/report"

	"github.com/Unity-Technologies/nemesis/pkg/resource/gcp"

	"google.golang.org/api/iterator"

	"github.com/Unity-Technologies/nemesis/pkg/utils"
	loggingpb "google.golang.org/genproto/googleapis/logging/v2"
)

// GetLoggingResources returns the logging config and log-based metric configurations
func (c *Client) GetLoggingResources() error {

	defer utils.Elapsed("GetLoggingResources")()

	worker := func(projectIDs <-chan string, results chan<- loggingClientResult) {

		id := <-projectIDs
		parent := fmt.Sprintf("projects/%s", id)

		ctx := context.Background()
		res := loggingClientResult{ProjectID: id, LogSinks: []*gcp.LoggingSinkResource{}}

		// Grab the project's logging sinks
		req1 := loggingpb.ListSinksRequest{
			Parent: parent,
		}
		it1 := c.logConfigClient.ListSinks(ctx, &req1)
		for {
			s, done := it1.Next()
			if done == iterator.Done {
				break
			}

			res.LogSinks = append(res.LogSinks, gcp.NewLoggingSinkResource(s))
		}

		// Grab the project's log-based metrics
		req2 := loggingpb.ListLogMetricsRequest{
			Parent: parent,
		}
		it2 := c.logMetricClient.ListLogMetrics(ctx, &req2)
		for {
			m, done := it2.Next()
			if done == iterator.Done {
				break
			}

			res.LogMetrics = append(res.LogMetrics, gcp.NewLoggingMetricResource(m))
		}

		results <- res
	}

	projectIDs := make(chan string, len(c.resourceprojects))
	results := make(chan loggingClientResult, len(c.resourceprojects))
	numWorkers := len(c.resourceprojects)
	for w := 0; w < numWorkers; w++ {
		go worker(projectIDs, results)
	}

	for _, p := range c.resourceprojects {
		projectIDs <- p.ProjectId
	}

	for i := 0; i < numWorkers; i++ {
		res := <-results
		c.logSinks[res.ProjectID] = res.LogSinks
		c.logMetrics[res.ProjectID] = res.LogMetrics
	}

	return nil
}

type loggingClientResult struct {
	ProjectID  string
	LogSinks   []*gcp.LoggingSinkResource
	LogMetrics []*gcp.LoggingMetricResource
}

// GenerateLoggingReports signals the client to process LoggingResources for reports.
// TODO - implement CIS 2.3
func (c *Client) GenerateLoggingReports() (reports []report.Report, err error) {

	reports = []report.Report{}

	for _, p := range c.computeprojects {

		r := report.NewReport(
			"logging_configuration",
			fmt.Sprintf("Project %s Logging Configuration", p.Name()),
		)

		// At least one sink in a project should ship all logs _somewhere_
		exportLogs := report.NewCISControl(
			"2.2",
			fmt.Sprintf("Project %s should have at least one export configured with no filters", p.Name()),
		)
		isExported := false
		for _, s := range c.logSinks[p.Name()] {
			isExported = s.ShipsAllLogs()
			if isExported {
				break
			}
		}
		if !isExported {
			exportLogs.Error = fmt.Sprintf("There is no logging sink that exports all logs for project %s", p.Name())
		} else {
			exportLogs.Passed()
		}

		// Helper function to determine if a list of log-based metrics contains a specific filter
		metricExists := func(metrics []*gcp.LoggingMetricResource, filter string) bool {
			for _, m := range metrics {
				if m.FilterMatches(filter) {
					return true
				}
			}
			return false
		}

		// Monitor Project Ownership Changes
		projectOwnerChanges := report.NewCISControl(
			"2.4",
			fmt.Sprintf("Project %s should monitor ownership changes", p.Name()),
		)

		// Monitor Project Audit Configuration Changes
		auditConfigChanges := report.NewCISControl(
			"2.5",
			fmt.Sprintf("Project %s should monitor audit log configuration changes", p.Name()),
		)

		// Monitor Project Custom Role Changes
		customRoleChanges := report.NewCISControl(
			"2.6",
			fmt.Sprintf("Project %s should monitor custom IAM role changes", p.Name()),
		)

		// Monitor VPC Firewall Changes
		vpcFirewallChanges := report.NewCISControl(
			"2.7",
			fmt.Sprintf("Project %s should monitor VPC firewall changes", p.Name()),
		)

		// Monitor VPC Route Changes
		vpcRouteChanges := report.NewCISControl(
			"2.8",
			fmt.Sprintf("Project %s should monitor VPC route changes", p.Name()),
		)

		// Monitor General Changes to VPC Configuration
		vpcNetworkChanges := report.NewCISControl(
			"2.9",
			fmt.Sprintf("Project %s should monitor VPC network changes", p.Name()),
		)

		// Monitor GCS IAM Policy Changes
		gcsIamChanges := report.NewCISControl(
			"2.10",
			fmt.Sprintf("Project %s should monitor GCS IAM changes", p.Name()),
		)

		// Monitor SQL Configuration Changes
		sqlConfigChanges := report.NewCISControl(
			"2.11",
			fmt.Sprintf("Project %s should monitor SQL config changes", p.Name()),
		)

		metricControls := []struct {
			Control *report.Control
			Filter  string
		}{
			{
				Control: &projectOwnerChanges,
				Filter:  `(protoPayload.serviceName="cloudresourcemanager.googleapis.com") AND (ProjectOwnership OR projectOwnerInvitee) OR (protoPayload.serviceData.policyDelta.bindingDeltas.action="REMOVE" AND protoPayload.serviceData.policyDelta.bindingDeltas.role="roles/owner") OR (protoPayload.serviceData.policyDelta.bindingDeltas.action="ADD" AND protoPayload.serviceData.policyDelta.bindingDeltas.role="roles/owner")`,
			},
			{
				Control: &auditConfigChanges,
				Filter:  `protoPayload.methodName="SetIamPolicy" AND protoPayload.serviceData.policyDelta.auditConfigDeltas:*`,
			},
			{
				Control: &customRoleChanges,
				Filter:  `resource.type="iam_role" AND protoPayload.methodName = "google.iam.admin.v1.CreateRole" OR protoPayload.methodName="google.iam.admin.v1.DeleteRole" OR protoPayload.methodName="google.iam.admin.v1.UpdateRole"`,
			},
			{
				Control: &vpcFirewallChanges,
				Filter:  `resource.type="gce_firewall_rule" AND jsonPayload.event_subtype="compute.firewalls.patch" OR jsonPayload.event_subtype="compute.firewalls.insert"`,
			},
			{
				Control: &vpcRouteChanges,
				Filter:  `resource.type="gce_route" AND jsonPayload.event_subtype="compute.routes.delete" OR jsonPayload.event_subtype="compute.routes.insert"`,
			},
			{
				Control: &vpcNetworkChanges,
				Filter:  `resource.type=gce_network AND jsonPayload.event_subtype="compute.networks.insert" OR jsonPayload.event_subtype="compute.networks.patch" OR jsonPayload.event_subtype="compute.networks.delete" OR jsonPayload.event_subtype="compute.networks.removePeering" OR jsonPayload.event_subtype="compute.networks.addPeering"`,
			},
			{
				Control: &gcsIamChanges,
				Filter:  `resource.type=gcs_bucket AND protoPayload.methodName="storage.setIamPermissions"`,
			},
			{
				Control: &sqlConfigChanges,
				Filter:  `protoPayload.methodName="cloudsql.instances.update"`,
			},
		}

		for _, m := range metricControls {
			if metricExists(c.logMetrics[p.Name()], m.Filter) {
				m.Control.Passed()
			} else {
				m.Control.Error = fmt.Sprintf("Project %s does not have the following filter monitored: %s", p.Name(), m.Filter)
			}
		}

		r.AddControls(
			exportLogs,
			projectOwnerChanges,
			auditConfigChanges,
			customRoleChanges,
			vpcFirewallChanges,
			vpcRouteChanges,
			vpcNetworkChanges,
			gcsIamChanges,
			sqlConfigChanges,
		)
		reports = append(reports, r)
	}

	return
}
