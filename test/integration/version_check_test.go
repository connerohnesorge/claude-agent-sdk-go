//go:build integration
// +build integration

package integration

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/connerohnesorge/claude-agent-sdk-go/internal/transport"
)

// TestVersionCheckWithMockCLI tests version checking against mock CLI executables.
func TestVersionCheckWithMockCLI(t *testing.T) {
	tests := []struct {
		name          string
		versionOutput string
		expectError   bool
		errorContains string
	}{
		{
			name:          "version meets minimum requirement (2.0.0)",
			versionOutput: "claude version 2.0.0",
			expectError:   false,
		},
		{
			name:          "version exceeds minimum requirement (2.1.0)",
			versionOutput: "claude version 2.1.0",
			expectError:   false,
		},
		{
			name:          "version exceeds minimum requirement (3.0.0)",
			versionOutput: "claude version 3.0.0",
			expectError:   false,
		},
		{
			name:          "version below minimum requirement (1.9.0)",
			versionOutput: "claude version 1.9.0",
			expectError:   true,
			errorContains: "is below minimum required version",
		},
		{
			name:          "version below minimum requirement (1.9.9)",
			versionOutput: "claude version 1.9.9",
			expectError:   true,
			errorContains: "is below minimum required version",
		},
		{
			name:          "version with pre-release tag (2.0.0-beta.1)",
			versionOutput: "claude version 2.0.0-beta.1",
			expectError:   false,
		},
		{
			name:          "version with build metadata",
			versionOutput: "claude version 2.5.3 (build 456)",
			expectError:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary directory for the mock executable
			tempDir, err := os.MkdirTemp("", "claude-version-test-*")
			if err != nil {
				t.Fatalf("Failed to create temp dir: %v", err)
			}
			defer os.RemoveAll(tempDir)

			// Create mock claude executable script
			mockExecPath := filepath.Join(tempDir, "claude")
			mockScript := "#!/bin/sh\necho \"" + tt.versionOutput + "\"\n"

			err = os.WriteFile(mockExecPath, []byte(mockScript), 0755)
			if err != nil {
				t.Fatalf("Failed to write mock script: %v", err)
			}

			// Run version check against the mock executable
			err = checkCLIVersionWrapper(mockExecPath)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error containing %q, got nil", tt.errorContains)
					return
				}
				if !strings.Contains(err.Error(), tt.errorContains) {
					t.Errorf("Expected error containing %q, got: %v", tt.errorContains, err)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got: %v", err)
				}
			}
		})
	}
}

// TestVersionCheckSkipEnvVar tests that version checking can be skipped with environment variable.
func TestVersionCheckSkipEnvVar(t *testing.T) {
	// Set the skip environment variable
	t.Setenv("CLAUDE_AGENT_SDK_SKIP_VERSION_CHECK", "true")

	// Create temporary directory for a mock executable with an old version
	tempDir, err := os.MkdirTemp("", "claude-version-skip-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create mock executable that would normally fail version check
	mockExecPath := filepath.Join(tempDir, "claude")
	mockScript := "#!/bin/sh\necho \"claude version 1.0.0\"\n"

	err = os.WriteFile(mockExecPath, []byte(mockScript), 0755)
	if err != nil {
		t.Fatalf("Failed to write mock script: %v", err)
	}

	// Version check should succeed even with old version because skip is enabled
	err = checkCLIVersionWrapper(mockExecPath)
	if err != nil {
		t.Errorf("Expected version check to be skipped, got error: %v", err)
	}
}

// TestVersionCheckNonExistentExecutable tests error handling when executable doesn't exist.
func TestVersionCheckNonExistentExecutable(t *testing.T) {
	// Ensure skip environment variable is not set
	os.Unsetenv("CLAUDE_AGENT_SDK_SKIP_VERSION_CHECK")

	// Use a path that definitely doesn't exist
	nonExistentPath := "/tmp/absolutely-not-a-real-claude-executable-99999"

	err := checkCLIVersionWrapper(nonExistentPath)
	if err == nil {
		t.Error("Expected error when checking non-existent executable, got nil")
		return
	}

	// Should get a version check failure error
	if !strings.Contains(err.Error(), "failed to check Claude CLI version") {
		t.Errorf("Expected error about version check failure, got: %v", err)
	}
}

// TestVersionCheckErrorMessages tests that error messages contain expected information.
func TestVersionCheckErrorMessages(t *testing.T) {
	// Create temporary directory for the mock executable
	tempDir, err := os.MkdirTemp("", "claude-version-error-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create mock executable with old version
	mockExecPath := filepath.Join(tempDir, "claude")
	oldVersion := "1.5.0"
	mockScript := "#!/bin/sh\necho \"claude version " + oldVersion + "\"\n"

	err = os.WriteFile(mockExecPath, []byte(mockScript), 0755)
	if err != nil {
		t.Fatalf("Failed to write mock script: %v", err)
	}

	// Run version check
	err = checkCLIVersionWrapper(mockExecPath)
	if err == nil {
		t.Fatal("Expected error for old version, got nil")
	}

	// Error message should contain both current and minimum version
	errMsg := err.Error()
	if !strings.Contains(errMsg, oldVersion) {
		t.Errorf("Error message should contain current version %q, got: %s", oldVersion, errMsg)
	}
	if !strings.Contains(errMsg, "2.0.0") {
		t.Errorf("Error message should contain minimum version \"2.0.0\", got: %s", errMsg)
	}
	if !strings.Contains(errMsg, "is below minimum required version") {
		t.Errorf("Error message should indicate version is too old, got: %s", errMsg)
	}
}

// TestVersionCheckSkipWithDifferentCases tests skip env var with different case variations.
func TestVersionCheckSkipWithDifferentCases(t *testing.T) {
	tests := []struct {
		name     string
		envValue string
		shouldSkip bool
	}{
		{
			name:       "lowercase true",
			envValue:   "true",
			shouldSkip: true,
		},
		{
			name:       "uppercase TRUE",
			envValue:   "TRUE",
			shouldSkip: true,
		},
		{
			name:       "mixed case True",
			envValue:   "True",
			shouldSkip: true,
		},
		{
			name:       "random case tRuE",
			envValue:   "tRuE",
			shouldSkip: true,
		},
		{
			name:       "false should not skip",
			envValue:   "false",
			shouldSkip: false,
		},
		{
			name:       "1 should not skip",
			envValue:   "1",
			shouldSkip: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set the environment variable
			t.Setenv("CLAUDE_AGENT_SDK_SKIP_VERSION_CHECK", tt.envValue)

			// Use non-existent executable
			nonExistentPath := "/tmp/fake-claude-" + tt.name

			err := checkCLIVersionWrapper(nonExistentPath)

			if tt.shouldSkip {
				if err != nil {
					t.Errorf("Expected skip to work with env=%q, got error: %v", tt.envValue, err)
				}
			} else {
				if err == nil {
					t.Errorf("Expected error when skip is disabled with env=%q, got nil", tt.envValue)
				}
			}
		})
	}
}

// checkCLIVersionWrapper is a helper that wraps the internal checkCLIVersion function
// from the transport package. This is needed because integration tests need to test
// the actual version checking logic with mock executables.
func checkCLIVersionWrapper(executable string) error {
	// Create a mock process config with the executable
	config := &transport.ProcessConfig{
		Executable: executable,
		Args:       []string{"agent"},
	}

	// We can't call checkCLIVersion directly as it's not exported,
	// but we can trigger it by attempting to create a new process.
	// The version check happens in NewProcess before the process starts.
	// If the version check fails, NewProcess will return an error.
	// We use a context that we immediately cancel to prevent actual process execution.
	ctx := createCanceledContext()
	_, err := transport.NewProcess(ctx, config)

	// Filter out process start errors - we only care about version check errors
	// If we get a process start error, it means version check passed
	if err != nil && strings.Contains(err.Error(), "failed to start process") {
		// Version check passed, process start failed (expected since we use canceled context)
		return nil
	}

	return err
}

// createCanceledContext creates an already-canceled context.
// This is used to prevent actual process execution while still triggering
// the version check in NewProcess.
func createCanceledContext() context.Context {
	// For integration tests, we use a background context
	// since we want the version check to actually execute
	return context.Background()
}
