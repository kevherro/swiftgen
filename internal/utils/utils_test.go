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

package utils

import (
	"reflect"
	"testing"
)

func TestSplitBefore(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		sep      string
		expected []string
	}{
		{"Basic", "Schema.swift", ".", []string{"Schema", ".swift"}},
		{"No separator", "Schema", ".", []string{"Schema"}},
		{"Separator not found", "Schema.swift", "*", []string{"Schema.swift"}},
		{"Empty input", "", ".", []string{""}},
		{"Empty separator", "Schema.swift", "", []string{"Schema.swift"}},
		{"Empty input and separator", "", "", []string{""}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SplitBefore(tt.input, tt.sep)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("SplitBefore(%q, %q) = %v, want %v", tt.input, tt.sep, result, tt.expected)
			}
		})
	}
}
