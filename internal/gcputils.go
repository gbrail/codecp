package internal

import (
	"context"
	"os"

	"golang.org/x/oauth2/google"
)

func GetGCPProject() string {
	if projectID := os.Getenv("GOOGLE_CLOUD_PROJECT"); projectID != "" {
		return projectID
	}

	ctx := context.Background()
	credentials, err := google.FindDefaultCredentials(ctx)
	if err == nil && credentials.ProjectID != "" {
		return credentials.ProjectID
	}

	return ""
}
