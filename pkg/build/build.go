package build

import "strings"

var (
	version string //+embed:version.txt
	// Version of the current build, a semver string.
	Version = strings.TrimSpace(version)
)
