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
		jsonType string
		expected string
	}{
		{"string", "String"},
		{"integer", "Int"},
		{"number", "Double"},
		{"boolean", "Bool"},
		{"array", "[String]"},
		{"unknown", "Any"},
	}

	for _, test := range tests {
		result := swiftType(test.jsonType)
		if result != test.expected {
			t.Errorf("swiftType(%s) = %s; want %s", test.jsonType, result, test.expected)
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
