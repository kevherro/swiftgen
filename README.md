[![Github Action CI](https://github.com/kevherro/swiftgen/workflows/ci/badge.svg)](https://github.com/kevherro/swiftgen/actions)
[![Codecov](https://codecov.io/gh/kevherro/swiftgen/graph/badge.svg)](https://codecov.io/gh/kevherro/swiftgen)
[![Go Reference](https://pkg.go.dev/badge/github.com/kevherro/swiftgen.svg)](https://pkg.go.dev/github.com/kevherro/swiftgen)

# Introduction

Are you over feeling the pain of manually updating your Swift models
after finding out (the hard way) that your dependency's JSON schema changed?

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

swiftgen can read a JSON schema file and convert it to Swift structs.

### Generate Swift code based on a JSON schema file

```
% swiftgen --src <json_schema_file> --dest <swift_file>
Where
    json_schema_file: Path to the JSON schema file. Required.
    swift_file: Path to the destination Swift file. Required.
```

---

swiftgen project structure is heavily influenced by [pprof](https://github.com/google/pprof). Thanks y'all!
