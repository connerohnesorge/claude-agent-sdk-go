package clauderrs

import (
	"errors"
	"fmt"
)

// ErrBufferSizeExceeded is a sentinel error for buffer size exceeded errors.
var ErrBufferSizeExceeded = errors.New("buffer size exceeded")

// BufferError represents buffer-related errors.
type BufferError struct {
	*BaseError
	limit      int
	actualSize int
	operation  string
}

// NewBufferSizeExceededError creates a new buffer size exceeded error.
// It takes the size limit in bytes, the actual size in bytes, and the operation
// that caused the buffer size to be exceeded.
func NewBufferSizeExceededError(limit, actualSize int, operation string) *BufferError {
	message := fmt.Sprintf("buffer size exceeded: limit=%d bytes, actual=%d bytes, operation=%s", limit, actualSize, operation)
	err := &BufferError{
		BaseError:  NewBaseError(CategoryTransport, ErrCodeBufferSizeExceeded, message, nil),
		limit:      limit,
		actualSize: actualSize,
		operation:  operation,
	}

	// Add buffer-specific metadata
	_ = err.WithMetadata("limit", limit)
	_ = err.WithMetadata("actual_size", actualSize)
	_ = err.WithMetadata("operation", operation)

	return err
}

// Limit returns the buffer size limit in bytes.
func (e *BufferError) Limit() int {
	return e.limit
}

// ActualSize returns the actual buffer size in bytes.
func (e *BufferError) ActualSize() int {
	return e.actualSize
}

// Operation returns the operation that caused the buffer size to be exceeded.
func (e *BufferError) Operation() string {
	return e.operation
}

// Is implements error matching for errors.Is().
func (e *BufferError) Is(target error) bool {
	return target == ErrBufferSizeExceeded
}
