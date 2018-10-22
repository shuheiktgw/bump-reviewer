package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

const (
	ReviewApprove = "approve"
	ReviewDismiss = "dismiss"
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

// GetPullRequestFiles gets files edited by a PR
func (c *GitHubClient) ListPullRequestsFiles(number int, opt *github.ListOptions) ([]*github.CommitFile, error) {
	cf, res, err := c.Client.PullRequests.ListFiles(context.TODO(), c.Owner, c.Repo, number, opt)

	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("PullRequests.ListFiles returns invalid status: %s", res.Status)
	}

	return cf, nil
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

// GetContent gets the specified file
func (c *GitHubClient) GetContent(path string, opt *github.RepositoryContentGetOptions) (*github.RepositoryContent, []*github.RepositoryContent, error) {
	fc, dc, res, err := c.Client.Repositories.GetContents(context.TODO(), c.Owner, c.Repo, path, opt)

	if err != nil {
		return nil, nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, nil, fmt.Errorf("Repositories.GetContents returns invalid status: %s", res.Status)
	}

	return fc, dc, nil
}

// CreateReview creates a review on a given PR
func (c *GitHubClient) CreateReview(number int, review *github.PullRequestReviewRequest) (*github.PullRequestReview, error) {
	prr, res, err := c.Client.PullRequests.CreateReview(context.TODO(), c.Owner, c.Repo, number, review)

	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("PullRequests.CreateReview returns invalid status: %s", res.Status)
	}

	return prr, nil
}
