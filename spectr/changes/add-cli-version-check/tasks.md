## 1. Add Error Code for Version Mismatch
- [x] 1.1 Add `ErrCodeVersionMismatch` constant to `pkg/clauderrs/types.go`
- [x] 1.2 Add helper function `NewVersionMismatchError` in `pkg/clauderrs/client.go`
- [x] 1.3 Document error code with example usage

## 2. Implement Version Checking Logic
- [x] 2.1 Create `internal/transport/version.go` file
- [x] 2.2 Add `MinimumClaudeCodeVersion` constant set to "2.0.0"
- [x] 2.3 Implement `checkCLIVersion()` function that:
  - [x] 2.3.1 Checks `CLAUDE_AGENT_SDK_SKIP_VERSION_CHECK` environment variable
  - [x] 2.3.2 Executes `claude --version` command
  - [x] 2.3.3 Parses version output (handles formats like "claude version 2.0.0")
  - [x] 2.3.4 Compares using semantic versioning rules
  - [x] 2.3.5 Returns appropriate error with version details
- [x] 2.4 Add helper function `parseVersion()` for version string parsing
- [x] 2.5 Add helper function `compareVersions()` for semantic version comparison

## 3. Integrate Version Check into Process Startup
- [x] 3.1 Modify `NewProcess()` in `internal/transport/process.go` to call version check
- [x] 3.2 Ensure version check happens before process spawn
- [x] 3.3 Handle version check errors appropriately
- [x] 3.4 Add godoc comments explaining version check behavior

## 4. Add Unit Tests
- [x] 4.1 Test version parsing with various formats:
  - [x] 4.1.1 "claude version 2.0.0"
  - [x] 4.1.2 "2.1.5"
  - [x] 4.1.3 "v2.0.0-beta.1"
  - [x] 4.1.4 Invalid formats
- [x] 4.2 Test version comparison logic:
  - [x] 4.2.1 Equal versions (2.0.0 == 2.0.0)
  - [x] 4.2.2 Greater versions (2.1.0 > 2.0.0)
  - [x] 4.2.3 Lesser versions (1.9.0 < 2.0.0)
  - [x] 4.2.4 Patch versions (2.0.1 > 2.0.0)
- [x] 4.3 Test environment variable skip mechanism:
  - [x] 4.3.1 Variable set to "true"
  - [x] 4.3.2 Variable set to "True" (case insensitive)
  - [x] 4.3.3 Variable set to "false"
  - [x] 4.3.4 Variable not set
- [x] 4.4 Test error cases:
  - [x] 4.4.1 CLI not found
  - [x] 4.4.2 Version command fails
  - [x] 4.4.3 Invalid version format
  - [x] 4.4.4 Version below minimum

## 5. Add Integration Tests
- [x] 5.1 Create `test/integration/version_check_test.go`
- [x] 5.2 Test with mock CLI executable that returns specific versions
- [x] 5.3 Test skip mechanism with environment variable set
- [x] 5.4 Verify error messages contain expected version information

## 6. Update Documentation
- [x] 6.1 Add version checking section to main README.md
- [x] 6.2 Document `CLAUDE_AGENT_SDK_SKIP_VERSION_CHECK` environment variable
- [x] 6.3 Add godoc comments to all new public constants and functions
- [x] 6.4 Update package-level documentation in `internal/transport/doc.go`

## 7. Validation
- [x] 7.1 Run all unit tests and ensure they pass
- [x] 7.2 Run integration tests with real Claude CLI
- [x] 7.3 Run golangci-lint and fix any issues
- [x] 7.4 Verify examples still work with version checking enabled
- [x] 7.5 Test with environment variable to skip version check
- [x] 7.6 Verify error messages are clear and actionable
