package splitter_test

import (
	"os"
	"testing"

	"golangci-test/splitter"
)

func TestGetSplits(t *testing.T) {
	rawData, err := os.ReadFile("testdata/out.txt")
	if err != nil {
		t.Fatal(err)
	}

	// Call GetSplits with the test output
	splitter.GetSplits(string(rawData))
}