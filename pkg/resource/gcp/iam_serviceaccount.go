package gcp

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang/glog"
	iam "google.golang.org/api/iam/v1"
)

// IamServiceAccountResource is resource for testing IAM Service Accounts in GCP
type IamServiceAccountResource struct {
	s    *iam.ServiceAccount
	Keys []*iam.ServiceAccountKey
}

// NewIamServiceAccountResource returns a new IamServiceAccountResource
func NewIamServiceAccountResource(s *iam.ServiceAccount) *IamServiceAccountResource {
	r := new(IamServiceAccountResource)
	r.s = s
	r.Keys = []*iam.ServiceAccountKey{}
	return r
}

// Email returns the email address of the service account
func (r *IamServiceAccountResource) Email() string {
	return r.s.Email
}

// Marshal returns the underlying resource's JSON representation
func (r *IamServiceAccountResource) Marshal() ([]byte, error) {
	return json.Marshal(&r.s)
}

// HasUserManagedKeys returns whether a service account has user-managed keys
func (r *IamServiceAccountResource) HasUserManagedKeys() bool {
	return len(r.Keys) != 0
}

// HasKeysNeedingRotation returns an error when the service account has keys older than the allowed time
func (r *IamServiceAccountResource) HasKeysNeedingRotation() (err error) {

	var errBuilder strings.Builder

	for _, k := range r.Keys {
		t, err := time.Parse(time.RFC3339, k.ValidAfterTime)
		if err != nil {
			glog.Fatalf("Failed to parse timestamp when checking keys: %v", err)
		}

		if t.Sub(time.Now()).Hours() > float64(saKeyExpirationTime*24) {
			errBuilder.WriteString(fmt.Sprintf("%v has key older than %v days. ", k.Name, saKeyExpirationTime))
		}

	}

	errString := errBuilder.String()
	if errString != "" {
		err = errors.New(errString)
	} else {
		err = nil
	}

	return
}
