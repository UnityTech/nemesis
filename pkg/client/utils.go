package client

import (
	"fmt"
	"strings"

	"github.com/golang/glog"
	cloudresourcemanager "google.golang.org/api/cloudresourcemanager/v1"
)

// listAllProjects returns the list of all Projects visible to the authenticated client
func listAllProjects(filter string, client *cloudresourcemanager.Service) []*cloudresourcemanager.Project {

	var projects []*cloudresourcemanager.Project

	projectListCall, err := client.Projects.List().Filter(fmt.Sprintf("name:%v", filter)).Do()
	if err != nil {
		glog.Fatalf("Error retreiving projects: %v", err)
	}

	projects = append(projects, projectListCall.Projects...)

	for projectListCall.NextPageToken != "" {
		projectListCall, err := client.Projects.List().PageToken(projectListCall.NextPageToken).Do()
		if err != nil {
			glog.Fatalf("Error retreiving projects: %v", err)
		}
		projects = append(projects, projectListCall.Projects...)
	}

	return projects
}

func (c *Client) isServiceEnabled(projectID, servicename string) bool {
	for _, api := range c.services[projectID] {
		if strings.Contains(api.Name(), servicename) {
			return true
		}
	}
	return false
}
