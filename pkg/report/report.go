// Package report outlines how reports are formatted and validated
package report

import (
	"encoding/json"

	"github.com/UnityTech/nemesis/pkg/cis"
	"github.com/golang/glog"
)

const (
	// Failed indicates that a resource did not match expected spec
	Failed = "failed"

	// Passed indicates that a resource met the expected spec
	Passed = "passed"
)

// Control is a measurable unit of an audit
type Control struct {
	Title  string `json:"title"`
	Desc   string `json:"desc"`
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

// NewControl returns a new Control with the given title
func NewControl(title string, desc string) Control {
	return Control{
		Title:  title,
		Desc:   desc,
		Status: Failed,
	}
}

// NewCISControl returns a new Control based on the CIS controls with a description
func NewCISControl(recommendationID string, desc string) Control {
	rec, ok := cis.Registry[recommendationID]
	if !ok {
		glog.Fatalf("Couldn't find CIS recommendation with ID '%v'", recommendationID)
	}
	return Control{
		Title:  rec.Format(),
		Desc:   desc,
		Status: Failed,
	}
}

// Passed changes the status of the control from `false` to `true`.
func (c *Control) Passed() {
	c.Status = Passed
}

// Report is a top-level structure for capturing information generated from an audit on a resource
type Report struct {
	Type     string          `json:"type"`
	Title    string          `json:"title"`
	Controls []Control       `json:"controls"`
	Data     json.RawMessage `json:"data"`
}

// NewReport returns a new top-level report with a given title
func NewReport(typ string, title string) Report {
	return Report{
		Type:     typ,
		Title:    title,
		Controls: []Control{},
	}
}

// Status returns whether a report passed all the controls it was assigned.
func (r *Report) Status() string {
	for _, c := range r.Controls {
		if c.Status == Failed {
			return Failed
		}
	}
	return Passed
}

// AddControls appends controls to the report. If we only report failures, then controls that pass are not included in the report
func (r *Report) AddControls(controls ...Control) {
	for _, c := range controls {
		if c.Status == Passed && *flagReportOnlyFailures {
			continue
		}
		r.Controls = append(r.Controls, c)
	}
}
