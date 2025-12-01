## 1. Add Error Code for Version Mismatch
- [ ] 1.1 Add `ErrCodeVersionMismatch` constant to `pkg/clauderrs/types.go`
- [ ] 1.2 Add helper function `NewVersionMismatchError` in `pkg/clauderrs/client.go`
- [ ] 1.3 Document error code with example usage

## 2. Implement Version Checking Logic
- [ ] 2.1 Create `internal/transport/version.go` file
- [ ] 2.2 Add `MinimumClaudeCodeVersion` constant set to "2.0.0"
- [ ] 2.3 Implement `checkCLIVersion()` function that:
  - [ ] 2.3.1 Checks `CLAUDE_AGENT_SDK_SKIP_VERSION_CHECK` environment variable
  - [ ] 2.3.2 Executes `claude --version` command
  - [ ] 2.3.3 Parses version output (handles formats like "claude version 2.0.0")
  - [ ] 2.3.4 Compares using semantic versioning rules
  - [ ] 2.3.5 Returns appropriate error with version details
- [ ] 2.4 Add helper function `parseVersion()` for version string parsing
- [ ] 2.5 Add helper function `compareVersions()` for semantic version comparison

## 3. Integrate Version Check into Process Startup
- [ ] 3.1 Modify `NewProcess()` in `internal/transport/process.go` to call version check
- [ ] 3.2 Ensure version check happens before process spawn
- [ ] 3.3 Handle version check errors appropriately
- [ ] 3.4 Add godoc comments explaining version check behavior

## 4. Add Unit Tests
- [ ] 4.1 Test version parsing with various formats:
  - [ ] 4.1.1 "claude version 2.0.0"
  - [ ] 4.1.2 "2.1.5"
  - [ ] 4.1.3 "v2.0.0-beta.1"
  - [ ] 4.1.4 Invalid formats
- [ ] 4.2 Test version comparison logic:
  - [ ] 4.2.1 Equal versions (2.0.0 == 2.0.0)
  - [ ] 4.2.2 Greater versions (2.1.0 > 2.0.0)
  - [ ] 4.2.3 Lesser versions (1.9.0 < 2.0.0)
  - [ ] 4.2.4 Patch versions (2.0.1 > 2.0.0)
- [ ] 4.3 Test environment variable skip mechanism:
  - [ ] 4.3.1 Variable set to "true"
  - [ ] 4.3.2 Variable set to "True" (case insensitive)
  - [ ] 4.3.3 Variable set to "false"
  - [ ] 4.3.4 Variable not set
- [ ] 4.4 Test error cases:
  - [ ] 4.4.1 CLI not found
  - [ ] 4.4.2 Version command fails
  - [ ] 4.4.3 Invalid version format
  - [ ] 4.4.4 Version below minimum

## 5. Add Integration Tests
- [ ] 5.1 Create `test/integration/version_check_test.go`
- [ ] 5.2 Test with mock CLI executable that returns specific versions
- [ ] 5.3 Test skip mechanism with environment variable set
- [ ] 5.4 Verify error messages contain expected version information

## 6. Update Documentation
- [ ] 6.1 Add version checking section to main README.md
- [ ] 6.2 Document `CLAUDE_AGENT_SDK_SKIP_VERSION_CHECK` environment variable
- [ ] 6.3 Add godoc comments to all new public constants and functions
- [ ] 6.4 Update package-level documentation in `internal/transport/doc.go`

## 7. Validation
- [ ] 7.1 Run all unit tests and ensure they pass
- [ ] 7.2 Run integration tests with real Claude CLI
- [ ] 7.3 Run golangci-lint and fix any issues
- [ ] 7.4 Verify examples still work with version checking enabled
- [ ] 7.5 Test with environment variable to skip version check
- [ ] 7.6 Verify error messages are clear and actionable
