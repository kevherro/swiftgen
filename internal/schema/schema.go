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
	"fmt"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type JSONSchemaProperty struct {
	Type    string `json:"type"`
	Format  string `json:"format,omitempty"`
	Default string `json:"default,omitempty"`
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
		return "[String]"
	default:
		return "Any"
	}
}

// convertToPascalCase converts an input string in snake_case or kebab-case
// to PascalCase. It handles dashes and underscores as word separators.
func convertToPascalCase(in string) string {
	in = strings.Replace(in, "-", "_", -1)
	words := strings.Split(in, "_")

	// Convert each word to TitleCase using the title caser.
	title := cases.Title(language.English)
	for i, word := range words {
		words[i] = title.String(word)
	}

	return strings.Join(words, "")
}

func (g *JSONSchemaToSwiftCodeGenerator) Generate() string {
	var b strings.Builder

	// Generate the struct definition.
	sName := g.Schema.Title
	b.WriteString(fmt.Sprintf("struct %s: Codable {\n", sName))

	// Generate properties.
	for pName, p := range g.Schema.Properties {
		n := convertToPascalCase(pName)
		t := swiftType(p.Type)

		isRequired := false
		for _, requiredProperty := range g.Schema.Required {
			if requiredProperty == pName {
				isRequired = true
				break
			}
		}

		if isRequired {
			b.WriteString(fmt.Sprintf("\tlet %s: %s\n", n, t))
		} else {
			b.WriteString(fmt.Sprintf("\tlet %s: %s?\n", n, t))
		}
	}

	// Close the struct.
	b.WriteString("}\n")

	return b.String()
}
