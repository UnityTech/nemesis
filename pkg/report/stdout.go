package report

import (
	"encoding/json"
	"fmt"

	"github.com/golang/glog"
)

// StdOutReporter is a reporter that prints audit reports to stdout
type StdOutReporter struct{}

// NewStdOutReporter returns a new StdOutReporter for outputting the findings of an audit
func NewStdOutReporter() *StdOutReporter {
	r := new(StdOutReporter)
	return r
}

// Publish prints a full list of reports to stdout
func (r *StdOutReporter) Publish(reports []Report) error {
	b, err := json.Marshal(&reports)
	if err != nil {
		glog.Fatalf("Failed to render report: %v", err)
	}
	fmt.Println(string(b))
	return nil
}
