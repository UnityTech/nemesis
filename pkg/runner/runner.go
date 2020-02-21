// Package runner executes a configured audit
package runner

import (
	"flag"

	"github.com/UnityTech/nemesis/pkg/client"
	"github.com/UnityTech/nemesis/pkg/report"
	"github.com/UnityTech/nemesis/pkg/utils"
	"github.com/golang/glog"
)

var (
	flagReportEnableStdout  = flag.Bool("reports.stdout.enable", utils.GetEnvBool("NEMESIS_ENABLE_STDOUT"), "Enable outputting report via stdout")
	flagReportEnablePubsub  = flag.Bool("reports.pubsub.enable", utils.GetEnvBool("NEMESIS_ENABLE_PUBSUB"), "Enable outputting report via Google Pub/Sub")
	flagReportPubsubProject = flag.String("reports.pubsub.project", utils.GetEnv("NEMESIS_PUBSUB_PROJECT", ""), "Indicate which GCP project to output Pub/Sub reports to")
	flagReportPubsubTopic   = flag.String("reports.pubsub.topic", utils.GetEnv("NEMESIS_PUBSUB_TOPIC", "nemesis"), "Indicate which topic to output Pub/Sub reports to")
)

// Audit is a runner that encapsulates the logic of an audit against GCP resources
type Audit struct {
	c         *client.Client
	reports   []report.Report
	reporters []report.Reporter
}

// NewAudit returns a new Audit runner
func NewAudit() *Audit {
	a := new(Audit)
	a.reports = []report.Report{}
	a.reporters = []report.Reporter{}
	return a
}

// Setup configures an Audit runner and sets up audit resources
func (a *Audit) Setup() {
	a.c = client.New()

	a.setupReporters()

	if err := a.c.GetProjects(); err != nil {
		glog.Fatalf("Failed to retrieve project resources: %v", err)
	}

	if err := a.c.GetIamResources(); err != nil {
		glog.Fatalf("Failed to retrieve iam resources: %v", err)
	}

	if err := a.c.GetComputeResources(); err != nil {
		glog.Fatalf("Failed to retrieve compute resources: %v", err)
	}

	if err := a.c.GetLoggingResources(); err != nil {
		glog.Fatalf("Failed to retrieve logging resources: %v", err)
	}

	if err := a.c.GetNetworkResources(); err != nil {
		glog.Fatalf("Failed to retrieve network resources: %v", err)
	}

	if err := a.c.GetContainerResources(); err != nil {
		glog.Fatalf("Failed to retrieve container resources: %v", err)
	}

	if err := a.c.GetStorageResources(); err != nil {
		glog.Fatalf("Failed to retrieve storage resources: %v", err)
	}
}

func (a *Audit) setupReporters() {
	// If pubsub client is required, create it here
	if *flagReportEnablePubsub {

		if *flagReportPubsubProject == "" {
			glog.Fatal("PubSub project not specified")

		}
		if *flagReportPubsubTopic == "" {
			glog.Fatal("PubSub topic not specified")
		}

		// Create the pubsub client
		a.reporters = append(a.reporters, report.NewPubSubReporter(*flagReportPubsubProject, *flagReportPubsubTopic))
	}

	// Setup stdout
	if *flagReportEnableStdout {
		a.reporters = append(a.reporters, report.NewStdOutReporter())
	}
}

// Execute performs the configured audits concurrently to completion
func (a *Audit) Execute() {

	// Setup goroutines for each set of reports we need to collect
	// TODO - how to make this list dynamic?
	generators := []func() (reports []report.Report, err error){
		a.c.GenerateComputeMetadataReports,
		a.c.GenerateComputeInstanceReports,
		a.c.GenerateLoggingReports,
		a.c.GenerateComputeNetworkReports,
		a.c.GenerateComputeSubnetworkReports,
		a.c.GenerateComputeFirewallRuleReports,
		a.c.GenerateComputeAddressReports,
		a.c.GenerateIAMPolicyReports,
		a.c.GenerateStorageBucketReports,
		a.c.GenerateContainerClusterReports,
		a.c.GenerateContainerNodePoolReports,
	}

	for _, f := range generators {
		reports, err := f()
		if err != nil {
			glog.Fatalf("Failed to generate reports: %v", err)
		}
		a.reports = append(a.reports, reports...)
	}

}

// Report exports the configured reports to their final destination
func (a *Audit) Report() {

	// Push metrics
	if err := a.c.PushMetrics(); err != nil {
		glog.Fatalf("Failed to push metrics: %v", err)
	}

	// Push outputs
	for _, r := range a.reporters {
		err := r.Publish(a.reports)
		if err != nil {
			glog.Fatalf("Failed to publish reports: %v", err)
		}
	}
}
