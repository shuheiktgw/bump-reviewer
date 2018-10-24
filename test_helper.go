package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
)

const (
	testGitHubOwner = "shuheiktgw"
	testGitHubRepo  = "bump-reviewer"
	testGitHubToken = "abcdefg12345"
)

// setup sets up a test HTTP server along with a GitHubClient that is
// configured to talk to that test server. Tests should register handlers on
// mux which provide mock responses for the API method being tested.
func setup() (client *GitHubClient, mux *http.ServeMux, serverURL string, tearDown func()) {
	// mux is the HTTP request multiplexer used with the test server.
	mux = http.NewServeMux()

	apiHandler := http.NewServeMux()
	apiHandler.Handle("/", mux)

	// server is a test HTTP server used to provide mock API responses.
	server := httptest.NewServer(apiHandler)

	// client is the GitHub client being tested and is
	// configured to use test server.
	client = NewGitHubClient(testGitHubOwner, testGitHubRepo, testGitHubToken)
	u, _ := url.Parse(server.URL + "/")
	client.Client.BaseURL = u

	return client, mux, server.URL, server.Close
}

func testMethod(t *testing.T, r *http.Request, want string) {
	if got := r.Method; got != want {
		t.Errorf("Request method: %v, want %v", got, want)
	}
}

type values map[string]string

func testFormValues(t *testing.T, r *http.Request, values values) {
	want := url.Values{}
	for k, v := range values {
		want.Set(k, v)
	}

	r.ParseForm()
	if got := r.Form; !reflect.DeepEqual(got, want) {
		t.Errorf("Request parameters: %v, want %v", got, want)
	}
}

func testBody(t *testing.T, r *http.Request, want string) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		t.Errorf("Error reading request body: %v", err)
	}
	if got := string(b); got != want {
		t.Errorf("request Body is %s, want %s", got, want)
	}
}

func setupReviewer() (reviewer *Reviewer, mux *http.ServeMux, url string, tearDown func()) {
	client, mux, url, tearDown := setup()
	return &Reviewer{client}, mux, url, tearDown
}

func setPullRequestFilesHandler(mux *http.ServeMux, number int, files string) {
	mux.HandleFunc(fmt.Sprintf("/repos/%v/%v/pulls/%d/files", testGitHubOwner, testGitHubRepo, number), func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, files)
	})
}

func setCreateReviewHandler(mux *http.ServeMux, number int, state string) {
	mux.HandleFunc(fmt.Sprintf("/repos/%v/%v/pulls/%d/reviews", testGitHubOwner, testGitHubRepo, number), func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, fmt.Sprintf(`{"state":"%s"}`, state))
	})
}

func setReleaseHandler(mux *http.ServeMux, tag string) {
	mux.HandleFunc(fmt.Sprintf("/repos/%s/%s/releases/latest", testGitHubOwner, testGitHubRepo), func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{"tag_name":"%s"}`, tag)
	})
}

func setGetContentHandler(mux *http.ServeMux, version string) {
	path := fmt.Sprintf("lib/%s/version.rb", testGitHubRepo)
	content := fmt.Sprintf(`
module BumpReviewer
  VERSION="%s"
end
`, version)

	mux.HandleFunc(fmt.Sprintf("/repos/%s/%s/contents/%s", testGitHubOwner, testGitHubRepo, path), func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{"content":"%s","encoding":"base64"}`, base64.StdEncoding.EncodeToString([]byte(content)))
	})
}
