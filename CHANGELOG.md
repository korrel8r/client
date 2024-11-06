# Change Log

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/)
and this project adheres to [Semantic Versioning](http://semver.org/).

## [0.0.4] - 2024-05-30

### Added
- Updated for API changes in korrel8r 0.7.3.
- feat: add tooltips to browser nodes.
- feat: Take bearer token from kube config by default.
- feat: --bearer-token option to set Authorization headers.
- feat: --insecure option, also fix arguments to goals command.

## [0.0.1] - 2024-05-30

### Added

First version of `korrel8rcli` command.
 - REST client, command line access to a remote korrel8r server. See `korrel8rcli --help`
 - Web browser API using data from remote korrel8r server, see `korrel8rcli web --help`
 - Client packages for 3rd party use, see ./pkg/swagger

