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

package driver

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGenerate(t *testing.T) {
	f, err := os.CreateTemp(os.TempDir(), "schema.json")
	if err != nil {
		t.Errorf("CreateTemp(%v, %v) error = %v", os.TempDir(), "schema.json", err)
	}
	defer os.Remove(f.Name())

	blob := `{
		"title": "Root",
        "type": "object",
        "properties": {
            "name": { "type": "string" },
            "age": { "type": "number" }
        }
    }`
	if _, err := f.Write([]byte(blob)); err != nil {
		f.Close()
		t.Errorf("f.Write() error = %v", err)
	}
	defer f.Close()

	// Call Generate with the temporary file paths.
	dest := filepath.Dir(f.Name()) + "/" + "Schema.swift"
	oldArgs := os.Args
	defer func() {
		os.Args = oldArgs
	}()
	os.Args = []string{"cmd", "--src", f.Name(), "--dest", dest}
	if err := SwiftGen(); err != nil {
		t.Errorf("Generate() error = %v", err)
	}
	defer os.Remove(dest) // clean up

	// Read the output file and check if it contains the generated code.
	generatedCode, err := os.ReadFile(dest)
	if err != nil {
		t.Errorf("ReadFile(%v) error = %v", dest, err)
	}

	expectedProperties := []string{
		"let Name: String?",
		"let Age: Double?",
	}

	// Check if all expected properties are present in the generated code.
	for _, property := range expectedProperties {
		if !strings.Contains(string(generatedCode), property) {
			t.Errorf("Generated code is missing property: %s\n\nGenerated code:\n\n%s", property, generatedCode)
		}
	}
}
