package gcp

import (
	loggingpb "google.golang.org/genproto/googleapis/logging/v2"
)

// LoggingMetricResource represents a StackDriver log-based metric
type LoggingMetricResource struct {
	m *loggingpb.LogMetric
}

// NewLoggingMetricResource returns a new LoggingMetricResource
func NewLoggingMetricResource(metric *loggingpb.LogMetric) *LoggingMetricResource {
	r := new(LoggingMetricResource)
	r.m = metric
	return r
}

// Filter returns the filter of the metric
func (r *LoggingMetricResource) Filter() string {
	return r.m.Filter
}

// FilterMatches returns whether the configured filter matches a given string
func (r *LoggingMetricResource) FilterMatches(filter string) bool {
	return r.m.Filter == filter
}
