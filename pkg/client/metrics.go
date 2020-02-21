package client

import (
	"errors"
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	push "github.com/prometheus/client_golang/prometheus/push"
)

var (
	promNamespace = "nemesis"

	// Prometheus metrics
	// Total resources scanned
	totalResourcesCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: promNamespace,
			Name:      "total_resources_scanned",
			Help:      "Total number of resources scanned",
		},
	)
	// Report summaries, reported by type, status, and project
	reportSummary = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: promNamespace,
			Name:      "report_summary",
			Help:      "Report summaries by type, status, and project",
		},
		[]string{"type", "name", "status", "project"},
	)
)

// configureMetrics is a helper function for configuring metrics.
// Since we use a push gateway, we must configure our metrics as a push model
func configureMetrics() *push.Pusher {

	// Only configure metrics collection if enabled
	if *flagMetricsEnabled {

		// Create the prometheus registry. We explicitly declare a registry rather than
		// depend on the default registry
		registry := prometheus.NewRegistry()

		// Register the necessary metrics
		registry.MustRegister(totalResourcesCounter)
		registry.MustRegister(reportSummary)

		// Configure the gateway and return the pusher
		pusher := push.New(*flagMetricsGateway, "nemesis_audit").Gatherer(registry)
		return pusher
	}

	return nil
}

// incrementMetrics is a small helper to consolidate reporting metrics that are reported for all resources
func (c *Client) incrementMetrics(typ string, name string, status string, projectID string) {
	totalResourcesCounter.Inc()
	reportSummary.WithLabelValues(typ, name, status, projectID).Inc()
}

// PushMetrics pushes the collected metrics from this client. Should only be called once.
func (c *Client) PushMetrics() error {

	// Only push metrics if we configured it
	if c.pusher != nil {

		if c.metricsArePushed {
			return errors.New("Metrics were already pushed, make sure client.PushMetrics is only called once")
		}

		if err := c.pusher.Add(); err != nil {
			return fmt.Errorf("Failed to push metrics to gateway: %v", err)
		}

		// Indicate that metrics for the client have already been pushed
		c.metricsArePushed = true
	}

	return nil
}
