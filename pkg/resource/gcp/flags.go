package gcp

import (
	"flag"
	"strconv"
	"strings"

	"github.com/UnityTech/nemesis/pkg/utils"
	"github.com/golang/glog"
)

var (
	// Default Values
	defaultOAuthMinimalScopes = "https://www.googleapis.com/auth/devstorage.read_only,https://www.googleapis.com/auth/logging.write,https://www.googleapis.com/auth/monitoring,https://www.googleapis.com/auth/servicecontrol,https://www.googleapis.com/auth/service.management.readonly,https://www.googleapis.com/auth/trace.append"

	// Values
	userDomains         = []string{}
	minimalOAuthScopes  = []string{}
	saKeyExpirationTime = 0

	// Flags
	flagContainerMinimalOAuthScopes = flag.String("container.oauth-scopes", utils.GetEnv("NEMESIS_CONTAINER_OAUTHSCOPES", defaultOAuthMinimalScopes), "A comma-seperated list of OAuth scopes to allow for GKE clusters")
	flagUserDomains                 = flag.String("iam.user-domains", utils.GetEnv("NEMESIS_IAM_USERDOMAINS", ""), "A comma-separated list of domains to allow users from")
	flagSAKeyExpirationTime         = flag.String("iam.sa-key-expiration-time", utils.GetEnv("NEMESIS_IAM_SA_KEY_EXPIRATION_TIME", "90"), "The time in days to allow service account keys to live before being rotated")
)

func init() {
	var err error

	minimalOAuthScopes = strings.Split(*flagContainerMinimalOAuthScopes, ",")
	userDomains = strings.Split(*flagUserDomains, ",")
	saKeyExpirationTime, err = strconv.Atoi(*flagSAKeyExpirationTime)

	if err != nil {
		glog.Fatalf("Failed to convert SA key expiration time to an integer: %v", err)
	}
}
