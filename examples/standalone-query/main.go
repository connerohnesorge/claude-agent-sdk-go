// Package main demonstrates the SimpleQuery function for one-shot queries.
//
// This example shows how to use the standalone SimpleQuery function for
// simple, stateless interactions with Claude. This is the recommended
// approach for one-off queries where you don't need multi-turn conversations
// or advanced control capabilities.
//
// Key benefits of SimpleQuery over ClaudeSDKClient:
//   - Simpler API: No need to manage client lifecycle
//   - Automatic cleanup: Resources are cleaned up automatically
//   - Fire-and-forget: Just send a prompt and iterate over responses
//
// When to use SimpleQuery:
//   - Simple one-off questions
//   - Batch processing independent prompts
//   - Automated scripts and CI/CD pipelines
//   - Code generation or analysis tasks
//
// When to use ClaudeSDKClient instead:
//   - Interactive conversations with follow-ups
//   - Need to interrupt processing
//   - Need dynamic model switching
//   - Long-running sessions with state
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/connerohnesorge/claude-agent-sdk-go/pkg/claude"
)

func main() {
	ctx := context.Background()

	// Simple query with default options
	fmt.Println("=== Simple Query Example ===")
	fmt.Println("Sending: What is 2+2? Just respond with the number.")
	fmt.Println()

	// Use SimpleQuery for fire-and-forget queries
	// The channel automatically closes when the query completes
	msgs, err := claude.SimpleQuery(ctx, "What is 2+2? Just respond with the number.", nil)
	if err != nil {
		log.Fatalf("Failed to start query: %v", err)
	}

	// Iterate over all messages until the channel closes
	for msg := range msgs {
		handleMessage(msg)
	}

	fmt.Println("\n=== Query with Options Example ===")

	// Query with custom options
	opts := &claude.Options{
		Model: "claude-sonnet-4-5",
		// PermissionMode: claude.PermissionModeBypassPermissions,
	}

	msgs, err = claude.SimpleQuery(ctx, "What are the first 5 prime numbers? List them.", opts)
	if err != nil {
		log.Fatalf("Failed to start query with options: %v", err)
	}

	for msg := range msgs {
		handleMessage(msg)
	}

	fmt.Println("\nDone!")
}

// handleMessage processes different types of SDK messages.
func handleMessage(msg claude.SDKMessage) {
	switch m := msg.(type) {
	case *claude.SDKSystemMessage:
		handleSystemMessage(m)
	case *claude.SDKAssistantMessage:
		handleAssistantMessage(m)
	case *claude.SDKResultMessage:
		handleResultMessage(m)
	case *claude.SDKUserMessage:
		// User messages are echoed back, typically the initial prompt
		fmt.Println("[User message received]")
	}
}

// handleSystemMessage processes system initialization messages.
func handleSystemMessage(m *claude.SDKSystemMessage) {
	if m.Subtype == "init" {
		fmt.Printf("[System] Initialized with model: %v\n", m.Data["model"])
	}
}

// handleAssistantMessage displays assistant response content.
func handleAssistantMessage(m *claude.SDKAssistantMessage) {
	fmt.Println("\n[Assistant]")
	for _, block := range m.Message.Content {
		displayContentBlock(block)
	}
}

// displayContentBlock prints text content from a content block.
func displayContentBlock(block claude.ContentBlock) {
	switch b := block.(type) {
	case claude.TextBlock:
		fmt.Printf("  %s\n", b.Text)
	case claude.TextContentBlock:
		fmt.Printf("  %s\n", b.Text)
	}
}

// handleResultMessage displays final result statistics.
func handleResultMessage(m *claude.SDKResultMessage) {
	fmt.Printf("\n[Result] Status: %s\n", m.Subtype)
	fmt.Printf("  Duration: %dms\n", m.DurationMS)
	fmt.Printf("  Cost: $%.4f\n", m.TotalCostUSD)
	fmt.Printf("  Turns: %d\n", m.NumTurns)
}
