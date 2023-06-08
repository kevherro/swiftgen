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
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/kevherro/swiftgen/internal/flags"
	"github.com/kevherro/swiftgen/internal/testutil"
)

func TestGenerateFromJSONSchema(t *testing.T) {
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
	cmdFlags := &flags.CmdFlags{
		Src:  f.Name(),
		Dest: dest,
	}
	if err := SwiftGen(cmdFlags); err != nil {
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

// TestGenerateFromJSONSchemaWithRefs tests schemas that are referenced in one
// file and exists in another file.
func TestGenerateFromJSONSchemaWithRefs(t *testing.T) {
	f2, cleanup2 := testutil.CreateTempTestFile(t, "bar_", ".json")
	blob := `{
		"title": "Address",
		"type": "object",
		"properties": {
			"street": { "type": "string" },
			"city": { "type": "number" },
			"country": { "type": "string" }
		},
		"required": ["street", "city", "country"]
	}`
	if _, err := f2.Write([]byte(blob)); err != nil {
		t.Cleanup(cleanup2)
		t.Errorf("f.Write() error = %v", err)
	}
	defer t.Cleanup(cleanup2)

	f1, cleanup1 := testutil.CreateTempTestFile(t, "foo_", ".json")
	blob = fmt.Sprintf(`{
		"title": "Foo",
		"type": "object",
		"properties": {
			"name": { "type": "string" },
			"age": { "type": "number" },
			"address": { "$ref": "%s" }
		},
		"required": ["name", "age", "address"]
	}`, f2.Name())
	if _, err := f1.Write([]byte(blob)); err != nil {
		t.Cleanup(cleanup1)
		t.Errorf("f.Write() error = %v", err)
	}
	defer t.Cleanup(cleanup1)

	// Call Generate with the temporary file paths.
	dest := filepath.Dir(f1.Name()) + "/" + "Schema.swift"
	cmdFlags := &flags.CmdFlags{
		Src:  f1.Name(),
		Dest: dest,
	}
	if err := SwiftGen(cmdFlags); err != nil {
		t.Errorf("Generate() error = %v", err)
	}
	defer os.Remove(dest) // clean up

	// Read the output file and check if it contains the generated code.
	generatedCode, err := os.ReadFile(dest)
	if err != nil {
		t.Errorf("ReadFile(%v) error = %v", dest, err)
	}

	expectedProperties := []string{
		"let Name: String",
		"let Age: Double",
		"let Address: Address",
	}

	// Check if all expected properties are present in the generated code.
	for _, property := range expectedProperties {
		if !strings.Contains(string(generatedCode), property) {
			t.Errorf("Generated code is missing property: %s\n\nGenerated code:\n\n%s", property, generatedCode)
		}
	}
}

func TestSwiftGen(t *testing.T) {
	// Write a test JSON schema to a temp file.
	f, err := os.CreateTemp(os.TempDir(), "schema.json")
	if err != nil {
		t.Errorf("CreateTemp(%v, %v) error = %v", os.TempDir(), "schema.json", err)
	}
	defer os.Remove(f.Name())

	blob := `{
		"title": "Person",
		"type": "object",
		"properties": {
			"name": { "type": "string" },
			"age": { "type": "number" },
			"address": { "$ref": "#/definitions/Address" },
			"hobbies": {
				"type": "array",
				"items": { "type": "string" }
			}
		},
		"required": ["name", "age", "address", "hobbies"],
		"definitions": {
			"Address": {
				"title": "Address",
				"type": "object",
				"properties": {
					"street": { "type": "string" },
					"city": { "type": "string" }
				},
				"required": ["street", "city"]
			}
		}
	}`
	if _, err := f.Write([]byte(blob)); err != nil {
		f.Close()
		t.Errorf("f.Write() error = %v", err)
	}
	defer f.Close()

	// Call SwiftGen with the temp file's directory
	// as the root directory.
	cmdFlags := &flags.CmdFlags{
		Src:  f.Name(),
		Dest: filepath.Join(os.TempDir(), "Person.swift"),
	}
	err = SwiftGen(cmdFlags)
	if err != nil {
		t.Errorf("SwiftGen() error = %v", err)
	}

	// Read the generated Swift file and verify that
	// it contains the expected code.
	generatedCode, err := os.ReadFile(cmdFlags.Dest)
	if err != nil {
		t.Errorf("ReadFile(%v) error = %v", cmdFlags.Dest, err)
	}

	// Delete the generated swift file after the test.
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
