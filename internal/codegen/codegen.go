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
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/kevherro/swiftgen/internal/schema"
	"github.com/kevherro/swiftgen/internal/utils"
)

func Generate() error {
	srcFlag := flag.String("src", "", "Path to JSON schema file")
	destFlag := flag.String("dest", "", "Path to write location")

	flag.Parse()

	if *srcFlag == "" {
		flag.Usage()
		err := fmt.Errorf("codegen: --src flag value not provided")
		return err
	}

	src, err := os.Open(*srcFlag)
	if err != nil {
		return err
	}
	defer src.Close()

	byteValue, _ := io.ReadAll(src)

	var jsonSchema schema.JSONSchema
	if err = json.Unmarshal(byteValue, &jsonSchema); err != nil {
		return err
	}

	generator := schema.JSONSchemaToSwiftCodeGenerator{Schema: jsonSchema}

	fParts := utils.SplitBefore(*destFlag, ".")
	if len(fParts) != 2 {
		err := fmt.Errorf("codegen: unable to parse dest flag: %v", *destFlag)
		return err
	}

	out, err := utils.NewTempFile(filepath.Dir(*destFlag), fParts[0], fParts[1])
	if err != nil {
		return err
	}

	code := generator.Generate()
	if err = os.WriteFile(out.Name(), []byte(code), 0644); err != nil {
		return err
	}

	return nil
}
