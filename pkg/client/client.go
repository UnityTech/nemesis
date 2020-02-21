package client

import (
	"github.com/UnityTech/nemesis/pkg/resource/gcp"

	"github.com/golang/glog"

	"context"

	logging "cloud.google.com/go/logging/apiv2"
	push "github.com/prometheus/client_golang/prometheus/push"
	cloudresourcemanager "google.golang.org/api/cloudresourcemanager/v1"
	compute "google.golang.org/api/compute/v1"
	container "google.golang.org/api/container/v1"
	iam "google.golang.org/api/iam/v1"
	serviceusage "google.golang.org/api/serviceusage/v1"
	storage "google.golang.org/api/storage/v1"
)

// Client is the client used for auditing Google Cloud Compute Engine resources
type Client struct {

	// API clients
	computeClient       *compute.Service
	cloudResourceClient *cloudresourcemanager.Service
	storageClient       *storage.Service
	containerClient     *container.Service
	serviceusageClient  *serviceusage.Service
	iamClient           *iam.Service
	logConfigClient     *logging.ConfigClient
	logMetricClient     *logging.MetricsClient

	// Root project
	resourceprojects []*cloudresourcemanager.Project

	// Resources
	services         map[string][]*gcp.ServiceAPIResource
	computeprojects  []*gcp.ComputeProjectResource
	computeMetadatas map[string]*gcp.ComputeProjectMetadataResource
	buckets          map[string][]*gcp.StorageBucketResource
	instances        map[string][]*gcp.ComputeInstanceResource

	// Container Resources
	clusters  map[string][]*gcp.ContainerClusterResource
	nodepools map[string][]*gcp.ContainerNodePoolResource

	// Compute Network Resources
	networks    map[string][]*gcp.ComputeNetworkResource
	subnetworks map[string][]*gcp.ComputeSubnetworkResource
	firewalls   map[string][]*gcp.ComputeFirewallRuleResource
	addresses   map[string][]*gcp.ComputeAddressResource

	// IAM Resources
	policies        map[string]*gcp.IamPolicyResource
	serviceaccounts map[string][]*gcp.IamServiceAccountResource

	// Logging resources
	logSinks   map[string][]*gcp.LoggingSinkResource
	logMetrics map[string][]*gcp.LoggingMetricResource

	// Metrics pusher
	pusher           *push.Pusher
	metricsArePushed bool
}

// New returns a new wrk conforming to the worker.W interface
func New() *Client {
	var cc *compute.Service
	var crm *cloudresourcemanager.Service
	var cs *storage.Service
	var con *container.Service
	var su *serviceusage.Service
	var i *iam.Service
	var lc *logging.ConfigClient
	var lm *logging.MetricsClient

	c := new(Client)
	ctx := context.Background()

	// Create compute client
	cc, err := compute.NewService(ctx)
	if err != nil {
		glog.Fatalf("Failed to create Google Cloud Engine client: %v", err)
	}

	// Create cloudresourcemanager client
	crm, err = cloudresourcemanager.NewService(ctx)
	if err != nil {
		glog.Fatalf("Failed to create Google Cloud Resource Manager client: %v", err)
	}

	// Create storage client
	cs, err = storage.NewService(ctx)
	if err != nil {
		glog.Fatalf("Failed to create Google Cloud Storage client: %v", err)
	}

	// Create container client
	con, err = container.NewService(ctx)
	if err != nil {
		glog.Fatalf("Failed to create Google Container client: %v", err)
	}

	// Create serviceusage client
	su, err = serviceusage.NewService(ctx)
	if err != nil {
		glog.Fatalf("Failed to create Google Service Usage client: %v", err)
	}

	i, err = iam.NewService(ctx)
	if err != nil {
		glog.Fatalf("Failed to create IAM client: %v", err)
	}

	lc, err = logging.NewConfigClient(ctx)
	if err != nil {
		glog.Fatalf("Failed to create logging config client: %v", err)
	}

	lm, err = logging.NewMetricsClient(ctx)
	if err != nil {
		glog.Fatalf("Failed to create logging metrics client: %v", err)
	}

	c.computeClient = cc
	c.cloudResourceClient = crm
	c.storageClient = cs
	c.containerClient = con
	c.serviceusageClient = su
	c.iamClient = i
	c.logConfigClient = lc
	c.logMetricClient = lm

	// Services
	c.resourceprojects = []*cloudresourcemanager.Project{}
	c.computeprojects = []*gcp.ComputeProjectResource{}

	// Resources
	c.services = make(map[string][]*gcp.ServiceAPIResource, 1)
	c.computeMetadatas = make(map[string]*gcp.ComputeProjectMetadataResource, 1)
	c.buckets = make(map[string][]*gcp.StorageBucketResource, 1)
	c.instances = make(map[string][]*gcp.ComputeInstanceResource, 1)
	c.clusters = make(map[string][]*gcp.ContainerClusterResource, 1)
	c.nodepools = make(map[string][]*gcp.ContainerNodePoolResource, 1)

	// Compute networking resources
	c.networks = make(map[string][]*gcp.ComputeNetworkResource, 1)
	c.subnetworks = make(map[string][]*gcp.ComputeSubnetworkResource, 1)
	c.firewalls = make(map[string][]*gcp.ComputeFirewallRuleResource, 1)
	c.addresses = make(map[string][]*gcp.ComputeAddressResource, 1)

	// IAM resources
	c.policies = make(map[string]*gcp.IamPolicyResource, 1)
	c.serviceaccounts = make(map[string][]*gcp.IamServiceAccountResource, 1)

	// Logging resources
	c.logSinks = make(map[string][]*gcp.LoggingSinkResource, 1)
	c.logMetrics = make(map[string][]*gcp.LoggingMetricResource, 1)

	// Configure metrics
	c.pusher = configureMetrics()

	return c
}
