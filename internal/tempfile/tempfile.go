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

// Package tempfile...
package tempfile

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

// NewTempFile returns a new output file in dir with the provided prefix and suffix.
func NewTempFile(dir, prefix string) (*os.File, error) {
	suffix := ".swift"
	switch f, err := os.OpenFile(filepath.Join(dir, fmt.Sprintf("%s%s", prefix, suffix)), os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666); {
	case err == nil:
		return f, nil
	case !os.IsExist(err):
		return nil, err
	}
	// "Just give up!" - Mr. S
	return nil, fmt.Errorf("could not create file of the form %s%s", prefix, suffix)
}

var tempFiles []string
var tempFilesMu = sync.Mutex{}

// deferDeleteTempFile marks a file to be deleted by the next call to Cleanup().
func deferDeleteTempFile(path string) {
	tempFilesMu.Lock()
	tempFiles = append(tempFiles, path)
	tempFilesMu.Unlock()
}

// cleanupTempFiles removes any temporary files selected for deferred cleanup.
func cleanupTempFiles() error {
	tempFilesMu.Lock()
	defer tempFilesMu.Unlock()
	var lastErr error
	for _, f := range tempFiles {
		if err := os.Remove(f); err != nil {
			lastErr = err
		}
	}
	tempFiles = nil
	return lastErr
}
