package gcp

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	storage "google.golang.org/api/storage/v1"
)

// Helper function for making fake bucket resources
func makeTestBucket(data []byte) *storage.Bucket {
	bucket := new(storage.Bucket)
	_ = json.Unmarshal(data, &bucket)
	return bucket
}

var (
	testValidBucketData = []byte(`
{
	"acl": [
		{
			"bucket": "my-test-bucket",
			"entity": "project-owners-01010101010101",
			"etag": "CAE=",
			"id": "my-test-bucket/project-owners-01010101010101",
			"kind": "storage#bucketAccessControl",
			"projectTeam": {
				"projectNumber": "01010101010101",
				"team": "owners"
			},
			"role": "OWNER",
			"selfLink": "https://www.googleapis.com/storage/v1/b/my-test-bucket/acl/project-owners-01010101010101"
		},
		{
			"bucket": "my-test-bucket",
			"entity": "project-editors-01010101010101",
			"etag": "CAE=",
			"id": "my-test-bucket/project-editors-01010101010101",
			"kind": "storage#bucketAccessControl",
			"projectTeam": {
				"projectNumber": "01010101010101",
				"team": "editors"
			},
			"role": "OWNER",
			"selfLink": "https://www.googleapis.com/storage/v1/b/my-test-bucket/acl/project-editors-01010101010101"
		},
		{
			"bucket": "my-test-bucket",
			"entity": "project-viewers-01010101010101",
			"etag": "CAE=",
			"id": "my-test-bucket/project-viewers-01010101010101",
			"kind": "storage#bucketAccessControl",
			"projectTeam": {
				"projectNumber": "01010101010101",
				"team": "viewers"
			},
			"role": "READER",
			"selfLink": "https://www.googleapis.com/storage/v1/b/my-test-bucket/acl/project-viewers-01010101010101"
		}
		],
		"etag": "CAE=",
		"iamConfiguration": {
			"bucketPolicyOnly": {}
		},
	"id": "my-test-bucket",
	"kind": "storage#bucket",
	"location": "US-CENTRAL1",
	"metageneration": "1",
	"name": "my-test-bucket",
	"projectNumber": "01010101010101",
	"selfLink": "https://www.googleapis.com/storage/v1/b/my-test-bucket",
	"storageClass": "REGIONAL",
	"timeCreated": "2019-01-18T14:14:07.472Z",
	"updated": "2019-01-18T14:14:07.472Z"
}
`)

	testValidBucket = makeTestBucket(testValidBucketData)
)

// Make sure that a bucket resource is created correctly
func TestNewStorageBucketResource(t *testing.T) {

	bucketResource := NewStorageBucketResource(testValidBucket)

	// Make sure the resource is not nil
	assert.NotNil(t, bucketResource)

	// Make sure the underlying datasource is not nil
	assert.NotNil(t, bucketResource.b)
}

func TestStorageBucketResourceName(t *testing.T) {

	bucketResource := NewStorageBucketResource(testValidBucket)

	// Make sure the name checks out
	assert.Equal(t, "my-test-bucket", bucketResource.Name())

}
func TestStorageBucketResourceMarshal(t *testing.T) {

	bucketResource := NewStorageBucketResource(testValidBucket)

	// Marshal the original bucket
	orig, err := json.Marshal(&testValidBucket)

	// Make sure the data returns the same as we put in
	data, err := bucketResource.Marshal()
	assert.Nil(t, err)

	assert.Equal(t, orig, data)
}

func TestStorageBucketResourceAllowAllUsers(t *testing.T) {

	// Assert that the bucket does not contain the `allUsers` entity
	bucketResource := NewStorageBucketResource(testValidBucket)
	exists := bucketResource.AllowAllUsers()
	assert.False(t, exists)

	// TODO - add a bucket with the `allUsers` entity and check that it works

}
func TestStorageBucketResourceAllowAllAuthenticatedUsers(t *testing.T) {

	// Assert that the bucket does not contain the `allAuthenticatedUsers` entity
	bucketResource := NewStorageBucketResource(testValidBucket)
	exists := bucketResource.AllowAllAuthenticatedUsers()
	assert.False(t, exists)

	// TODO - add a bucket with the `allAuthenticatedUsers` entity and check that it works

}
func TestStorageBucketResourceHasBucketPolicyOnlyEnabled(t *testing.T) {

	// Assert that the bucket does not have the BucketPolicyOnly IAM configuration
	bucketResource := NewStorageBucketResource(testValidBucket)
	exists, err := bucketResource.HasBucketPolicyOnlyEnabled()
	assert.Nil(t, err)
	assert.False(t, exists)

	// TODO - add a bucket with the BucketPolicyOnly IAM configuration and test it

}
