package client

import (
	"context"
	"net/http"

	"github.com/golang/glog"
	"golang.org/x/oauth2/google"
)

// getOAuthClient configures an http.client capable of authenticating to Google APIs
func getOAuthClient(ctx context.Context, scope string) *http.Client {

	// Configure client with OAuth
	oauthClient, err := google.DefaultClient(ctx, scope)
	if err != nil {
		glog.Fatalf("Failed to create OAuth client: %v", err)
	}

	return oauthClient
}
