package client

import (
	"fmt"

	"github.com/UnityTech/nemesis/pkg/report"
	"github.com/UnityTech/nemesis/pkg/resource/gcp"
	"github.com/UnityTech/nemesis/pkg/utils"
	"github.com/golang/glog"
	"google.golang.org/api/cloudresourcemanager/v1"
)

// GetIamResources gathers the list of IAM resources for the projects
func (c *Client) GetIamResources() error {

	defer utils.Elapsed("GetIamResources")()

	worker := func(projectIDs <-chan string, results chan<- iamCallResult) {

		id := <-projectIDs
		projectID := fmt.Sprintf("projects/%v", id)
		res := iamCallResult{ProjectID: id, Policy: nil, ServiceAccounts: []*gcp.IamServiceAccountResource{}}

		req := cloudresourcemanager.GetIamPolicyRequest{}
		policy, err := c.cloudResourceClient.Projects.GetIamPolicy(id, &req).Do()
		if err != nil {
			glog.Fatalf("Failed to retrieve IAM policy for project %v: %v", id, err)
		}

		res.Policy = gcp.NewIamPolicyResource(policy)

		saList, err := c.iamClient.Projects.ServiceAccounts.List(projectID).Do()
		if err != nil {
			glog.Fatalf("Failed to retrieve service accounts from project %v: %v", id, err)
		}

		for _, a := range saList.Accounts {

			acct := gcp.NewIamServiceAccountResource(a)
			saKeySearch := fmt.Sprintf("%v/serviceAccounts/%v", projectID, a.UniqueId)
			keys, err := c.iamClient.Projects.ServiceAccounts.Keys.List(saKeySearch).KeyTypes("USER_MANAGED").Do()
			if err != nil {
				glog.Fatalf("Failed to retrieve service account keys from project %v: %v", id, err)
			}
			for _, k := range keys.Keys {
				acct.Keys = append(acct.Keys, k)
			}

			res.ServiceAccounts = append(res.ServiceAccounts, acct)
		}

		results <- res
	}

	// Setup worker pool
	projectIDs := make(chan string, len(c.resourceprojects))
	results := make(chan iamCallResult, len(c.resourceprojects))
	numWorkers := len(c.resourceprojects)
	for w := 0; w < numWorkers; w++ {
		go worker(projectIDs, results)
	}

	// Feed the workers and collect the cluster info
	for _, p := range c.resourceprojects {
		projectIDs <- p.ProjectId
	}

	// Collect the info
	for i := 0; i < numWorkers; i++ {
		res := <-results
		c.policies[res.ProjectID] = res.Policy
		c.serviceaccounts[res.ProjectID] = res.ServiceAccounts
	}

	return nil
}

type iamCallResult struct {
	ProjectID       string
	Policy          *gcp.IamPolicyResource
	ServiceAccounts []*gcp.IamServiceAccountResource
}

// GenerateIAMPolicyReports signals the client to process IamPolicyResource's for reports.
func (c *Client) GenerateIAMPolicyReports() (reports []report.Report, err error) {

	reports = []report.Report{}
	typ := "iam_policy"

	for _, p := range c.computeprojects {
		projectID := p.Name()
		policy := c.policies[projectID]
		serviceAccounts := c.serviceaccounts[projectID]

		r := report.NewReport(
			typ,
			fmt.Sprintf("Project %v IAM Policy", projectID),
		)
		r.Data, err = policy.Marshal()
		if err != nil {
			glog.Fatalf("Failed to marshal IAM policy: %v", err)
		}

		// Corporate login credentials should be used
		corpCreds := report.NewCISControl(
			"1.1",
			fmt.Sprintf("Project %v should only allow corporate login credentials", p.Name()),
		)
		if err := policy.PolicyViolatesUserDomainWhitelist(); err != nil {
			corpCreds.Error = err.Error()
		} else {
			corpCreds.Passed()
		}
		r.AddControls(corpCreds)

		for _, sa := range serviceAccounts {

			// Service account keys should be GCP-managed
			saManagedKeys := report.NewCISControl(
				"1.3",
				fmt.Sprintf("%v should not have user-managed keys", sa.Email()),
			)
			if sa.HasUserManagedKeys() {
				saManagedKeys.Error = "Service account has user-managed keys"
			} else {
				saManagedKeys.Passed()
			}
			r.AddControls(saManagedKeys)

			// Service accounts should not have admin privileges
			saAdminRole := report.NewCISControl(
				"1.4",
				fmt.Sprintf("%v should not have admin roles", sa.Email()),
			)
			if err := policy.MemberHasAdminRole(fmt.Sprintf("serviceAccount:%v", sa.Email())); err != nil {
				saAdminRole.Error = err.Error()
			} else {
				saAdminRole.Passed()
			}
			r.AddControls(saAdminRole)

		}

		// IAM Users should not be able to impersonate service accounts at the project level
		saServiceAccountUserRole := report.NewCISControl(
			"1.5",
			fmt.Sprintf("Project %v should not allow project-wide use of Service Account User role", p.Name()),
		)
		if err := policy.PolicyAllowsIAMUserServiceAccountUserRole(); err != nil {
			saServiceAccountUserRole.Error = err.Error()
		} else {
			saServiceAccountUserRole.Passed()
		}
		r.AddControls(saServiceAccountUserRole)

		// Service account keys should be rotated on a regular interval
		for _, sa := range serviceAccounts {
			saKeyExpired := report.NewCISControl(
				"1.6",
				fmt.Sprintf("%v should not have expired keys", sa.Email()),
			)
			if err := sa.HasKeysNeedingRotation(); err != nil {
				saKeyExpired.Error = err.Error()
			} else {
				saKeyExpired.Passed()
			}
			r.AddControls(saKeyExpired)
		}

		// Users should not be allowed to administrate and impersonate service accounts
		saSeperateDuties := report.NewCISControl(
			"1.7",
			fmt.Sprintf("Project %v should have separation of duties with respect to service account usage", p.Name()),
		)
		if err := policy.PolicyViolatesServiceAccountSeparationoOfDuties(); err != nil {
			saSeperateDuties.Error = err.Error()
		} else {
			saSeperateDuties.Passed()
		}
		r.AddControls(saSeperateDuties)

		// Users should not be allowed to administrate and utilize KMS functionality
		kmsSeperateDuties := report.NewCISControl(
			"1.9",
			fmt.Sprintf("Project %v should have separation of duties with respect to KMS usage", p.Name()),
		)
		if err := policy.PolicyViolatesKMSSeparationoOfDuties(); err != nil {
			kmsSeperateDuties.Error = err.Error()
		} else {
			kmsSeperateDuties.Passed()
		}
		r.AddControls(kmsSeperateDuties)

		// Project IAM Policies should define audit configurations
		auditConfig := report.NewCISControl(
			"2.1",
			fmt.Sprintf("Project %v should proper audit logging configurations", p.Name()),
		)
		if err := policy.PolicyConfiguresAuditLogging(); err != nil {
			auditConfig.Error = err.Error()
		} else {
			if err := policy.PolicyDoesNotHaveAuditLogExceptions(); err != nil {
				auditConfig.Error = err.Error()
			} else {
				auditConfig.Passed()
			}
		}
		r.AddControls(auditConfig)

		reports = append(reports, r)
	}

	return
}
