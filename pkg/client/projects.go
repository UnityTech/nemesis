package client

import (
	"fmt"

	"github.com/Unity-Technologies/nemesis/pkg/resource/gcp"
	"github.com/Unity-Technologies/nemesis/pkg/utils"
	"github.com/golang/glog"
)

// GetProjects gathers the list of projects and active API resources for the project
func (c *Client) GetProjects() error {

	if *flagProjectFilter == "" {
		glog.Exitf("No project filter was provided. Either specify --project.filter or set NEMESIS_PROJECT_FILTER to the appropriate regex (e.g. my-cool-projects-*)")
	}

	defer utils.Elapsed("GetProjects")()

	// Get list of all projects.
	// Additionally we must make sure that the project is ACTIVE. Any other state will return errors
	projectFilter := fmt.Sprintf("%v AND lifecycleState=ACTIVE", *flagProjectFilter)
	projects := listAllProjects(projectFilter, c.cloudResourceClient)

	// Return an error that we retrieved no projects
	if len(projects) == 0 {
		return fmt.Errorf("No projects found when matching against '%v'", projectFilter)
	}

	// Create a short-lived goroutine for retrieving project services
	servicesWorker := func(workerID int, projectIDs <-chan string, results chan<- serviceCallResult) {
		id := <-projectIDs
		projectID := fmt.Sprintf("projects/%v", id)

		servicesList, err := c.serviceusageClient.Services.List(projectID).Filter("state:ENABLED").Do()
		if err != nil {
			glog.Fatalf("Failed to retrieve list of services for project %v: %v", projectID, err)
		}

		projectServices := []*gcp.ServiceAPIResource{}
		for _, s := range servicesList.Services {
			projectServices = append(projectServices, gcp.NewServiceAPIResource(s))
		}

		res := serviceCallResult{ProjectID: id, Services: projectServices}

		results <- res
	}

	// Setup worker pool
	projectIDs := make(chan string, len(projects))
	results := make(chan serviceCallResult, len(projects))
	numWorkers := len(projects)
	for w := 0; w < numWorkers; w++ {
		go servicesWorker(w, projectIDs, results)
	}

	// Feed the workers and collect the projects for reuse
	for _, p := range projects {
		projectIDs <- p.ProjectId
		c.resourceprojects = append(c.resourceprojects, p)
	}
	close(projectIDs)

	// Collect the results
	for i := 0; i < len(projects); i++ {
		res := <-results
		c.services[res.ProjectID] = res.Services
	}

	return nil
}

type serviceCallResult struct {
	ProjectID string
	Services  []*gcp.ServiceAPIResource
}
