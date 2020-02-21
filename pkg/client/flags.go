package client

import (
	"flag"

	"github.com/UnityTech/nemesis/pkg/utils"
)

var (
	// Metrics
	flagMetricsEnabled = flag.Bool("metrics.enabled", utils.GetEnvBool("NEMESIS_METRICS_ENABLED"), "Enable Prometheus metrics")
	flagMetricsGateway = flag.String("metrics.gateway", utils.GetEnv("NEMESIS_METRICS_GATEWAY", "127.0.0.1:9091"), "Prometheus metrics Push Gateway")

	// Projects
	flagProjectFilter = flag.String("project.filter", utils.GetEnv("NEMESIS_PROJECT_FILTER", ""), "REQUIRED - the project filter to perform audits on.")

	// Compute
	flagComputeInstanceNumInterfaces     = flag.Int("compute.instance.num-interfaces", utils.GetEnvInt("NEMESIS_COMPUTE_NUM_NICS", 1), "The number of network interfaces (NIC) that an instance should have")
	flagComputeInstanceAllowNat          = flag.Bool("compute.instance.allow-nat", utils.GetEnvBool("NEMESIS_COMPUTE_ALLOW_NAT"), "Indicate whether instances should be allowed to have external (NAT) IP addresses")
	flagComputeInstanceAllowIPForwarding = flag.Bool("compute.instance.allow-ip-forwarding", utils.GetEnvBool("NEMESIS_COMPUTE_ALLOW_IP_FORWARDING"), "Indicate whether instances should be allowed to perform IP forwarding")
)
