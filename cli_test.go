package main

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

func TestCLI_Run(t *testing.T) {
	cases := []struct {
		command           string
		expectedOutStream string
		expectedErrStream string
		expectedExitCode  int
	}{
		{
			command:           "bump-reviewer",
			expectedOutStream: "",
			expectedErrStream: "Failed to set up bump-reviewer: GitHub owner is missing\nPlease set it via `-o` option\n\n",
			expectedExitCode:  ExitCodeInvalidFlagError,
		},
		{
			command:           "bump-reviewer -o shuheiktgw",
			expectedOutStream: "",
			expectedErrStream: "Failed to set up bump-reviewer: GitHub repository is missing\nPlease set it via `-r` option\n\n",
			expectedExitCode:  ExitCodeInvalidFlagError,
		},
		{
			command:           "bump-reviewer -o shuheiktgw -r bump-reviewer",
			expectedOutStream: "",
			expectedErrStream: "Failed to set up bump-reviewer: GitHub Personal Access Token is missing\nPlease set it via `-t` option\n\n",
			expectedExitCode:  ExitCodeInvalidFlagError,
		},
		{
			command:           "bump-reviewer -o shuheiktgw -r bump-reviewer -t 1234abcd",
			expectedOutStream: "",
			expectedErrStream: "Failed to set up bump-reviewer: Pull Request number is missing\nPlease set it via `-n` option\n\n",
			expectedExitCode:  ExitCodeInvalidFlagError,
		},
		{
			command:           "bump-reviewer -o shuheiktgw -r bump-reviewer -t 1234abcd -n 1",
			expectedOutStream: "",
			expectedErrStream: "bump-reviewer failed to review because of the following error.\n\n" +
				"GET https://api.github.com/repos/shuheiktgw/bump-reviewer/pulls/1/files: 401 Bad credentials []\n\n" +
				"You might encounter a bug with bump-reviewer, so please report it to https://github.com/shuheiktgw/bump-reviewer/issues\n\n",
			expectedExitCode: ExitCodeError,
		},
		{
			command:           "bump-reviewer -v",
			expectedOutStream: fmt.Sprintf("bump-reviewer current version v%s\n", Version),
			expectedErrStream: "",
			expectedExitCode:  ExitCodeOK,
		},
	}

	for i, tc := range cases {
		outStream := new(bytes.Buffer)
		errStream := new(bytes.Buffer)

		cli := CLI{outStream: outStream, errStream: errStream}
		args := strings.Split(tc.command, " ")

		if got := cli.Run(args); got != tc.expectedExitCode {
			t.Fatalf("#%d %q exits with %d, want %d", i, tc.command, got, tc.expectedExitCode)
		}

		if got := outStream.String(); got != tc.expectedOutStream {
			t.Fatalf("#%d Unexpected outStream has returned: want: %s, got: %s", i, tc.expectedOutStream, got)
		}

		if got := errStream.String(); got != tc.expectedErrStream {
			t.Fatalf("#%d Unexpected errStream has returned: want: %s, got: %s", i, tc.expectedErrStream, got)
		}
	}
}
