# Implementation Tasks

## 1. Add User Field to Options Struct
- [x] 1.1 Add `User string` field to `Options` struct in `pkg/claude/options.go`
- [x] 1.2 Add godoc comment explaining the field's purpose and platform limitations
- [x] 1.3 Document that this feature is Unix-specific and requires appropriate permissions

## 2. Update Transport Layer Configuration
- [x] 2.1 Add `User string` field to `ProcessConfig` struct in `internal/transport/process.go`
- [x] 2.2 Implement `resolveUserCredential(username string) (*syscall.Credential, error)` helper function
- [x] 2.3 Update `createCommand` function to configure `SysProcAttr.Credential` when User is specified
- [x] 2.4 Add proper error handling for user lookup failures

## 3. Wire User Option Through Client
- [x] 3.1 Update client code to pass `User` field from `Options` to `ProcessConfig`
- [x] 3.2 Ensure the User value flows through to process creation

## 4. Add Error Handling
- [x] 4.1 Define new error types in `internal/transport/errors.go` for user resolution failures
- [x] 4.2 Add validation for empty user strings (handled via nil return in resolveUserCredential)
- [x] 4.3 Handle permission errors when switching to non-privileged users (wrapped errors)

## 5. Documentation and Examples
- [x] 5.1 Document security considerations in godoc (comprehensive docs in options.go)
- [x] 5.2 Note platform-specific behavior (Unix-only) - documented in Options.User field
- [ ] 5.3 Add example demonstrating User option in `examples/` directory (optional - deferred)

## 6. Testing
- [x] 6.1 Add unit tests for user credential resolution
- [x] 6.2 Add test cases for invalid usernames
- [x] 6.3 Add test cases for current user resolution
- [x] 6.4 Add test cases for root user resolution
- [x] 6.5 Verify tests pass on Unix systems
- [ ] 6.6 Integration tests verifying process runs as specified user (requires root - deferred)

## 7. Validation
- [x] 7.1 Run `go test ./...` to verify all tests pass
- [x] 7.2 Run `go build ./...` to verify compilation
- [ ] 7.3 Run `golangci-lint run` to ensure linting passes (optional)
- [ ] 7.4 Test on multiple Unix platforms (Linux, macOS) if possible (optional)
