package main

import (
	"fmt"
	"testing"
)

func TestVersion_OutputVersion(t *testing.T) {
	if got, want := OutputVersion(), fmt.Sprintf("%s current version v%s\n", Name, Version); got != want {
		t.Fatalf("#OutputVersion returnes unexpected string, want: %s, got: %s", want, got)
	}
}
