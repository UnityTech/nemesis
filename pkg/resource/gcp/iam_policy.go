package gcp

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	cloudresourcemanager "google.golang.org/api/cloudresourcemanager/v1"
)

const (
	// Rules to check
	editorRole              = "roles/editor"
	serviceAccountUserRole  = "roles/iam.serviceAccountUser"
	serviceAccountAdminRole = "roles/iam.serviceAccountAdmin"
	kmsAdminRole            = "roles/cloudkms.admin"
	kmsRoleMatcher          = "roles/cloudkms."
)

var (
	// Cloud Audit log types
	logTypes = []string{"ADMIN_READ", "DATA_READ", "DATA_WRITE"}
)

// Helper functions for identifying various types of users or roles
func isGCPAccount(member string) bool {
	return strings.Contains(member, "developer.gserviceaccount.com") || strings.Contains(member, "appspot.gserviceaccount.com")
}

func isAdminRole(role string) bool {
	return strings.Contains(role, "admin") || strings.Contains(role, "owner") || strings.Contains(role, "editor")
}

func isIAMUserMember(member string) bool {
	return strings.Contains(member, "user:") || strings.Contains(member, "group:") || strings.Contains(member, "domain:")
}

// IamPolicyResource is resource for testing IAM Policies in GCP
type IamPolicyResource struct {
	p *cloudresourcemanager.Policy
}

// NewIamPolicyResource returns a new IamPolicyResource
func NewIamPolicyResource(p *cloudresourcemanager.Policy) *IamPolicyResource {
	r := new(IamPolicyResource)
	r.p = p
	return r
}

// Marshal returns the underlying resource's JSON representation
func (r *IamPolicyResource) Marshal() ([]byte, error) {
	return json.Marshal(&r.p)
}

// PolicyViolatesUserDomainWhitelist returns whether the policy contains a user or domain that is not part of the domain whitelist
func (r *IamPolicyResource) PolicyViolatesUserDomainWhitelist() (err error) {

	var errBuilder strings.Builder
	for _, b := range r.p.Bindings {
		for _, member := range b.Members {
			if isIAMUserMember(member) {
				for _, domain := range userDomains {
					if !strings.Contains(member, domain) {
						errBuilder.WriteString(fmt.Sprintf("%v is not allowed by your domain whitelist. ", member))
					}
				}
			}
		}
	}

	// If we collected errors, report a failure
	errString := errBuilder.String()
	if errString != "" {
		err = errors.New(errString)
	} else {
		err = nil
	}

	return
}

// MemberHasAdminRole returns whether a given member has an admin role
func (r *IamPolicyResource) MemberHasAdminRole(member string) (err error) {

	for _, b := range r.p.Bindings {
		for _, m := range b.Members {
			if member == m {

				// Found the member, now check if it has an admin role
				if isAdminRole(b.Role) {

					// Allow the default compute and appengine service accounts
					// to have "editor" role
					if isGCPAccount(member) && b.Role == editorRole {
						return
					}

					err = fmt.Errorf("Member has admin role %v", b.Role)
				}
				break
			}
		}
	}

	return
}

// PolicyAllowsIAMUserServiceAccountUserRole checks whether the policy allows non-service account
// users to impersonate a service account (privelage escalation)
func (r *IamPolicyResource) PolicyAllowsIAMUserServiceAccountUserRole() (err error) {

	var errBuilder strings.Builder

	for _, b := range r.p.Bindings {
		if b.Role == serviceAccountUserRole {
			for _, member := range b.Members {
				if isIAMUserMember(member) {
					errBuilder.WriteString(fmt.Sprintf("%v has Service Account User role. ", member))
				}
			}
			break
		}
	}

	errString := errBuilder.String()
	if errString != "" {
		err = errors.New(errString)
	}

	return
}

func (r *IamPolicyResource) findMembersWithOverlappingRoles(roleA string, roleB string) []string {

	aMembers := []string{}
	bMembers := []string{}

	for _, b := range r.p.Bindings {

		// Check for members that have the A role. If we don't that role,
		// then there's nothing to check
		if b.Role == roleA {

			// Now check for a binding with the user role
			for _, bb := range r.p.Bindings {

				// If we find a binding, then we need to check for overlap.
				if bb.Role == roleB {
					aMembers = b.Members
					bMembers = bb.Members
				}
				break
			}
			break
		}
	}

	overlap := []string{}

	// Now compare memberships for overlap
	for _, m := range aMembers {
		for _, mm := range bMembers {
			if m == mm {
				overlap = append(overlap, m)
			}
		}
	}

	return overlap
}

// PolicyViolatesServiceAccountSeparationoOfDuties returns whether the policy allows for IAM users
// to both administrate and impersonate service accounts
func (r *IamPolicyResource) PolicyViolatesServiceAccountSeparationoOfDuties() (err error) {

	// We should report errors when we see a member that has both roles:
	//  -- roles/iam.serviceAccountUser
	//  -- roles/iam.serviceAccountAdmin
	overlap := r.findMembersWithOverlappingRoles(serviceAccountAdminRole, serviceAccountUserRole)

	var errBuilder strings.Builder

	// Now compare memberships. If there is overlap, report these as errors
	for _, m := range overlap {
		errBuilder.WriteString(fmt.Sprintf("%v can both administrate and impersonate service accounts. ", m))
	}

	errString := errBuilder.String()
	if errString != "" {
		err = errors.New(errString)
	}

	return
}

// PolicyViolatesKMSSeparationoOfDuties returns whether the policy allows for KMS users
// to both administrate keyrings and encrypt/decrypt with keys
func (r *IamPolicyResource) PolicyViolatesKMSSeparationoOfDuties() (err error) {

	// We should report errors when we see a member that has the KMS admin role and a non-admin role:
	//  -- roles/cloudkms.admin
	//  -- roles/cloudkms.*

	kmsRolesDefined := []string{}

	for _, b := range r.p.Bindings {

		// If we have no admin role, then there's nothing to check
		if b.Role == kmsAdminRole {
			for _, bb := range r.p.Bindings {

				if bb.Role != kmsAdminRole && strings.Contains(bb.Role, kmsRoleMatcher) {
					kmsRolesDefined = append(kmsRolesDefined, bb.Role)
				}
			}

			break
		}
	}

	var errBuilder strings.Builder

	// Now check each for overlap
	for _, role := range kmsRolesDefined {
		overlap := r.findMembersWithOverlappingRoles(kmsAdminRole, role)
		for _, member := range overlap {
			errBuilder.WriteString(fmt.Sprintf("%v can both administrate and perform actions with %v. ", member, role))
		}
	}

	errString := errBuilder.String()
	if errString != "" {
		err = errors.New(errString)
	}

	return
}

// PolicyConfiguresAuditLogging returns whether the IAM policy defines Cloud Audit logging
func (r *IamPolicyResource) PolicyConfiguresAuditLogging() error {

	// Do we even have an auditConfig?
	if r.p.AuditConfigs == nil {
		return errors.New("Policy does not define auditConfigs")
	}

	if r.p.AuditConfigs[0].Service != "allServices" {
		return errors.New("allServices is not the default audit config policy")
	}

	if r.p.AuditConfigs[0].AuditLogConfigs == nil {
		return errors.New("Policy does not define auditLogConfigs")
	}

	// We must have the required number of audit log config types
	if len(r.p.AuditConfigs[0].AuditLogConfigs) != len(logTypes) {
		return errors.New("Policy does not define all required log types (requires ADMIN_READ, DATA_READ, DATA_WRITE)")
	}

	for _, cfg := range r.p.AuditConfigs[0].AuditLogConfigs {
		found := false
		for _, typ := range logTypes {
			if cfg.LogType == typ {
				found = true
				break
			}
		}

		if !found {
			return errors.New("Policy has an unrecognized auditLogConfig type")
		}
	}

	return nil
}

// PolicyDoesNotHaveAuditLogExceptions returns whether the IAM policy allows for exceptions to audit logging
func (r *IamPolicyResource) PolicyDoesNotHaveAuditLogExceptions() error {

	// Do we even have an auditConfig?
	if r.p.AuditConfigs == nil {
		return errors.New("Policy does not define auditConfigs")
	}

	if r.p.AuditConfigs[0].AuditLogConfigs == nil {
		return errors.New("Policy does not define auditLogConfigs")
	}

	var errBuilder strings.Builder

	for _, cfg := range r.p.AuditConfigs[0].AuditLogConfigs {
		if len(cfg.ExemptedMembers) != 0 {
			errBuilder.WriteString(fmt.Sprintf("%s has the following exceptions: ", cfg.LogType))
			for _, exempt := range cfg.ExemptedMembers {
				errBuilder.WriteString(exempt)
				errBuilder.WriteString(",")
			}
			errBuilder.WriteString(". ")
		}
	}

	errString := errBuilder.String()

	if len(errString) != 0 {
		return errors.New(errString)
	}

	return nil
}
