# Change: Add Standalone Query() Function

## Why

The Go SDK currently only provides the `ClaudeSDKClient` for interacting with Claude Agent, which is a stateful, bidirectional client designed for multi-turn conversations. This is overly complex for simple, one-shot queries where users just want to send a prompt and receive a response without managing client state, sessions, or follow-up interactions.

The Python SDK provides a standalone `query()` function that offers a simpler, stateless interface for fire-and-forget queries. Adding a similar `Query()` function to the Go SDK would improve API ergonomics and provide better alignment with the Python SDK, making it easier for developers to choose the right tool for their use case.

## What Changes

- Add package-level `Query()` function to `pkg/claude/` for simple, stateless queries
- Function signature: `Query(ctx context.Context, prompt string, opts *Options) (<-chan SDKMessage, error)`
- Handles complete lifecycle: connect, send prompt, stream messages, and close
- Returns a channel that automatically closes when the query completes (on ResultMessage or error)
- Internal implementation creates a `queryImpl`, sends the prompt, and manages the message stream
- Add documentation explaining when to use `Query()` vs `ClaudeSDKClient`

## Impact

- **Affected specs**: `query-interface` (new standalone function capability)
- **Affected code**:
  - New function in `pkg/claude/query.go` or `pkg/claude/standalone.go`
  - Leverages existing `newQueryImpl()` internal constructor
  - No changes to existing `Query` interface or `ClaudeSDKClient`
- **Breaking changes**: None - this is a purely additive change
- **Benefits**:
  - Simpler API for one-shot queries
  - Better Python SDK parity
  - Reduced boilerplate for simple use cases
