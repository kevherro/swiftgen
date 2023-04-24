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

// Package testutil implements some useful helper functions for testing.
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
