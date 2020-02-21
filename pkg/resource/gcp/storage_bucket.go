package gcp

import (
	"encoding/json"
	"fmt"

	storage "google.golang.org/api/storage/v1"
)

// StorageBucketResource represents a Google Storage bucket resource
type StorageBucketResource struct {
	b *storage.Bucket
}

// NewStorageBucketResource returns a new StorageBucketResource
func NewStorageBucketResource(b *storage.Bucket) *StorageBucketResource {
	r := new(StorageBucketResource)
	r.b = b
	return r
}

// Name returns the bucket's name
func (r *StorageBucketResource) Name() string {
	return r.b.Name
}

// Marshal returns the underlying resource's JSON representation
func (r *StorageBucketResource) Marshal() ([]byte, error) {
	return json.Marshal(&r.b)
}

// AllowAllUsers checks whether a bucket is configured to be world readable
func (r *StorageBucketResource) AllowAllUsers() (result bool) {

	acls := r.b.Acl

	for _, acl := range acls {

		// The `allUsers` entity denotes public access
		if acl.Entity == "allUsers" {
			result = true
			return
		}
	}

	return
}

// AllowAllAuthenticatedUsers checks whether a bucket is configured to be readable by anyone with a Google account
func (r *StorageBucketResource) AllowAllAuthenticatedUsers() (result bool) {

	acls := r.b.Acl

	for _, acl := range acls {

		// The `allAuthenticatedUsers` entity denotes access to any authenticated user to Google
		if acl.Entity == "allAuthenticatedUsers" {
			result = true
			return
		}
	}

	return
}

// HasBucketPolicyOnlyEnabled checks whether a bucket is configured to use permissions across the entire bucket
func (r *StorageBucketResource) HasBucketPolicyOnlyEnabled() (result bool, err error) {

	result = false
	iamConfig := r.b.IamConfiguration

	if iamConfig == nil {
		err = fmt.Errorf("Could not retrieve IAM configuration for gs://%v", r.b.Name)
		return
	}

	// Check if the policy exists. If not, then pass
	if bucketPolicyOnly := iamConfig.BucketPolicyOnly; bucketPolicyOnly != nil {
		// If the policy exists, return whether it is enabled
		result = bucketPolicyOnly.Enabled
	}
	return result, err
}
