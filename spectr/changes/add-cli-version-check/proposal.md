# Change: Add CLI Version Checking

## Why
The Go SDK currently lacks version compatibility validation with the Claude Code CLI, which can lead to cryptic failures when protocol mismatches occur. The Python and TypeScript SDKs already implement version checking at connection time to ensure compatibility and provide clear error messages when the CLI is outdated. This creates an inconsistency across SDK implementations and leaves Go users vulnerable to confusing runtime errors.

Without version checking, users experience:
- Cryptic protocol errors when CLI is too old
- Wasted debugging time trying to understand failures
- No clear guidance on what minimum CLI version is required
- Inconsistent behavior compared to other official SDKs

## What Changes
- Add CLI version validation at process startup in the transport layer
- Add minimum supported CLI version constant (matching Python SDK: 2.0.0)
- Add environment variable `CLAUDE_AGENT_SDK_SKIP_VERSION_CHECK` to bypass validation
- Add new error type `ErrCodeVersionMismatch` to the clauderrs package
- Add version checking logic that executes `claude --version` before establishing connection
- Update client initialization to perform version check unless explicitly skipped

The version check will:
1. Execute `claude --version` to get the installed CLI version
2. Parse and compare against minimum required version (2.0.0)
3. Return a clear error if CLI is too old
4. Skip check if `CLAUDE_AGENT_SDK_SKIP_VERSION_CHECK` environment variable is set to `true`

## Impact
- **Affected specs**: client-options (new CLI version check option)
- **Affected code**:
  - `internal/transport/process.go` - Add version check before process start
  - `internal/transport/version.go` - New file for version checking logic
  - `pkg/clauderrs/types.go` - Add new error code for version mismatch
  - `pkg/claude/options.go` - Document version check behavior
  - Integration tests to verify version checking works correctly

**Breaking Changes**: None. This is purely additive functionality with a skip mechanism for edge cases.

**Benefits**:
- Clearer error messages when CLI is incompatible
- Better user experience matching other SDK implementations
- Prevents wasted time debugging protocol mismatches
- Provides clear upgrade path when CLI is outdated
