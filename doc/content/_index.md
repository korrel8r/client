---
title: Korrel8r Client
weight: 1
hideChildren: true
---

Client library and command for the [Korrel8r](https://korrel8r.github.io/korrel8r/) REST API.

## Installation

```terminal
go install github.com/korrel8r/client/cmd/korrel8rcli@latest
korrel8rcli --help
```

## Usage

`korrel8rcli` is a command-line client for the Korrel8r REST API.
It automatically uses your cluster login credentials for authentication.

```terminal
# List available domains
korrel8rcli -u $KORREL8R_URL domains

# Find all data related to a deployment
korrel8rcli -u $KORREL8R_URL neighbors --query 'k8s:Deployment:{namespace: korrel8r}'

# Find logs related to a deployment
korrel8rcli -u $KORREL8R_URL goals --start 'k8s:Deployment:{namespace: korrel8r}' --goal 'log:application'
```

