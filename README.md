[![Github Action CI](https://github.com/kevherro/swiftgen/workflows/ci/badge.svg)](https://github.com/kevherro/swiftgen/actions)
[![Go Reference](https://pkg.go.dev/badge/github.com/kevherro/swiftgen.svg)](https://pkg.go.dev/github.com/kevherro/swiftgen)

# Introduction

Are you over feeling the pain of manually updating your Swift models
after you find out (the hard way) that your dependency's API changed?

swiftgen is a code generation tool for Swift that can help alleviate this pain
point by automatically generating models based on a specified schema or structure.
This can help ensure that your models are always up-to-date with the latest version of the API,
and can save time and effort that would otherwise be spent manually updating models.

# Building swiftgen

Prerequisites:

- Go development kit of a [supported version](https://golang.org/doc/devel/release.html#policy).
  Follow [these instructions](http://golang.org/doc/code.html) to prepare
  the environment.

To build and install it:

    go install github.com/kevherro/swiftgen@latest

The binary will be installed in `$GOPATH/bin` (`$HOME/go/bin` by default).

# Basic usage

swiftgen can read a JSON schema file (`[json_schema_file].json`).
Specify the input file in the command line. It writes the Swift code (`[generated_swift_code].swift`)
to your current working directory.

### Generate Swift code based on a JSON schema file

```
% swiftgen [json_schema_file]
Where
    json_schema_file: Local path to the JSON schema file.
```
