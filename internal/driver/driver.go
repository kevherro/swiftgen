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

// Package driver implements the core swiftgen functionality.
package driver

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/kevherro/swiftgen/internal/schema"
)

// SwiftGen converts JSON schema data located at src to Swift code and writes
// it to dest.
func SwiftGen() error {
	srcFlag := flag.String("src", "", "Path to JSON schema file")
	destFlag := flag.String("dest", "", "Path to write location")

	flag.Parse()

	if *srcFlag == "" {
		flag.Usage()
		err := fmt.Errorf("driver: --src flag value not provided")
		return err
	}

	src, err := os.Open(*srcFlag)
	if err != nil {
		return err
	}
	defer src.Close()

	b, err := io.ReadAll(src)
	if err != nil {
		return err
	}

	var s schema.JSONSchema
	if err = json.Unmarshal(b, &s); err != nil {
		return err
	}

	if err = os.MkdirAll(filepath.Dir(*destFlag), 0750); err != nil {
		return err
	}

	g := schema.JSONSchemaToSwiftCodeGenerator{Schema: s}
	c := g.Generate()
	if err = os.WriteFile(*destFlag, []byte(c), 0666); err != nil {
		return err
	}

	return nil
}
