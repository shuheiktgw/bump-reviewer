package main

import (
	"strings"
	"testing"
)

func TestReviewer_Review_FailWithTooManyFiles(t *testing.T) {
	reviewer, mux, _, tearDown := setupReviewer()
	defer tearDown()

	number := 1
	setPullRequestFilesHandler(mux, number, `[{"filename":"version.rb"}, {"filename":"version_spec.rb"}]`)
	setCreateReviewHandler(mux, number, "COMMENT")

	err := reviewer.Review(number)
	r, ok := err.(review)
	if !ok {
		t.Fatalf("Reviewer.Review returned unexpected error: %s", err)
	}

	if !strings.Contains(r.review(), "edited more than one file") {
		t.Fatalf("Reviewer.Review returned unexpected error: %s", err)
	}
}

func TestReviewer_Review_FailWithNonVersionFile(t *testing.T) {
	reviewer, mux, _, tearDown := setupReviewer()
	defer tearDown()

	number := 1
	setPullRequestFilesHandler(mux, number, `[{"filename":"test.rb"}]`)
	setCreateReviewHandler(mux, number, "COMMENT")

	err := reviewer.Review(number)
	r, ok := err.(review)
	if !ok {
		t.Fatalf("Reviewer.Review returned unexpected error: %s", err)
	}

	if !strings.Contains(r.review(), "edited unexpected file") {
		t.Fatalf("Reviewer.Review returned unexpected error: %s", err)
	}
}

func TestReviewer_Review_FailWithNonMatchVersion(t *testing.T) {
	reviewer, mux, _, tearDown := setupReviewer()
	defer tearDown()

	number := 1
	setPullRequestFilesHandler(mux, number, `[{"filename":"lib/bump-reviewer/version.rb"}]`)
	setCreateReviewHandler(mux, number, "COMMENT")
	setReleaseHandler(mux, "v1.0.1")
	setGetContentHandler(mux, "1.0.3")

	err := reviewer.Review(number)
	r, ok := err.(review)
	if !ok {
		t.Fatalf("Reviewer.Review returned unexpected error: %s", err)
	}

	if !strings.Contains(r.review(), "version.rb does not match with the following regex") {
		t.Fatalf("Reviewer.Review returned unexpected error: %s", err)
	}
}

func TestReviewer_Review_Success(t *testing.T) {
	reviewer, mux, _, tearDown := setupReviewer()
	defer tearDown()

	number := 1
	setPullRequestFilesHandler(mux, number, `[{"filename":"lib/bump-reviewer/version.rb"}]`)
	setCreateReviewHandler(mux, number, "COMMENT")
	setReleaseHandler(mux, "v1.0.1")
	setGetContentHandler(mux, "1.0.2")

	err := reviewer.Review(number)
	if err != nil {
		t.Fatalf("Reviewer.Review returned unexpected error: %s", err)
	}
}
