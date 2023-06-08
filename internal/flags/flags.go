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

// Package flags manages command-line flags for swiftgen.
package flags

import (
	"flag"
	"fmt"
)

type CmdFlags struct {
	Src  string // Path to the JSON schema file
	Dest string // Path to the destination Swift file
}

// Parse parses flag definitions, which should not include the command name.
func (cf *CmdFlags) Parse() {
	flag.StringVar(&cf.Src, "src", "", "Path to the JSON schema file")
	flag.StringVar(&cf.Dest, "dest", "", "Path to the destination Swift file")

	flag.Parse()

	if cf.Src == "" || cf.Dest == "" {
		flag.Usage()
		panic("both -src and -dest flags are required")
	}
}

// Validate validates the flags provided by the user.
func (cf *CmdFlags) Validate() error {
	if cf.Src == "" {
		return fmt.Errorf("src cannot be empty")
	}
	if cf.Dest == "" {
		return fmt.Errorf("dest cannot be empty")
	}
	return nil
}
