package transport

import (
	"context"
	"io"
	"strings"
	"testing"
	"errors"

	"github.com/connerohnesorge/claude-agent-sdk-go/pkg/clauderrs"
)

// mockReadCloser implements io.ReadCloser for testing
type mockReadCloser struct {
	*strings.Reader
}

func (m *mockReadCloser) Close() error {
	return nil
}

// mockWriteCloser implements io.WriteCloser for testing
type mockWriteCloser struct {
	io.Writer
}

func (m *mockWriteCloser) Close() error {
	return nil
}

func TestStdioTransport_MaxBufferSize_Unlimited(t *testing.T) {
	// Test that maxBufferSize=0 means unlimited
	input := strings.Repeat("a", 10000) + "\n"
	stdout := &mockReadCloser{strings.NewReader(input)}
	stdin := &mockWriteCloser{io.Discard}
	stderr := &mockReadCloser{strings.NewReader("")}

	transport := NewStdioTransport(stdin, stdout, stderr, 0)

	ctx := context.Background()
	data, err := transport.Read(ctx)
	if err != nil {
		t.Fatalf("Expected no error with maxBufferSize=0, got: %v", err)
	}

	if len(data) != len(input) {
		t.Errorf("Expected to read %d bytes, got %d", len(input), len(data))
	}
}

func TestStdioTransport_MaxBufferSize_WithinLimit(t *testing.T) {
	// Test that reading within the limit works
	input := "hello world\n"
	maxBufferSize := 100
	stdout := &mockReadCloser{strings.NewReader(input)}
	stdin := &mockWriteCloser{io.Discard}
	stderr := &mockReadCloser{strings.NewReader("")}

	transport := NewStdioTransport(stdin, stdout, stderr, maxBufferSize)

	ctx := context.Background()
	data, err := transport.Read(ctx)
	if err != nil {
		t.Fatalf("Expected no error when within limit, got: %v", err)
	}

	if string(data) != input {
		t.Errorf("Expected to read %q, got %q", input, string(data))
	}
}

func TestStdioTransport_MaxBufferSize_ExceedsLimit(t *testing.T) {
	// Test that exceeding the limit returns an error
	maxBufferSize := 10
	input := strings.Repeat("a", maxBufferSize+5) + "\n" // Exceeds limit
	stdout := &mockReadCloser{strings.NewReader(input)}
	stdin := &mockWriteCloser{io.Discard}
	stderr := &mockReadCloser{strings.NewReader("")}

	transport := NewStdioTransport(stdin, stdout, stderr, maxBufferSize)

	ctx := context.Background()
	_, err := transport.Read(ctx)
	if err == nil {
		t.Fatal("Expected error when exceeding buffer size limit, got nil")
	}

	// Check that it's a BufferSizeExceeded error
	if !errors.Is(err, clauderrs.ErrBufferSizeExceeded) {
		t.Errorf("Expected ErrBufferSizeExceeded, got: %v", err)
	}

	// Check the error details
	var bufferErr *clauderrs.BufferError
	if errors.As(err, &bufferErr) {
		if bufferErr.Limit() != maxBufferSize {
			t.Errorf("Expected limit %d, got %d", maxBufferSize, bufferErr.Limit())
		}
		if bufferErr.ActualSize() <= maxBufferSize {
			t.Errorf("Expected actual size > %d, got %d", maxBufferSize, bufferErr.ActualSize())
		}
		if bufferErr.Operation() != "json_accumulation" {
			t.Errorf("Expected operation 'json_accumulation', got %q", bufferErr.Operation())
		}
	} else {
		t.Error("Expected error to be *clauderrs.BufferError")
	}
}

func TestStdioTransport_MaxBufferSize_ExactLimit(t *testing.T) {
	// Test behavior at exact limit including newline (should succeed)
	maxBufferSize := 11 // 10 bytes of data + 1 newline
	input := strings.Repeat("a", 10) + "\n"
	stdout := &mockReadCloser{strings.NewReader(input)}
	stdin := &mockWriteCloser{io.Discard}
	stderr := &mockReadCloser{strings.NewReader("")}

	transport := NewStdioTransport(stdin, stdout, stderr, maxBufferSize)

	ctx := context.Background()
	data, err := transport.Read(ctx)
	if err != nil {
		t.Fatalf("Expected no error at exact limit, got: %v", err)
	}

	if len(data) != len(input) {
		t.Errorf("Expected to read %d bytes, got %d", len(input), len(data))
	}
}

func TestStdioTransport_MaxBufferSize_OneByteOver(t *testing.T) {
	// Test that one byte over the limit triggers error
	maxBufferSize := 10
	input := strings.Repeat("a", maxBufferSize+1) + "\n"
	stdout := &mockReadCloser{strings.NewReader(input)}
	stdin := &mockWriteCloser{io.Discard}
	stderr := &mockReadCloser{strings.NewReader("")}

	transport := NewStdioTransport(stdin, stdout, stderr, maxBufferSize)

	ctx := context.Background()
	_, err := transport.Read(ctx)
	if err == nil {
		t.Fatal("Expected error when one byte over limit, got nil")
	}

	if !errors.Is(err, clauderrs.ErrBufferSizeExceeded) {
		t.Errorf("Expected ErrBufferSizeExceeded, got: %v", err)
	}
}
