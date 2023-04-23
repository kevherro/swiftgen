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

// swiftgen is a tool for generating Swift code based on JSON schema data.
package main

import (
	"fmt"
	"os"

	"github.com/kevherro/swiftgen/driver"
)

func main() {
	if err := driver.SwiftGen(); err != nil {
		fmt.Fprintf(os.Stderr, "swiftgen: %v", err)
		os.Exit(2)
	}
}
