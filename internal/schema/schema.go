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

// Package schema implements JSON schema serialization and code generation.
package schema

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type JSONSchemaProperty struct {
	Type    string              `json:"type"`
	Format  string              `json:"format,omitempty"`
	Default string              `json:"default,omitempty"`
	Ref     string              `json:"$ref,omitempty"`
	Items   *JSONSchemaProperty `json:"items"`
}

type JSONSchema struct {
	Title       string                        `json:"title"`
	Type        string                        `json:"type"`
	Properties  map[string]JSONSchemaProperty `json:"properties"`
	Required    []string                      `json:"required"`
	Definitions map[string]JSONSchema         `json:"definitions"`
}

type JSONSchemaToSwiftCodeGenerator struct {
	Schema JSONSchema
}

// Generate recursively processes referenced schemas.
func (g *JSONSchemaToSwiftCodeGenerator) Generate() (string, error) {
	var b strings.Builder
	err := g.generate(g.Schema, &b)
	return b.String(), err
}

// generate does not support anonymous structs and complex $ref URIs,
// including network URIs or URIs pointing to a location in an external file.
func (g *JSONSchemaToSwiftCodeGenerator) generate(s JSONSchema, b *strings.Builder) error {
	// Open the struct.
	b.WriteString(fmt.Sprintf("struct %s: Codable {\n", s.Title))

	// Generate properties.
	for n, p := range s.Properties {
		if p.Ref != "" {
			if strings.HasPrefix(p.Ref, "#") {
				// The schema for p is referenced in the
				// definitions of the same file.
				definitionKey := strings.TrimPrefix(p.Ref, "#/definitions/")
				definedSchema, exists := s.Definitions[definitionKey]
				if !exists {
					return fmt.Errorf("unable to find definition for $ref: %s", p.Ref)
				}

				// Recursively generate the defined schema.
				err := g.generate(definedSchema, b)
				if err != nil {
					return err
				}

				name := convertToPascalCase(n)
				b.WriteString(fmt.Sprintf("\tlet %s: %s\n", name, definedSchema.Title))
			} else {
				// The schema for p is located in an
				// external file referenced by p.Ref.
				s, err := loadJSONSchemaFromFile(p.Ref)
				if err != nil {
					return err
				}

				err = g.generate(s, b)
				if err != nil {
					return err
				}

				name := convertToPascalCase(n)
				b.WriteString(fmt.Sprintf("\tlet %s: %s\n", name, s.Title))
			}
		} else {
			name := convertToPascalCase(n)
			t := swiftType(p)

			req := contains(s.Required, name)

			if req {
				b.WriteString(fmt.Sprintf("\tlet %s: %s\n", name, t))
			} else {
				b.WriteString(fmt.Sprintf("\tlet %s: %s?\n", name, t))
			}
		}
	}

	// Close the struct.
	b.WriteString("}\n")

	return nil
}

// Contains checks if a string is present in a slice of strings.
func contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}

// convertToPascalCase converts an input string in snake_case or kebab-case
// to PascalCase. It handles dashes and underscores as word separators.
func convertToPascalCase(in string) string {
	in = strings.Replace(in, "-", "_", -1)
	words := strings.Split(in, "_")

	// Convert each word to TitleCase using the title caser.
	title := cases.Title(language.English)
	for i, w := range words {
		words[i] = title.String(w)
	}

	return strings.Join(words, "")
}

// swiftType returns the Swift type for the jsonType.
func swiftType(p JSONSchemaProperty) string {
	switch p.Type {
	case "string":
		return "String"
	case "integer":
		return "Int"
	case "number":
		return "Double"
	case "boolean":
		return "Bool"
	case "array":
		if p.Items != nil {
			// Recursively handle the type of array items.
			return "[" + swiftType(*p.Items) + "]"
		} else {
			// If the type of the array items is not specified,
			// fall back to Any.
			return "[Any]"
		}
	default:
		return "Any"
	}
}

// loadJSONSchemaFromFile loads JSON schema data from name.
// In case of an error, an empty JSONSchema is returned.
func loadJSONSchemaFromFile(name string) (JSONSchema, error) {
	var s JSONSchema

	f, err := os.Open(name)
	if err != nil {
		return s, err
	}
	defer f.Close()

	b, err := io.ReadAll(f)
	if err != nil {
		return s, err
	}
	err = json.Unmarshal(b, &s)

	if err != nil {
		return s, err
	}

	// Validate the schema.
	for _, p := range s.Properties {
		if p.Type == "" {
			return s, errors.New("invalid schema: missing type in properties")
		}
	}

	return s, nil
}
