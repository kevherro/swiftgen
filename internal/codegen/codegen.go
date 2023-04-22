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

// Package codegen implements code generation.
package codegen

import (
	"encoding/json"
	"errors"
	"io"
	"os"

	"github.com/kevherro/swiftgen/internal/schema"
	"github.com/kevherro/swiftgen/internal/tempfile"
)

func Generate() error {
	args := os.Args
	if len(args) != 2 { // Executable + input file
		return errors.New("usage: ./swiftgen <source>")
	}

	source, err := os.Open(args[1])
	if err != nil {
		return err
	}
	defer source.Close()

	byteValue, _ := io.ReadAll(source)

	var jsonSchema schema.JSONSchema
	if err = json.Unmarshal(byteValue, &jsonSchema); err != nil {
		return err
	}

	generator := schema.JSONSchemaToSwiftCodeGenerator{Schema: jsonSchema}

	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	out, err := tempfile.NewTempFile(wd, generator.Schema.Title)
	if err != nil {
		return err
	}

	code := generator.Generate()
	if err = os.WriteFile(out.Name(), []byte(code), 0644); err != nil {
		return err
	}

	// Success!
	return nil
}
