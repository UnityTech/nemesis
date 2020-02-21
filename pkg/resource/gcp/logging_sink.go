package gcp

import (
	loggingpb "google.golang.org/genproto/googleapis/logging/v2"
)

// LoggingSinkResource represents a StackDriver logging sink
type LoggingSinkResource struct {
	s *loggingpb.LogSink
}

// NewLoggingSinkResource returns a new LoggingSinkResource
func NewLoggingSinkResource(sink *loggingpb.LogSink) *LoggingSinkResource {
	r := new(LoggingSinkResource)
	r.s = sink
	return r
}

// ShipsAllLogs indicates whether there is no filter (and thus all logs are shipped)
func (r *LoggingSinkResource) ShipsAllLogs() bool {

	// An empty string indicates that there is no filter - thus all logs
	// that are generated are shipped to the logging sink destination
	return r.s.Filter == ""
}
