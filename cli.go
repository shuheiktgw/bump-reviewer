package main

import (
	"flag"
	"fmt"
	"io"
)

const (
	ExitCodeOK = iota
	ExitCodeError
	ExitCodeReviewFailed
	ExitCodeParseFlagsError
	ExitCodeInvalidFlagError
)

type CLI struct {
	outStream, errStream io.Writer
}

func (cli *CLI) Run(args []string) int {
	var (
		owner   string
		repo    string
		token   string
		number  int
		version bool
	)

	flags := flag.NewFlagSet(Name, flag.ContinueOnError)
	flags.Usage = func() {
		fmt.Fprint(cli.outStream, usage)
	}

	flags.StringVar(&owner, "owner", "", "")
	flags.StringVar(&owner, "o", "", "")

	flags.StringVar(&repo, "repo", "", "")
	flags.StringVar(&repo, "r", "", "")

	flags.StringVar(&token, "token", "", "")
	flags.StringVar(&token, "t", "", "")

	flags.IntVar(&number, "number", 0, "")
	flags.IntVar(&number, "n", 0, "")

	flags.BoolVar(&version, "version", false, "")
	flags.BoolVar(&version, "v", false, "")

	if err := flags.Parse(args[1:]); err != nil {
		return ExitCodeParseFlagsError
	}

	if version {
		fmt.Fprint(cli.outStream, OutputVersion())
		return ExitCodeOK
	}

	if len(owner) == 0 {
		fmt.Fprintf(cli.errStream, "Failed to set up bump-reviewer: GitHub owner is missing\n"+
			"Please set it via `-o` option\n\n")
		return ExitCodeInvalidFlagError
	}

	if len(repo) == 0 {
		fmt.Fprintf(cli.errStream, "Failed to set up bump-reviewer: GitHub repository is missing\n"+
			"Please set it via `-r` option\n\n")
		return ExitCodeInvalidFlagError
	}

	if len(token) == 0 {
		fmt.Fprintf(cli.errStream, "Failed to set up bump-reviewer: GitHub Personal Access Token is missing\n"+
			"Please set it via `-t` option\n\n")
		return ExitCodeInvalidFlagError
	}

	if number == 0 {
		fmt.Fprintf(cli.errStream, "Failed to set up bump-reviewer: Pull Request number is missing\n"+
			"Please set it via `-n` option\n\n")
		return ExitCodeInvalidFlagError
	}

	client := NewGitHubClient(owner, repo, token)
	reviewer := Reviewer{client}

	if err := reviewer.Review(number); err != nil {
		if r, ok := err.(review); ok {
			fmt.Fprintf(cli.errStream, "Pull Request #%d did not pass the review because of the following reason\n\n%s", number, r.review())
			return ExitCodeReviewFailed
		}
		fmt.Fprintf(cli.errStream, `bump-reviewer failed to review because of the following error.

%s

You might encounter a bug with bump-reviewer, so please report it to https://github.com/shuheiktgw/bump-reviewer/issues

`, err)
		return ExitCodeError
	}

	fmt.Fprintf(cli.outStream, "bump-reviewer successfully approved your Pull Request.\n\n")
	return ExitCodeOK
}

var usage = `Usage: bump-reviewer [options...]

bump-reviewer is a command to review and approve bump up Pull Requests

OPTIONS:
  --number value, -n value  specifies the number of last lines of file (default 10)
  --quiet, -q               suppresses printing of headers when multiple files are being examined
  --version, -v             prints the current version
  --help, -h                prints help

`
