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

func swiftType(jsonType, format string) string {
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

func (g *JSONSchemaToSwiftCodeGenerator) Generate() string {
	var code strings.Builder

	// Generate the struct definition.
	structName := g.Schema.Title
	code.WriteString(fmt.Sprintf("struct %s: Codable {\n", structName))

	// Generate properties.
	for propertyName, property := range g.Schema.Properties {
		swiftPropertyName := cases.Title(language.English).String(propertyName)
		swiftPropertyType := swiftType(property.Type, property.Format)

		isRequired := false
		for _, requiredProperty := range g.Schema.Required {
			if requiredProperty == propertyName {
				isRequired = true
				break
			}
		}

		if isRequired {
			code.WriteString(fmt.Sprintf("\tlet %s: %s\n", swiftPropertyName, swiftPropertyType))
		} else {
			code.WriteString(fmt.Sprintf("\tlet %s: %s?\n", swiftPropertyName, swiftPropertyType))
		}
	}

	// Close the struct.
	code.WriteString("}\n")

	return code.String()
}
