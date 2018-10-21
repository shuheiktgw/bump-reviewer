// +build integration

package main

import (
	"testing"
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
