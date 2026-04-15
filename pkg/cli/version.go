package cli

import (
	"github.com/github/gh-aw/pkg/logger"
	"github.com/github/gh-aw/pkg/workflow"
)

var versionLog = logger.New("cli:version")

// Package-level version information
var (
	version = "dev"
)

func init() {
	// Set the version in the workflow package so NewCompiler() auto-detects it
	workflow.SetVersion(version)
}

// SetVersionInfo sets the version information for the CLI and workflow package
func SetVersionInfo(v string) {
	versionLog.Printf("SetVersionInfo: setting version to %q", v)
	version = v
	workflow.SetVersion(v) // Keep workflow package in sync
}

// GetVersion returns the current version
func GetVersion() string {
	versionLog.Printf("GetVersion: returning %q", version)
	return version
}
