# Client library and command for Korrel8r

[![Build](https://github.com/korrel8r/client/actions/workflows/build.yml/badge.svg)](https://github.com/korrel8r/client/actions/workflows/build.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/korrel8r/client.svg)](https://pkg.go.dev/github.com/korrel8r/client)

Client library and command for the [Korrel8r](http://github.com/korrel8r/korrel8r) REST API.

Installation:

    go install github.com/korrel8r/client/cmd/korrel8rcli@latest

Command help:

    korrel8rcli help


**NOTE**: This code _must not_ depend on `github.com/korrel8r/korrel8r`.
It is deliberately separated from the korrel8r to keep its dependency set small,
to avoid potential package clashes when included in 3rd party applications.
