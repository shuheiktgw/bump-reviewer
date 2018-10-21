package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

// GitHubClient is a clint to interact with Github API
type GitHubClient struct {
	Owner, Repo string
	Client      *github.Client
}

// NewGitHubClient creates and initializes a new GitHubClient
func NewGitHubClient(owner, repo, token string) *GitHubClient {
	ts := oauth2.StaticTokenSource(&oauth2.Token{
		AccessToken: token,
	})
	tc := oauth2.NewClient(context.Background(), ts)

	client := github.NewClient(tc)

	return &GitHubClient{
		Owner:  owner,
		Repo:   repo,
		Client: client,
	}
}

// CompareCommits gets diffs between base and head
func (c *GitHubClient) CompareCommits(base, head string) (*github.CommitsComparison, error) {
	cc, res, err := c.Client.Repositories.CompareCommits(context.TODO(), c.Owner, c.Repo, base, head)

	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Repositories.CompareCommits returns invalid status: %s", res.Status)
	}

	return cc, nil
}

// GetLatestRelease gets the latest release of the repository
func (c *GitHubClient) GetLatestRelease() (*github.RepositoryRelease, error) {
	rr, res, err := c.Client.Repositories.GetLatestRelease(context.TODO(), c.Owner, c.Repo)

	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Repositories.GetLatestRelease returns invalid status: %s", res.Status)
	}

	return rr, nil
}
