// Package claude provides a high-level SDK for interacting with Claude AI
// agents through the Claude Code CLI.
//
// This package wraps the lower-level API and protocol implementations
// to provide a simple interface for creating and managing Claude agent queries
// and conversations.
//
// # Quick Start
//
// The package provides two main ways to interact with Claude:
//
// 1. [SimpleQuery] - For simple, one-shot queries (recommended for most use cases)
// 2. [ClaudeSDKClient] - For stateful, multi-turn conversations
//
// # Using SimpleQuery
//
// For simple queries where you just need to send a prompt and receive a response:
//
//	ctx := context.Background()
//	msgs, err := claude.SimpleQuery(ctx, "What is 2+2?", nil)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	for msg := range msgs {
//	    switch m := msg.(type) {
//	    case *claude.SDKAssistantMessage:
//	        fmt.Println("Response:", m.Message.Content)
//	    case *claude.SDKResultMessage:
//	        fmt.Printf("Completed in %dms\n", m.DurationMS)
//	    }
//	}
//
// # Using ClaudeSDKClient
//
// For interactive sessions with follow-up messages and control capabilities:
//
//	client, err := claude.NewClient(&claude.Options{
//	    Model: "claude-sonnet-4-5",
//	})
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer client.Close()
//
//	// Start a conversation
//	if err := client.Query(ctx, "Help me understand this code"); err != nil {
//	    log.Fatal(err)
//	}
//
//	// Process responses
//	for msg := range client.ReceiveResponse(ctx) {
//	    // Handle messages...
//	}
//
//	// Send follow-up
//	if err := client.Query(ctx, "Can you explain the error handling?"); err != nil {
//	    log.Fatal(err)
//	}
//
// # Choosing Between SimpleQuery and ClaudeSDKClient
//
//	| Feature                  | SimpleQuery | ClaudeSDKClient |
//	|--------------------------|-------------|-----------------|
//	| One-shot queries         | ✓           | ✓               |
//	| Multi-turn conversations | ✗           | ✓               |
//	| Send follow-up messages  | ✗           | ✓               |
//	| Interrupt processing     | ✗           | ✓               |
//	| Dynamic model switching  | ✗           | ✓               |
//	| Automatic cleanup        | ✓           | Manual Close()  |
//	| Complexity               | Simple      | Full-featured   |
//
// Use [SimpleQuery] when:
//   - You have simple one-off questions
//   - Batch processing independent prompts
//   - Building automated scripts or CI/CD pipelines
//   - You know all inputs upfront
//
// Use [ClaudeSDKClient] when:
//   - Building interactive chat applications
//   - Need to send messages based on responses
//   - Need interrupt or dynamic control capabilities
//   - Managing long-running sessions with state
package claude
