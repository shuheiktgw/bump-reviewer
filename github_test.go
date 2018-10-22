package main

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/google/go-github/github"
)

func TestGitHubClient_CompareCommits(t *testing.T) {
	client, mux, _, tearDown := setup()
	defer tearDown()

	base := "master"
	head := "develop"

	mux.HandleFunc(fmt.Sprintf("/repos/%v/%v/compare/%v...%v", testGitHubOwner, testGitHubRepo, base, head), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `{"status":"ahead"}`)
	})

	cc, err := client.CompareCommits(base, head)
	if err != nil {
		t.Fatalf("GitHubClient.CompareCommits returned unexpected error: %v", err)
	}

	want := &github.CommitsComparison{Status: github.String("ahead")}
	if !reflect.DeepEqual(cc, want) {
		t.Errorf("GitHubClient.CompareCommits returned %+v, want %+v", cc, want)
	}
}

func TestGitHubClient_ListPullRequestsFiles(t *testing.T) {
	client, mux, _, tearDown := setup()
	defer tearDown()

	number := 1

	mux.HandleFunc(fmt.Sprintf("/repos/%v/%v/pulls/%d/files", testGitHubOwner, testGitHubRepo, number), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `[{"filename":"version.rb"}]`)
	})

	cc, err := client.ListPullRequestsFiles(number, nil)
	if err != nil {
		t.Fatalf("GitHubClient.ListPullRequestsFiles returned unexpected error: %v", err)
	}

	want := []*github.CommitFile{{Filename: github.String("version.rb")}}
	if !reflect.DeepEqual(cc, want) {
		t.Errorf("GitHubClient.ListPullRequestsFiles returned %+v, want %+v", cc, want)
	}
}

func TestGitHubClient_GetLatestRelease(t *testing.T) {
	client, mux, _, tearDown := setup()
	defer tearDown()

	mux.HandleFunc(fmt.Sprintf("/repos/%s/%s/releases/latest", testGitHubOwner, testGitHubRepo), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `{"tag_name":"v0.0.1","name":"Release v0.0.1"}`)
	})

	rr, err := client.GetLatestRelease()
	if err != nil {
		t.Fatalf("GitHubClient.GetLatestRelease returned unexpected error: %v", err)
	}

	want := &github.RepositoryRelease{TagName: github.String("v0.0.1"), Name: github.String("Release v0.0.1")}
	if !reflect.DeepEqual(rr, want) {
		t.Errorf("GitHubClient.GetLatestRelease returned %+v, want %+v", rr, want)
	}
}

func TestGitHubClient_GetContent(t *testing.T) {
	client, mux, _, tearDown := setup()
	defer tearDown()

	path := "version.rb"
	u := fmt.Sprintf("/repos/%s/%s/contents/%s", testGitHubOwner, testGitHubRepo, path)

	mux.HandleFunc(u, func(w http.ResponseWriter, r *http.Request) {
		testFormValues(t, r, values{"ref": "develop"})
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `{"name":"version.rb"}`)
	})

	opt := github.RepositoryContentGetOptions{Ref: "develop"}
	fc, _, err := client.GetContent(path, &opt)
	if err != nil {
		t.Fatalf("GitHubClient.GetContent returned unexpected error: %v", err)
	}

	want := &github.RepositoryContent{Name: github.String("version.rb")}
	if !reflect.DeepEqual(fc, want) {
		t.Errorf("GitHubClient.GetContent returned %+v, want %+v", fc, want)
	}
}
