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

package codegen

import (
	"os"
	"strings"
	"testing"
)

func TestGenerate(t *testing.T) {
	// Create a temporary JSON schema file.
	tempDirName := "/tmp/"
	tempFile, err := os.CreateTemp(tempDirName, "schema.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDirName)

	// Write the JSON schema to the temporary file.
	jsonSchema := `{
		"title": "Root",
        "type": "object",
        "properties": {
            "name": { "type": "string" },
            "age": { "type": "number" }
        }
    }`
	if _, err := tempFile.Write([]byte(jsonSchema)); err != nil {
		t.Fatal(err)
	}

	// Close the temporary file.
	if err := tempFile.Close(); err != nil {
		t.Fatal(err)
	}

	dest := "Schema.swift"

	// Call Generate with the temporary file paths
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"cmd", "--src", tempFile.Name(), "--dest", dest}
	if err := Generate(); err != nil {
		t.Errorf("Generate() error = %v", err)
	}
	defer os.Remove(dest)

	// Read the output file and check if it contains the generated code
	generatedCode, err := os.ReadFile(dest)
	if err != nil {
		t.Fatal(err)
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
