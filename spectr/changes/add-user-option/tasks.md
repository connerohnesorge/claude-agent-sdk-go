# Implementation Tasks

## 1. Add User Field to Options Struct
- [ ] 1.1 Add `User string` field to `Options` struct in `pkg/claude/options.go`
- [ ] 1.2 Add godoc comment explaining the field's purpose and platform limitations
- [ ] 1.3 Document that this feature is Unix-specific and requires appropriate permissions

## 2. Update Transport Layer Configuration
- [ ] 2.1 Add `User string` field to `ProcessConfig` struct in `internal/transport/process.go`
- [ ] 2.2 Implement `resolveUserCredential(username string) (*syscall.Credential, error)` helper function
- [ ] 2.3 Update `createCommand` function to configure `SysProcAttr.Credential` when User is specified
- [ ] 2.4 Add proper error handling for user lookup failures

## 3. Wire User Option Through Client
- [ ] 3.1 Update client code to pass `User` field from `Options` to `ProcessConfig`
- [ ] 3.2 Ensure the User value flows through to process creation

## 4. Add Error Handling
- [ ] 4.1 Define new error types in `pkg/clauderrs/` for user resolution failures
- [ ] 4.2 Add validation for empty user strings
- [ ] 4.3 Handle permission errors when switching to non-privileged users

## 5. Documentation and Examples
- [ ] 5.1 Add example demonstrating User option in `examples/` directory
- [ ] 5.2 Document security considerations in godoc
- [ ] 5.3 Note platform-specific behavior (Unix-only)

## 6. Testing
- [ ] 6.1 Add unit tests for user credential resolution
- [ ] 6.2 Add integration tests verifying process runs as specified user
- [ ] 6.3 Add test cases for invalid usernames
- [ ] 6.4 Add test cases for permission errors
- [ ] 6.5 Verify tests pass on Unix systems
- [ ] 6.6 Document test requirements (may need specific user accounts)

## 7. Validation
- [ ] 7.1 Run `go test ./...` to verify all tests pass
- [ ] 7.2 Run `golangci-lint run` to ensure linting passes
- [ ] 7.3 Verify example code compiles and runs correctly
- [ ] 7.4 Test on multiple Unix platforms (Linux, macOS) if possible
