package claude

import (
	"context"
	"io"
	"strings"
	"testing"

	"github.com/connerohnesorge/claude-agent-sdk-go/pkg/clauderrs"
)

type stubSimpleQuery struct {
	sessionID string
	responses []stubSimpleQueryResponse
	closed    bool
}

type stubSimpleQueryResponse struct {
	msg SDKMessage
	err error
}

func (s *stubSimpleQuery) Next(context.Context) (SDKMessage, error) {
	if len(s.responses) == 0 {
		return nil, io.EOF
	}

	resp := s.responses[0]
	s.responses = s.responses[1:]

	return resp.msg, resp.err
}

func (s *stubSimpleQuery) Close() error {
	s.closed = true

	return nil
}

func (s *stubSimpleQuery) SessionID() string {
	return s.sessionID
}

func TestStreamSimpleQuery_EmitsStreamErrors(t *testing.T) {
	ctx := context.Background()
	streamErr := clauderrs.NewBufferSizeExceededError(1024, 2048, "json_accumulation")

	stub := &stubSimpleQuery{
		sessionID: "session-123",
		responses: []stubSimpleQueryResponse{
			{err: streamErr},
		},
	}

	out := streamSimpleQuery(ctx, stub)

	msg, ok := <-out
	if !ok {
		t.Fatal("expected error message before channel close")
	}

	result, ok := msg.(*SDKResultMessage)
	if !ok {
		t.Fatalf("expected *SDKResultMessage, got %T", msg)
	}

	if !result.IsError {
		t.Fatal("expected result message to be marked as error")
	}
	if result.Subtype != ResultSubtypeErrorDuringExecution {
		t.Fatalf("unexpected error subtype: %s", result.Subtype)
	}
	if result.SessionID() != stub.sessionID {
		t.Fatalf("expected session ID %s, got %s", stub.sessionID, result.SessionID())
	}
	if len(result.Errors) == 0 || !strings.Contains(result.Errors[0], "buffer size exceeded") {
		t.Fatalf("expected buffer overflow details in error message, got %v", result.Errors)
	}

	if _, ok := <-out; ok {
		t.Fatal("expected channel to be closed after streaming error")
	}
	if !stub.closed {
		t.Fatal("expected query Close() to be invoked")
	}
}
