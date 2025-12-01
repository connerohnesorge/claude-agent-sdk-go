# Implementation Tasks

## 1. Core Implementation
- [x] 1.1 Add `Query(ctx context.Context, prompt string, opts *Options) (<-chan SDKMessage, error)` function
- [x] 1.2 Implement message streaming logic that auto-closes channel on completion
- [x] 1.3 Ensure proper cleanup/close of underlying queryImpl on error or completion
- [x] 1.4 Handle context cancellation correctly

## 2. Documentation
- [x] 2.1 Add comprehensive godoc for `Query()` function
- [x] 2.2 Add usage examples in godoc
- [x] 2.3 Document when to use `Query()` vs `ClaudeSDKClient`
- [x] 2.4 Add comparison table to package documentation

## 3. Testing
- [x] 3.1 Write unit tests for happy path (successful query)
- [x] 3.2 Write unit tests for error handling
- [x] 3.3 Write unit tests for context cancellation
- [x] 3.4 Write integration test with actual Claude Code CLI
- [x] 3.5 Create example in `examples/standalone-query/`
- [x] 3.6 Verify example runs successfully

## 4. Validation
- [x] 4.1 Run `go test ./...` and verify all tests pass
- [x] 4.2 Run linter (`golangci-lint run`) and fix any issues
- [x] 4.3 Verify examples compile and run
- [x] 4.4 Review against Python SDK `query()` for API consistency
