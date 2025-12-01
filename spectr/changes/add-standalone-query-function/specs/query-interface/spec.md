# Specification: Query Interface

## ADDED Requirements

### Requirement: Standalone Query Function

The SDK SHALL provide a package-level `Query()` function for simple, stateless, one-shot queries that returns a channel of messages and handles the complete lifecycle automatically.

#### Scenario: Simple one-shot query

- **WHEN** user calls `Query(ctx, "What is 2+2?", nil)`
- **THEN** function creates internal query, sends prompt, and returns channel that streams all messages until completion
- **AND** channel automatically closes when ResultMessage is received or an error occurs
- **AND** underlying query resources are cleaned up automatically

#### Scenario: Query with custom options

- **WHEN** user calls `Query(ctx, prompt, &Options{Model: "claude-3-5-sonnet-20241022"})`
- **THEN** query is created with specified options
- **AND** messages are streamed through returned channel
- **AND** options are applied to the underlying query session

#### Scenario: Context cancellation during query

- **WHEN** user cancels context while Query() is streaming
- **THEN** channel closes gracefully
- **AND** underlying query is cleaned up
- **AND** no goroutine leaks occur

#### Scenario: Error during query initialization

- **WHEN** query creation fails (e.g., invalid options, CLI not found)
- **THEN** Query() returns nil channel and non-nil error
- **AND** no resources are leaked

#### Scenario: Automatic cleanup on completion

- **WHEN** query completes successfully with ResultMessage
- **THEN** channel is closed automatically
- **AND** underlying queryImpl is closed and cleaned up
- **AND** goroutines are terminated properly

### Requirement: Query Function Documentation

The `Query()` function godoc SHALL clearly explain when to use it versus `ClaudeSDKClient` and provide usage examples.

#### Scenario: Documentation includes use case guidance

- **WHEN** developer reads `Query()` godoc
- **THEN** documentation explains "When to use Query()"
- **AND** documentation explains "When to use ClaudeSDKClient"
- **AND** differences between stateless and stateful approaches are clear

#### Scenario: Documentation includes code examples

- **WHEN** developer reads `Query()` godoc
- **THEN** documentation includes basic usage example
- **AND** documentation shows how to iterate over messages
- **AND** documentation demonstrates error handling pattern

### Requirement: Python SDK Parity

The Go `Query()` function SHALL provide equivalent functionality to Python's `query()` function for one-shot, unidirectional queries.

#### Scenario: Functional equivalence with Python SDK

- **WHEN** comparing Go `Query()` to Python `query()`
- **THEN** both support prompt string as input
- **AND** both support optional configuration options
- **AND** both return message stream/iterator
- **AND** both handle complete lifecycle automatically
- **AND** both are stateless (each call is independent)

#### Scenario: API ergonomics match Python patterns

- **WHEN** developer uses Go `Query()`
- **THEN** experience is similar to Python's `async for message in query(...)`
- **AND** Go channel iteration provides equivalent message-by-message processing
- **AND** cleanup is automatic like Python's async context management
