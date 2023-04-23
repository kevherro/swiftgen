#!/usr/bin/env bash

set -e
set -x
MODE=atomic
echo "mode: $MODE" > coverage.txt

if [ "$RUN_STATICCHECK" != "false" ]; then
  staticcheck ./...
fi

# Packages that have any tests.
PKG=$(go list -f '{{if .TestGoFiles}}{{.ImportPath}}{{end}}' ./...)

go test -v $PKG

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
