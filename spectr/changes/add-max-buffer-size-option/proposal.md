# Change: Add MaxBufferSize Option for Transport Buffer Control

## Why
The Go SDK is missing the `MaxBufferSize` option that exists in the Python SDK. This option is critical for controlling memory usage when processing large outputs from the Claude CLI, preventing out-of-memory (OOM) errors in constrained environments, and configuring behavior for long-running operations that produce substantial stdout data.

The Python SDK implements this as a configurable limit (default 1MB) that prevents unbounded memory growth when buffering incomplete JSON messages from the CLI's stdout stream. Without this option, the Go SDK cannot handle large outputs safely or provide users with memory usage guarantees.

## What Changes
- Add `MaxBufferSize int` field to the `Options` struct in `pkg/claude/options.go`
- Add documentation explaining the field's purpose, default value, and when to customize it
- Update the transport layer in `internal/transport/transport.go` to:
  - Accept and use the `MaxBufferSize` configuration
  - Implement buffer size checking during JSON message accumulation
  - Return appropriate errors when buffer limits are exceeded
- Add default constant `DefaultMaxBufferSize = 1024 * 1024` (1MB) to match Python SDK
- Add error type in `pkg/clauderrs/` for buffer size exceeded scenarios

## Impact
- **Affected specs**: client-options
- **Affected code**:
  - `pkg/claude/options.go` - Add field to Options struct
  - `internal/transport/transport.go` - Add buffer size enforcement
  - `internal/transport/process.go` - Pass MaxBufferSize to transport
  - `pkg/clauderrs/` - Add new error type for buffer exceeded
- **Breaking changes**: None (additive change, backward compatible)
- **Dependencies**: None
