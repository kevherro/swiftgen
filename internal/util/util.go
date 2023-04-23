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

// Package util provides a set of common functions that can be used across
// different packages.
package util

import "strings"

// SplitBefore slices s into all substrings before the first instance of sep and
// returns a slice of those substrings.
//
// If s does not contain sep and sep is not empty, SplitBefore returns
// a slice of length 1 whose only element is s.
//
// If s or sep are empty, SplitBefore returns an empty slice.
func SplitBefore(s string, sep string) []string {
	if len(sep) == 0 || len(s) == 0 {
		return []string{s}
	}

	var result []string
	index := strings.Index(s, sep)
	if index != -1 {
		if index > 0 {
			result = append(result, s[:index])
		}
		result = append(result, s[index:])
	} else {
		result = append(result, s)
	}

	// s does not contain sep.
	if len(result) == 0 {
		return []string{s}
	}

	return result
}
