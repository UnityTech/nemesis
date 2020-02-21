// Package cis is a schema for organizing CIS controls for Google Cloud
package cis

import (
	"encoding/json"
	"fmt"
)

// Recommendation is a CIS recommendation for GCP
type Recommendation struct {
	// The name of the CIS recommendation
	Name string `json:"name"`

	// Indicates whether compliance with the recommendation should
	// be attributeable to the overall compliance of the relevant resource
	Scored bool `json:"scored"`

	// The CIS identifier for the recommendation. They are formatted in a major-minor
	// string (E.g. "1.12")
	CisID string `json:"cisId"`

	// The CIS level for the recommendation.
	Level int `json:"level"`
}

var (
	// Registry is the registry of CIS recommendations
	Registry = make(map[string]Recommendation, 1)
)

// Marshal returns the JSON formatted bytes for a recommendation
func (r *Recommendation) Marshal() ([]byte, error) {
	return json.Marshal(&r)
}

// Format returns the fully formatted CIS descriptive name
func (r *Recommendation) Format() string {
	score := "Scored"
	if !r.Scored {
		score = "Not Scored"
	}
	return fmt.Sprintf("CIS %v - %v (%v)", r.CisID, r.Name, score)
}

func init() {
	// IAM controls
	Registry[iam1.CisID] = iam1
	Registry[iam2.CisID] = iam2
	Registry[iam3.CisID] = iam3
	Registry[iam4.CisID] = iam4
	Registry[iam5.CisID] = iam5
	Registry[iam6.CisID] = iam6
	Registry[iam7.CisID] = iam7
	Registry[iam8.CisID] = iam8
	Registry[iam9.CisID] = iam9
	Registry[iam10.CisID] = iam10
	Registry[iam11.CisID] = iam11
	Registry[iam12.CisID] = iam12
	Registry[iam13.CisID] = iam13

	// Logging & Monitoring
	Registry[logmon1.CisID] = logmon1
	Registry[logmon2.CisID] = logmon2
	Registry[logmon3.CisID] = logmon3
	Registry[logmon4.CisID] = logmon4
	Registry[logmon5.CisID] = logmon5
	Registry[logmon6.CisID] = logmon6
	Registry[logmon7.CisID] = logmon7
	Registry[logmon8.CisID] = logmon8
	Registry[logmon9.CisID] = logmon9
	Registry[logmon10.CisID] = logmon10
	Registry[logmon11.CisID] = logmon11

	// Networking
	Registry[network1.CisID] = network1
	Registry[network2.CisID] = network2
	Registry[network3.CisID] = network3
	Registry[network4.CisID] = network4
	Registry[network5.CisID] = network5
	Registry[network6.CisID] = network6
	Registry[network7.CisID] = network7
	Registry[network8.CisID] = network8
	Registry[network9.CisID] = network9

	// VM & Compute
	Registry[compute1.CisID] = compute1
	Registry[compute2.CisID] = compute2
	Registry[compute3.CisID] = compute3
	Registry[compute4.CisID] = compute4
	Registry[compute5.CisID] = compute5
	Registry[compute6.CisID] = compute6

	// GCS Storage
	Registry[storage1.CisID] = storage1
	Registry[storage2.CisID] = storage2
	Registry[storage3.CisID] = storage3

	// SQL
	Registry[sql1.CisID] = sql1
	Registry[sql2.CisID] = sql2
	Registry[sql3.CisID] = sql3
	Registry[sql4.CisID] = sql4

	// Kubernetes Engine
	Registry[gke1.CisID] = gke1
	Registry[gke2.CisID] = gke2
	Registry[gke3.CisID] = gke3
	Registry[gke4.CisID] = gke4
	Registry[gke5.CisID] = gke5
	Registry[gke6.CisID] = gke6
	Registry[gke7.CisID] = gke7
	Registry[gke8.CisID] = gke8
	Registry[gke9.CisID] = gke9
	Registry[gke10.CisID] = gke10
	Registry[gke11.CisID] = gke11
	Registry[gke12.CisID] = gke12
	Registry[gke13.CisID] = gke13
	Registry[gke14.CisID] = gke14
	Registry[gke15.CisID] = gke15
	Registry[gke16.CisID] = gke16
	Registry[gke17.CisID] = gke17
	Registry[gke18.CisID] = gke18
}
