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

	"github.com/kevherro/swiftgen/internal/flags"
)

func TestGenerateFromJSONSchemaFile(t *testing.T) {
	schemaFile := "testdata/person.schema.json"

	// Check that the file exists.
	if _, err := os.Stat(schemaFile); os.IsNotExist(err) {
		t.Fatalf("File does not exist: %s", schemaFile)
	}

	// Call SwiftGen with the schema file.
	destFile := filepath.Join(os.TempDir(), "Person.swift")
	cmdFlags := &flags.CmdFlags{
		Src:  schemaFile,
		Dest: destFile,
	}

	err := SwiftGen(cmdFlags)
	if err != nil {
		t.Errorf("SwiftGen() error = %v", err)
	}

	// Read the generated Swift file and verify that
	// it contains the expected code.
	generatedCode, err := os.ReadFile(cmdFlags.Dest)
	if err != nil {
		t.Errorf("ReadFile(%v) error = %v", cmdFlags.Dest, err)
	}

	// Delete the generated Swift file after the test.
	defer os.Remove(cmdFlags.Dest)

	expectedProperties := []string{
		"struct Person: Codable {",
		"let Name: String?",
		"let Age: Double?",
		"struct Address: Codable {",
		"let Street: String?",
		"let City: String?",
		"}",
		"let Address: Address",
		"let Hobbies: [String]?",
		"}",
	}

	// Check if all expected properties are present in the generated code.
	for _, property := range expectedProperties {
		if !strings.Contains(string(generatedCode), property) {
			t.Errorf("Generated code is missing property: %s\n\nGenerated code:\n\n%s", property, string(generatedCode))
		}
	}
}
