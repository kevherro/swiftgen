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

package schema

import (
	"strings"
	"testing"
)

func TestSwiftType(t *testing.T) {
	tests := []struct {
		jsonProperty JSONSchemaProperty
		expected     string
	}{
		{JSONSchemaProperty{Type: "string"}, "String"},
		{JSONSchemaProperty{Type: "integer"}, "Int"},
		{JSONSchemaProperty{Type: "number"}, "Double"},
		{JSONSchemaProperty{Type: "boolean"}, "Bool"},
		{JSONSchemaProperty{Type: "array", Items: &JSONSchemaProperty{Type: "string"}}, "[String]"},
		{JSONSchemaProperty{Type: "array", Items: &JSONSchemaProperty{Type: "array", Items: &JSONSchemaProperty{Type: "integer"}}}, "[[Int]]"},
		{JSONSchemaProperty{Type: "unknown"}, "Any"},
	}

	for _, test := range tests {
		result := swiftType(test.jsonProperty)
		if result != test.expected {
			t.Errorf("swiftType(%v) = %s; want %s", test.jsonProperty, result, test.expected)
		}
	}
}

func TestGenerate(t *testing.T) {
	schema := JSONSchema{
		Title: "User",
		Type:  "object",
		Properties: map[string]JSONSchemaProperty{
			"id":        {Type: "integer"},
			"name":      {Type: "string"},
			"age":       {Type: "integer"},
			"email":     {Type: "string"},
			"is_active": {Type: "boolean"},
		},
		Required: []string{"id", "name"},
	}

	g := JSONSchemaToSwiftCodeGenerator{Schema: schema}
	c, err := g.Generate()
	if err != nil {
		t.Errorf("Generate() error = %v", err)
	}

	expectedProperties := []string{
		"let Id: Int",
		"let Name: String",
		"let Age: Int?",
		"let Email: String?",
		"let IsActive: Bool?",
	}

	// Check if all expected properties are present in the generated code.
	for _, p := range expectedProperties {
		if !strings.Contains(c, p) {
			t.Errorf("Generated code is missing property: %s\n\nGenerated code:\n\n%s", p, c)
		}
	}
}

func TestLoadJSONSchemaFromFile(t *testing.T) {
	tests := []struct {
		name    string
		file    string
		wantErr bool
	}{
		{
			name:    "valid JSON schema",
			file:    "testdata/valid.schema.json",
			wantErr: false,
		},
		{
			name:    "invalid JSON schema",
			file:    "testdata/invalid.schema.json",
			wantErr: true,
		},
		{
			name:    "non-existent JSON schema",
			file:    "testdata/non_existent.schema.json",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := loadJSONSchemaFromFile(tt.file)
			if (err != nil) != tt.wantErr {
				t.Errorf("loadJSONSchemaFromFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
