// +build integration

package main

import (
	"fmt"
	"testing"

	"github.com/google/go-github/github"
)

func TestGitHubClient_Integration_CompareCommits(t *testing.T) {
	base := "master"
	head := "pull/1/head"

	cc, err := integrationGitHubClient.CompareCommits(base, head)

	if err != nil {
		t.Fatalf("GitHubClient.CompareCommits returns unexpected error: %s", err)
	}

	if got, want := len(cc.Files), 1; got != want {
		t.Fatalf("GitHubClient.CompareCommits returns unexpected number of files: want: %d, got: %d", want, got)
	}

	if got, want := *cc.Files[0].Filename, "lib/bump-reviewer-test/version.rb"; got != want {
		t.Fatalf("GitHubClient.CompareCommits returns unexpected file: want: %s, got: %s", want, got)
	}
}

func TestGitHubClient_Integration_ListPullRequestsFiles(t *testing.T) {
	cf, err := integrationGitHubClient.ListPullRequestsFiles(1, nil)

	if err != nil {
		t.Fatalf("GitHubClient.ListPullRequestsFiles returns unexpected error: %s", err)
	}

	if got, want := len(cf), 1; got != want {
		t.Fatalf("GitHubClient.ListPullRequestsFiles returns unexpected number of files: want: %d, got: %d", want, got)
	}

	if got, want := *cf[0].Filename, "lib/bump-reviewer-test/version.rb"; got != want {
		t.Fatalf("GitHubClient.CompareCommits returns unexpected file: want: %s, got: %s", want, got)
	}
}

func TestGitHubClient_Integration_GetLatestRelease(t *testing.T) {
	rr, err := integrationGitHubClient.GetLatestRelease()

	if err != nil {
		t.Fatalf("GitHubClient.GetLatestRelease returns unexpected error: %s", err)
	}

	if got, want := *rr.TagName, "v0.0.1"; got != want {
		t.Fatalf("GitHubClient.GetLatestRelease returns unexpected TagName: want: %s, got: %s", want, got)
	}

	if got, want := *rr.Name, "Release v0.0.1"; got != want {
		t.Fatalf("GitHubClient.GetLatestRelease returns unexpected Name: want: %s, got: %s", want, got)
	}
}

func TestGitHubClient_Integration_GetContent(t *testing.T) {
	path := fmt.Sprintf("lib/%s/version.rb", integrationGitHubRepo)
	opt := github.RepositoryContentGetOptions{Ref: "pull/1/head"}
	fc, _, err := integrationGitHubClient.GetContent(path, &opt)

	if err != nil {
		t.Fatalf("GitHubClient.GetContent returns unexpected error: %s", err)
	}

	if got, want := *fc.Name, "version.rb"; got != want {
		t.Fatalf("GitHubClient.GetContent returns unexpected Name: want: %s, got: %s", want, got)
	}
}

func TestGitHubClient_Integration_CreateReview(t *testing.T) {
	number := 3
	review := github.PullRequestReviewRequest{Body: github.String("LGTM"), Event: github.String(ReviewApprove)}

	prr, err := integrationGitHubClient.CreateReview(number, &review)

	if err != nil {
		t.Fatalf("GitHubClient.CreateReview returns unexpected error: %s", err)
	}

	if got, want := *prr.State, ReviewApprove; got != want {
		t.Fatalf("GitHubClient.CreateReview returns unexpected status: want: %s, got %s", want, got)
	}
}
