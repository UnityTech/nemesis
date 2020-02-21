package client

import (
	"fmt"

	"github.com/UnityTech/nemesis/pkg/report"
	"github.com/UnityTech/nemesis/pkg/resource/gcp"
	"github.com/UnityTech/nemesis/pkg/utils"
	"github.com/golang/glog"
)

// GetStorageResources launches the process retrieving storage buckets and other storage resources
func (c *Client) GetStorageResources() error {

	defer utils.Elapsed("GetStorageResources")()

	worker := func(projectIDs <-chan string, results chan<- storageCallResult) {

		id := <-projectIDs
		res := storageCallResult{ProjectID: id, Buckets: []*gcp.StorageBucketResource{}}

		// Get the project's buckets
		bucketList, err := c.storageClient.Buckets.List(id).Do()
		if err != nil {
			glog.Fatalf("Error retrieving project %v's bucket list: %v", id, err)
		}

		for _, b := range bucketList.Items {

			// Get the ACLs for the bucket, as they are not included by default in the bucket list call
			acls, err := c.storageClient.BucketAccessControls.List(b.Name).Do()
			if err != nil {

				continue
				/*
					// The call above will throw a 400 error if Bucket Policy Only is enabled
					if strings.Contains(
						err.Error(),
						"googleapi: Error 400: Cannot get legacy ACLs for a bucket that has enabled Bucket Policy Only",
					) {
						continue
					}
				*/

				// Otherwise, we hit a real error
				//glog.Fatalf("Error retrieving bucket %v's ACLs: %v", b.Name, err)
			}

			// Store the ACLs with the bucket
			b.Acl = acls.Items

			// Append a new bucket resource
			res.Buckets = append(res.Buckets, gcp.NewStorageBucketResource(b))
		}

		results <- res
	}

	// Setup worker pool
	projectIDs := make(chan string, len(c.resourceprojects))
	results := make(chan storageCallResult, len(c.resourceprojects))
	numWorkers := len(c.resourceprojects)
	for w := 0; w < numWorkers; w++ {
		go worker(projectIDs, results)
	}

	// Feed the workers and collect the storage info
	for _, p := range c.resourceprojects {
		projectIDs <- p.ProjectId
	}

	// Collect the info
	for i := 0; i < numWorkers; i++ {
		res := <-results
		c.buckets[res.ProjectID] = res.Buckets
	}

	return nil
}

type storageCallResult struct {
	ProjectID string
	Buckets   []*gcp.StorageBucketResource
}

// GenerateStorageBucketReports signals the client to process ComputeStorageBucket's for reports.
// If there are keys configured for buckets in the configuration, no reports will be created.
func (c *Client) GenerateStorageBucketReports() (reports []report.Report, err error) {

	reports = []report.Report{}
	typ := "storage_bucket"

	for _, p := range c.computeprojects {

		projectID := p.Name()
		projectBuckets := c.buckets[projectID]

		for _, b := range projectBuckets {
			r := report.NewReport(typ, fmt.Sprintf("Project %v Storage Bucket %v", projectID, b.Name()))
			if r.Data, err = b.Marshal(); err != nil {
				glog.Fatalf("Failed to marshal storage bucket: %v", err)
			}

			allUsersControl := report.NewCISControl(
				"5.1",
				"Bucket ACL should not include entity 'allUsers'",
			)

			if !b.AllowAllUsers() {
				allUsersControl.Passed()
			} else {
				allUsersControl.Error = "Bucket ACL includes entity 'allUsers'"
			}

			// Add the `allAuthenticatedUsers` entity control if the spec says that allAuthenticatedUsers == false
			allAuthenticatedUsersControl := report.NewCISControl(
				"5.1",
				"Bucket ACL should not include entity 'allAuthenticatedUsers'",
			)

			if !b.AllowAllAuthenticatedUsers() {
				allAuthenticatedUsersControl.Passed()
			} else {
				allAuthenticatedUsersControl.Error = "Bucket ACL includes entity 'allAuthenticatedUsers'"
			}

			r.AddControls(allUsersControl, allAuthenticatedUsersControl)

			// Add the bucket report to the final list of bucket reports
			reports = append(reports, r)
			c.incrementMetrics(typ, b.Name(), r.Status(), projectID)
		}

	}

	return
}
