package transport

import (
	"testing"
)

// TestParseVersion tests the parseVersion function with various version string formats.
func TestParseVersion(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantMajor   int
		wantMinor   int
		wantPatch   int
		wantErr     bool
		errContains string
	}{
		{
			name:      "standard claude version output",
			input:     "claude version 2.0.0",
			wantMajor: 2,
			wantMinor: 0,
			wantPatch: 0,
			wantErr:   false,
		},
		{
			name:      "plain semantic version",
			input:     "2.1.5",
			wantMajor: 2,
			wantMinor: 1,
			wantPatch: 5,
			wantErr:   false,
		},
		{
			name:      "version with pre-release tag",
			input:     "v2.0.0-beta.1",
			wantMajor: 2,
			wantMinor: 0,
			wantPatch: 0,
			wantErr:   false,
		},
		{
			name:      "version with build metadata",
			input:     "claude 2.3.4 (build 123)",
			wantMajor: 2,
			wantMinor: 3,
			wantPatch: 4,
			wantErr:   false,
		},
		{
			name:        "no version in string",
			input:       "no version here",
			wantErr:     true,
			errContains: "no valid version found",
		},
		{
			name:        "empty string",
			input:       "",
			wantErr:     true,
			errContains: "no valid version found",
		},
		{
			name:      "version with v prefix",
			input:     "v3.2.1",
			wantMajor: 3,
			wantMinor: 2,
			wantPatch: 1,
			wantErr:   false,
		},
		{
			name:      "version with extra text before",
			input:     "Claude Code version 1.9.8 is installed",
			wantMajor: 1,
			wantMinor: 9,
			wantPatch: 8,
			wantErr:   false,
		},
		{
			name:      "double digit version numbers",
			input:     "10.20.30",
			wantMajor: 10,
			wantMinor: 20,
			wantPatch: 30,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			major, minor, patch, err := parseVersion(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Errorf("parseVersion(%q) expected error containing %q, got nil", tt.input, tt.errContains)
				} else if tt.errContains != "" && !containsString(err.Error(), tt.errContains) {
					t.Errorf("parseVersion(%q) error = %v, want error containing %q", tt.input, err, tt.errContains)
				}
				return
			}

			if err != nil {
				t.Errorf("parseVersion(%q) unexpected error: %v", tt.input, err)
				return
			}

			if major != tt.wantMajor {
				t.Errorf("parseVersion(%q) major = %d, want %d", tt.input, major, tt.wantMajor)
			}
			if minor != tt.wantMinor {
				t.Errorf("parseVersion(%q) minor = %d, want %d", tt.input, minor, tt.wantMinor)
			}
			if patch != tt.wantPatch {
				t.Errorf("parseVersion(%q) patch = %d, want %d", tt.input, patch, tt.wantPatch)
			}
		})
	}
}

// TestCompareVersions tests semantic version comparison logic.
func TestCompareVersions(t *testing.T) {
	tests := []struct {
		name    string
		current string
		minimum string
		want    int // -1: current < minimum, 0: equal, 1: current > minimum
		wantErr bool
	}{
		{
			name:    "versions are equal",
			current: "2.0.0",
			minimum: "2.0.0",
			want:    0,
			wantErr: false,
		},
		{
			name:    "current version greater - minor",
			current: "2.1.0",
			minimum: "2.0.0",
			want:    1,
			wantErr: false,
		},
		{
			name:    "current version lesser - major",
			current: "1.9.0",
			minimum: "2.0.0",
			want:    -1,
			wantErr: false,
		},
		{
			name:    "current version greater - patch",
			current: "2.0.1",
			minimum: "2.0.0",
			want:    1,
			wantErr: false,
		},
		{
			name:    "current version greater - major",
			current: "3.0.0",
			minimum: "2.9.9",
			want:    1,
			wantErr: false,
		},
		{
			name:    "current version greater - minor with double digits",
			current: "2.10.0",
			minimum: "2.9.0",
			want:    1,
			wantErr: false,
		},
		{
			name:    "current version lesser - patch",
			current: "2.0.0",
			minimum: "2.0.1",
			want:    -1,
			wantErr: false,
		},
		{
			name:    "current version lesser - minor",
			current: "2.5.0",
			minimum: "2.6.0",
			want:    -1,
			wantErr: false,
		},
		{
			name:    "complex comparison - current greater",
			current: "3.2.1",
			minimum: "2.9.8",
			want:    1,
			wantErr: false,
		},
		{
			name:    "complex comparison - current lesser",
			current: "1.9.9",
			minimum: "2.0.0",
			want:    -1,
			wantErr: false,
		},
		{
			name:    "invalid current version",
			current: "invalid",
			minimum: "2.0.0",
			wantErr: true,
		},
		{
			name:    "invalid minimum version",
			current: "2.0.0",
			minimum: "not-a-version",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := compareVersions(tt.current, tt.minimum)

			if tt.wantErr {
				if err == nil {
					t.Errorf("compareVersions(%q, %q) expected error, got nil", tt.current, tt.minimum)
				}
				return
			}

			if err != nil {
				t.Errorf("compareVersions(%q, %q) unexpected error: %v", tt.current, tt.minimum, err)
				return
			}

			if got != tt.want {
				t.Errorf("compareVersions(%q, %q) = %d, want %d", tt.current, tt.minimum, got, tt.want)
			}
		})
	}
}

// TestCheckCLIVersionSkip tests that the version check can be skipped via environment variable.
func TestCheckCLIVersionSkip(t *testing.T) {
	tests := []struct {
		name       string
		envValue   string
		shouldSkip bool
	}{
		{
			name:       "skip with true lowercase",
			envValue:   "true",
			shouldSkip: true,
		},
		{
			name:       "skip with True mixed case",
			envValue:   "True",
			shouldSkip: true,
		},
		{
			name:       "skip with TRUE uppercase",
			envValue:   "TRUE",
			shouldSkip: true,
		},
		{
			name:       "skip with tRuE random case",
			envValue:   "tRuE",
			shouldSkip: true,
		},
		{
			name:       "no skip with false",
			envValue:   "false",
			shouldSkip: false,
		},
		{
			name:       "no skip with empty string",
			envValue:   "",
			shouldSkip: false,
		},
		{
			name:       "no skip with random value",
			envValue:   "yes",
			shouldSkip: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set the environment variable for this test
			if tt.envValue != "" {
				t.Setenv(SkipVersionCheckEnvVar, tt.envValue)
			}

			// Use a non-existent executable path to test that skip works
			// If skip is enabled, it should return nil without trying to execute
			// If skip is disabled, it will try to execute and fail
			nonExistentExec := "/tmp/definitely-not-a-real-claude-executable-12345"
			err := checkCLIVersion(nonExistentExec)

			if tt.shouldSkip {
				// When skip is enabled, should return nil even with non-existent executable
				if err != nil {
					t.Errorf("checkCLIVersion() with skip enabled should return nil, got error: %v", err)
				}
			} else {
				// When skip is disabled, should try to execute and fail with non-existent executable
				if err == nil {
					t.Errorf("checkCLIVersion() with skip disabled and non-existent executable should return error, got nil")
				}
			}
		})
	}
}

// TestCheckCLIVersionUnset tests the behavior when the environment variable is not set.
func TestCheckCLIVersionUnset(t *testing.T) {
	// When the environment variable is not set, the version check should be performed.
	// Using a non-existent executable should cause an error.
	nonExistentExec := "/tmp/definitely-not-a-real-claude-executable-99999"
	err := checkCLIVersion(nonExistentExec)

	if err == nil {
		t.Error("checkCLIVersion() with unset env var and non-existent executable should return error, got nil")
	}

	// The error should be a version check failure
	if !containsString(err.Error(), "failed to check Claude CLI version") {
		t.Errorf("checkCLIVersion() error should mention version check failure, got: %v", err)
	}
}

// containsString checks if a string contains a substring (case-sensitive).
func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		func() bool {
			for i := 0; i <= len(s)-len(substr); i++ {
				if s[i:i+len(substr)] == substr {
					return true
				}
			}
			return false
		}())
}
