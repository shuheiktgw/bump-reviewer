package main

import (
	"encoding/base64"
	"fmt"
	"regexp"
	"strings"

	"github.com/blang/semver"
	"github.com/google/go-github/github"
	"github.com/iancoleman/strcase"
)

type review interface {
	review() string
}

type reviewError struct {
	Message string
}

func (r *reviewError) Error() string {
	return r.Message
}

func (r *reviewError) review() string {
	return r.Message
}

// Reviewer reviews bump up PRs
type Reviewer struct {
	*GitHubClient
}

// Review reviews a bump up PR
func (r *Reviewer) Review(number int) error {
	// Check if the PR changes only the version.rb file
	if err := r.reviewFile(number); err != nil {
		return r.handleReviewError(number, err)
	}

	// Check if the PR's version.rb follows the expected pattern
	if err := r.reviewVersion(number); err != nil {
		return r.handleReviewError(number, err)
	}

	// Approve the PR
	if err := r.approvePullRequest(number); err != nil {
		return err
	}

	return nil
}

func (r *Reviewer) reviewFile(number int) error {
	files, err := r.ListPullRequestsFiles(number, nil)
	if err != nil {
		return err
	}

	if len(files) != 1 {
		return &reviewError{Message: fmt.Sprintf("Pull Request #%d edited more than one file. bump-reviewer only allows to edit one file, which is `version.rb`.", number)}
	}

	filename := fmt.Sprintf("lib/%s/version.rb", r.Repo)
	if *files[0].Filename != filename {
		return &reviewError{Message: fmt.Sprintf("Pull Request #%d edited edited unexpected file, bump-reviewer only allows to edit %s.", number, filename)}
	}

	return nil
}

func (r *Reviewer) reviewVersion(number int) error {
	release, err := r.GetLatestRelease()
	if err != nil {
		return err
	}
	tag := *release.TagName

	opt := github.RepositoryContentGetOptions{Ref: fmt.Sprintf("pull/%d/head", number)}
	fc, _, err := r.GetContent(fmt.Sprintf("lib/%s/version.rb", r.Repo), &opt)
	if err != nil {
		return err
	}

	content, err := decodeContent(fc)
	if err != nil {
		return err
	}

	// Trim the prefix "v" or "V"
	trimmedTag := strings.TrimPrefix(tag, "v")
	trimmedTag = strings.TrimPrefix(trimmedTag, "V")

	if err := r.checkVersionRegex(trimmedTag, content); err != nil {
		return err
	}

	return nil
}

func (r *Reviewer) handleReviewError(number int, err error) error {
	if review, ok := err.(review); ok {
		if err := r.postComment(number, review.review()); err != nil {
			return err
		}
	}

	return err
}

func (r *Reviewer) postComment(number int, comment string) error {
	review := github.PullRequestReviewRequest{Event: github.String(ReviewComment), Body: github.String(comment)}
	_, err := r.CreateReview(number, &review)
	if err != nil {
		return err
	}

	return nil
}

func (r *Reviewer) approvePullRequest(number int) error {
	body := `LGTM

bump-reviewer checks the following two points.

- PR changes only version.rb
- PR increments patch version by one 
`
	approve := github.PullRequestReviewRequest{Event: github.String(ReviewApprove), Body: github.String(body)}
	_, err := r.CreateReview(number, &approve)
	if err != nil {
		return err
	}

	return nil
}

func decodeContent(rc *github.RepositoryContent) (string, error) {
	if *rc.Encoding != "base64" {
		return "", fmt.Errorf("unexpected encoding: %s", *rc.Encoding)
	}

	decoded, err := base64.StdEncoding.DecodeString(*rc.Content)

	if err != nil {
		return "", err
	}

	return string(decoded), nil
}

func (r *Reviewer) checkVersionRegex(tag, content string) error {
	appName := strcase.ToCamel(r.Repo)

	v, err := semver.New(tag)
	if err != nil {
		return err
	}

	v.Patch = v.Patch + 1
	newTag := v.String()

	regStr := fmt.Sprintf(`\s*module\s+%s\s+VERSION\s*=\s*['"]%s['"](\.freeze)?\s+end\s*`, appName, newTag)
	reg := regexp.MustCompile(regStr)
	if !reg.Match([]byte(content)) {
		return &reviewError{Message: fmt.Sprintf("version.rb does not match with the following regex: `%s`. bump-reviewer expects you increment patch version by one.", regStr)}
	}

	return nil
}
