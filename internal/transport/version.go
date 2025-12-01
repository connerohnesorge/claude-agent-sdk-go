// Package transport provides process management and communication
// for Claude Code subprocesses.
package transport

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/connerohnesorge/claude-agent-sdk-go/pkg/clauderrs"
)

const (
	// MinimumClaudeCodeVersion specifies the minimum required version of the Claude CLI.
	MinimumClaudeCodeVersion = "2.0.0"

	// SkipVersionCheckEnvVar is the environment variable name that, when set to "true",
	// skips the CLI version check.
	SkipVersionCheckEnvVar = "CLAUDE_AGENT_SDK_SKIP_VERSION_CHECK"
)

var (
	// ErrVersionCheckFailed is returned when the version check command fails to execute.
	ErrVersionCheckFailed = errors.New("failed to check Claude CLI version")

	// ErrVersionParseFailed is returned when parsing the version string fails.
	ErrVersionParseFailed = errors.New("failed to parse Claude CLI version")

	// ErrVersionTooOld is returned when the Claude CLI version is below the minimum required.
	ErrVersionTooOld = errors.New("Claude CLI version is below minimum required")
)

// versionRegex matches semantic version patterns like "2.0.0" or "v2.0.0-beta.1".
// It captures the first semver pattern found in the output, which is intentional
// since CLI version output typically places the primary version at or near the start
// (e.g., "claude version 2.0.0" or "v2.0.0-beta.1").
//
// Note: This regex extracts the FIRST occurrence of a semver pattern. In edge cases
// where output contains multiple versions (e.g., "v1.2.3 updated from 0.9.8"),
// the first match (1.2.3) is returned. This is correct for CLI --version output
// which consistently places the current version first.
var versionRegex = regexp.MustCompile(`(\d+)\.(\d+)\.(\d+)`)

// checkCLIVersion verifies that the Claude CLI version meets the minimum requirements.
// It can be skipped by setting the CLAUDE_AGENT_SDK_SKIP_VERSION_CHECK environment
// variable to "true".
func checkCLIVersion(executable string) error {
	// Check if version check should be skipped
	if strings.EqualFold(os.Getenv(SkipVersionCheckEnvVar), "true") {
		return nil
	}

	// Execute claude --version
	cmd := exec.Command(executable, "--version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf(errWrapFormat, ErrVersionCheckFailed, err)
	}

	// Parse the version from output
	major, minor, patch, err := parseVersion(string(output))
	if err != nil {
		return fmt.Errorf(errWrapFormat, ErrVersionParseFailed, err)
	}

	currentVersion := fmt.Sprintf("%d.%d.%d", major, minor, patch)

	// Compare versions
	cmp, err := compareVersions(currentVersion, MinimumClaudeCodeVersion)
	if err != nil {
		return fmt.Errorf(errWrapFormat, ErrVersionParseFailed, err)
	}

	if cmp < 0 {
		return clauderrs.NewVersionMismatchError(currentVersion, MinimumClaudeCodeVersion)
	}

	return nil
}

// parseVersion extracts and parses the semantic version from the CLI output.
// It handles formats like "claude version 2.0.0" or just "2.0.0", and strips
// pre-release suffixes like "-beta.1".
func parseVersion(output string) (major, minor, patch int, err error) {
	matches := versionRegex.FindStringSubmatch(output)
	if len(matches) != 4 {
		return 0, 0, 0, fmt.Errorf("no valid version found in output: %s", output)
	}

	major, err = strconv.Atoi(matches[1])
	if err != nil {
		return 0, 0, 0, fmt.Errorf("invalid major version: %w", err)
	}

	minor, err = strconv.Atoi(matches[2])
	if err != nil {
		return 0, 0, 0, fmt.Errorf("invalid minor version: %w", err)
	}

	patch, err = strconv.Atoi(matches[3])
	if err != nil {
		return 0, 0, 0, fmt.Errorf("invalid patch version: %w", err)
	}

	return major, minor, patch, nil
}

// compareVersions compares two semantic version strings.
// Returns:
//   - 1 if current > minimum
//   - 0 if current == minimum
//   - -1 if current < minimum
func compareVersions(current, minimum string) (int, error) {
	currMajor, currMinor, currPatch, err := parseVersion(current)
	if err != nil {
		return 0, fmt.Errorf("failed to parse current version: %w", err)
	}

	minMajor, minMinor, minPatch, err := parseVersion(minimum)
	if err != nil {
		return 0, fmt.Errorf("failed to parse minimum version: %w", err)
	}

	// Compare major version
	if currMajor > minMajor {
		return 1, nil
	}
	if currMajor < minMajor {
		return -1, nil
	}

	// Major versions equal, compare minor
	if currMinor > minMinor {
		return 1, nil
	}
	if currMinor < minMinor {
		return -1, nil
	}

	// Major and minor equal, compare patch
	if currPatch > minPatch {
		return 1, nil
	}
	if currPatch < minPatch {
		return -1, nil
	}

	// All equal
	return 0, nil
}
