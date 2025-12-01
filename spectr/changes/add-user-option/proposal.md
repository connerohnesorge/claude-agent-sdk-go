# Change: Add User Option for Process Isolation

## Why
The Go SDK currently lacks the `User` field present in the Python SDK (`ClaudeAgentOptions.user`), preventing users from running the Claude Code CLI subprocess as a different user. This capability is essential for security isolation, running with reduced privileges, multi-tenant environments, and containerized deployments where processes should run as non-root users.

## What Changes
- Add `User string` field to the `Options` struct in `pkg/claude/options.go`
- Update `ProcessConfig` in `internal/transport/process.go` to include a `User` field
- Implement user credential resolution and process attribute configuration in the transport layer
- Add documentation explaining the security implications and platform limitations
- Create tests verifying user switching functionality on Unix-like systems

## Impact
- Affected specs: `client-options`
- Affected code:
  - `pkg/claude/options.go` - Add `User` field to Options struct
  - `internal/transport/process.go` - Add `User` field to ProcessConfig and implement credential setting via `syscall.SysProcAttr.Credential`
  - `pkg/claude/client.go` - Pass User option through to transport layer
- Platform consideration: This feature uses Unix-specific `syscall.SysProcAttr.Credential`, which requires user lookup and UID/GID resolution
- Breaking changes: None (additive change only)
- Python SDK parity: Achieves parity with Python SDK's `user` option
