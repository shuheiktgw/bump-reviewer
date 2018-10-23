// +build integration

package main

import (
	"testing"
)

func TestReviewer_Integration_Review_Success(t *testing.T) {
	r := Reviewer{integrationGitHubClient}
	err := r.Review(3)
	if err != nil {
		t.Fatalf("Unexpected error has returned reviewer.Review: %s", err)
	}
}

func TestReviewer_Integration_Review_Fail(t *testing.T) {
	cases := []struct {
		prNum int
	}{
		{prNum: 4},
		{prNum: 5},
	}

	for i, tc := range cases {
		r := Reviewer{integrationGitHubClient}
		err := r.Review(tc.prNum)
		if _, ok := err.(review); !ok {
			t.Fatalf("#%d Unexpected error has returned from reviewer.Review: %s", i, err)
		}
	}
}
