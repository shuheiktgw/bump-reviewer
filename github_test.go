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
