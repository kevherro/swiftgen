#!/usr/bin/env bash

# MIT License
#
# Copyright (c) 2023 Kevin Herro
#
# Permission is hereby granted, free of charge, to any person obtaining a copy
# of this software and associated documentation files (the "Software"), to deal
# in the Software without restriction, including without limitation the rights
# to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
# copies of the Software, and to permit persons to whom the Software is
# furnished to do so, subject to the following conditions:
#
# The above copyright notice and this permission notice shall be included in
# all copies or substantial portions of the Software.

set -e
set -x
MODE=atomic
echo "mode: $MODE" > coverage.txt

if [ "$RUN_STATICCHECK" != "false" ]; then
  staticcheck ./...
fi

# Packages that have any tests.
PKG=$(go list -f '{{if .TestGoFiles}}{{.ImportPath}}{{end}}' ./...)

go test -v "$PKG"

for d in $PKG; do
  go test -race -coverprofile=profile.out -covermode=$MODE "$d"
  if [ -f profile.out ]; then
    grep -vh "^mode: " profile.out >> coverage.txt
    rm profile.out
  fi
done

go vet -all ./...
if [ "$RUN_GOLANGCI_LINTER" != "false" ];  then
  golangci-lint run -D errcheck ./...
fi

gofmt -s -d .
