package testutil

import (
	"os"
	"testing"
)

func CreateTempTestFile(t *testing.T, prefix, suffix string) (*os.File, func()) {
	tempFile, err := os.CreateTemp("", prefix+"*"+suffix)
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	return tempFile, func() {
		tempFile.Close()
		os.Remove(tempFile.Name())
	}
}
