# Implementation Tasks

## 1. Add MaxBufferSize to Options struct
- [ ] 1.1 Add `MaxBufferSize int` field to `Options` struct in `pkg/claude/options.go`
- [ ] 1.2 Add comprehensive godoc comments explaining:
  - Purpose: Controls max bytes for CLI stdout buffering
  - Default: 1MB (1024 * 1024 bytes)
  - When to customize: Large outputs, constrained memory environments
  - Zero value behavior: Use default
- [ ] 1.3 Add `DefaultMaxBufferSize` constant (1024 * 1024) to `pkg/claude/options.go`

## 2. Add buffer exceeded error type
- [ ] 2.1 Add `ErrBufferSizeExceeded` error type to `pkg/clauderrs/`
- [ ] 2.2 Implement error with context (buffer size, limit exceeded)
- [ ] 2.3 Add godoc documentation for the error type

## 3. Update transport layer
- [ ] 3.1 Modify `NewStdioTransport` to accept `maxBufferSize int` parameter
- [ ] 3.2 Add `maxBufferSize` field to `StdioTransport` struct
- [ ] 3.3 Update `Read()` method to track accumulated buffer size
- [ ] 3.4 Implement buffer size checking when accumulating incomplete JSON
- [ ] 3.5 Return `ErrBufferSizeExceeded` when limit exceeded with details
- [ ] 3.6 Update `internal/transport/process.go` to pass MaxBufferSize from Options

## 4. Update client initialization
- [ ] 4.1 Update `Client.Query()` to pass `MaxBufferSize` (or default) to transport
- [ ] 4.2 Ensure zero value uses `DefaultMaxBufferSize`

## 5. Add tests
- [ ] 5.1 Add unit test for buffer size enforcement (large incomplete JSON)
- [ ] 5.2 Add test for custom buffer size configuration
- [ ] 5.3 Add test for default buffer size behavior
- [ ] 5.4 Add test for error details when buffer exceeded
- [ ] 5.5 Verify tests pass with `go test ./...`

## 6. Update examples (if needed)
- [ ] 6.1 Review existing examples to ensure they work with new option
- [ ] 6.2 Add example comment showing MaxBufferSize usage if warranted
