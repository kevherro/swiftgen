// MIT License
//
// Copyright (c) 2023 Kevin Herro
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

package util

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const (
	prefix = "testfile"
	suffix = ".swift"
)

func TestNewTempFile(t *testing.T) {
	t.Cleanup(func() {
		os.Remove(os.TempDir() + prefix + suffix)
	})

	f, err := NewTempFile(os.TempDir(), prefix, suffix)
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer f.Close()

	// Verify that the file was created with the expected name.
	expectedName := filepath.Join(os.TempDir(), prefix+suffix)
	if got := f.Name(); got != expectedName {
		t.Errorf("unexpected file name: got %q, expected %q", got, expectedName)
	}

	// Verify that the file contains the expected content.
	expectedContent := "hello world"
	if _, err := f.WriteString(expectedContent); err != nil {
		t.Fatalf("failed to write to file: %v", err)
	}
	if err := f.Sync(); err != nil {
		t.Fatalf("failed to sync file: %v", err)
	}
	content, err := os.ReadFile(expectedName)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}
	if got := strings.TrimSpace(string(content)); got != expectedContent {
		t.Errorf("unexpected file content: got %q, expected %q", got, expectedContent)
	}
}
