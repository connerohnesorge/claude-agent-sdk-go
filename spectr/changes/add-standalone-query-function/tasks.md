# Implementation Tasks

## 1. Core Implementation
- [ ] 1.1 Add `Query(ctx context.Context, prompt string, opts *Options) (<-chan SDKMessage, error)` function
- [ ] 1.2 Implement message streaming logic that auto-closes channel on completion
- [ ] 1.3 Ensure proper cleanup/close of underlying queryImpl on error or completion
- [ ] 1.4 Handle context cancellation correctly

## 2. Documentation
- [ ] 2.1 Add comprehensive godoc for `Query()` function
- [ ] 2.2 Add usage examples in godoc
- [ ] 2.3 Document when to use `Query()` vs `ClaudeSDKClient`
- [ ] 2.4 Add comparison table to package documentation

## 3. Testing
- [ ] 3.1 Write unit tests for happy path (successful query)
- [ ] 3.2 Write unit tests for error handling
- [ ] 3.3 Write unit tests for context cancellation
- [ ] 3.4 Write integration test with actual Claude Code CLI
- [ ] 3.5 Create example in `examples/standalone-query/`
- [ ] 3.6 Verify example runs successfully

## 4. Validation
- [ ] 4.1 Run `go test ./...` and verify all tests pass
- [ ] 4.2 Run linter (`golangci-lint run`) and fix any issues
- [ ] 4.3 Verify examples compile and run
- [ ] 4.4 Review against Python SDK `query()` for API consistency
