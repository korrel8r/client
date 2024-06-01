package build

import (
	_ "embed"
	"strings"
)

var (
	//go:embed version.txt
	version string
	// Version of the current build, a semver string.
	Version = strings.TrimSpace(version)
)
