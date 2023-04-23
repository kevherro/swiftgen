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

// Package driver provides an external entrypoint to the swiftgen driver.
package driver

import internalDriver "github.com/kevherro/swiftgen/internal/driver"

// SwiftGen converts JSON schema data located at src to Swift code and writes
// it to dest.
func SwiftGen() error {
	return internalDriver.SwiftGen()
}
