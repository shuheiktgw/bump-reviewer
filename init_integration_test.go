// +build integration

package main

import "os"

var (
	integrationGitHubClient *GitHubClient
	integrationGitHubToken  = os.Getenv("BUMP_REVIEWER_INTEGRATION_GITHUB_TOKEN")
	integrationGitHubOwner  = os.Getenv("BUMP_REVIEWER_INTEGRATION_GITHUB_OWNER")
	integrationGitHubRepo   = os.Getenv("BUMP_REVIEWER_INTEGRATION_GITHUB_REPO")
)

func init() {
	integrationGitHubClient = NewGitHubClient(integrationGitHubOwner, integrationGitHubRepo, integrationGitHubToken)
}
