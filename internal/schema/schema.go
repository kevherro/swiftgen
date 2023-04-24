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
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type JSONSchemaProperty struct {
	Type    string `json:"type"`
	Format  string `json:"format,omitempty"`
	Default string `json:"default,omitempty"`
	Ref     string `json:"$ref,omitempty"`
}

type JSONSchema struct {
	Title      string                        `json:"title"`
	Type       string                        `json:"type"`
	Properties map[string]JSONSchemaProperty `json:"properties"`
	Required   []string                      `json:"required"`
}

type JSONSchemaToSwiftCodeGenerator struct {
	Schema JSONSchema
}

// Generate recursively processes referenced schemas.
func (g *JSONSchemaToSwiftCodeGenerator) Generate() (string, error) {
	var b strings.Builder
	err := g.generate(g.Schema, &b)
	if err != nil {
		return "", err
	}
	return b.String(), nil
}

func (g *JSONSchemaToSwiftCodeGenerator) generate(s JSONSchema, b *strings.Builder) error {
	// Open the struct.
	b.WriteString(fmt.Sprintf("struct %s: Codable {\n", s.Title))

	// Generate properties without $ref.
	for n, p := range s.Properties {
		if p.Ref == "" {
			name := convertToPascalCase(n)
			t := swiftType(p.Type)

			req := false
			for _, rp := range s.Required {
				if rp == name {
					req = true
					break
				}
			}

			if req {
				b.WriteString(fmt.Sprintf("\tlet %s: %s\n", name, t))
			} else {
				b.WriteString(fmt.Sprintf("\tlet %s: %s?\n", name, t))
			}
		}
	}

	// Generate properties with $ref and their respective structs.
	for n, p := range s.Properties {
		if p.Ref != "" {
			// The schema for p is located in file p.Ref.
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
	}

	// Close the struct.
	b.WriteString("}\n")

	return nil
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
func swiftType(jsonType string) string {
	switch jsonType {
	case "string":
		return "String"
	case "integer":
		return "Int"
	case "number":
		return "Double"
	case "boolean":
		return "Bool"
	case "array":
		// Assuming the array element type is known, e.g., String.
		return "[String]" // TODO: don't assume this!
	default:
		return "Any"
	}
}

// loadJSONSchemaFromFile loads JSON schema data from name. In case of an error,
// an empty JSONSchema is returned.
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

	return s, err
}
