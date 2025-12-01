package clauderrs

import (
	"errors"
	"testing"
)

// TestBufferError_Creation tests that NewBufferSizeExceededError creates error with correct fields
func TestBufferError_Creation(t *testing.T) {
	limit := 1024
	actualSize := 2048
	operation := "test_operation"

	err := NewBufferSizeExceededError(limit, actualSize, operation)

	if err == nil {
		t.Fatal("NewBufferSizeExceededError should not return nil")
	}

	// Verify fields via accessor methods
	if err.Limit() != limit {
		t.Errorf("Limit() = %d, want %d", err.Limit(), limit)
	}
	if err.ActualSize() != actualSize {
		t.Errorf("ActualSize() = %d, want %d", err.ActualSize(), actualSize)
	}
	if err.Operation() != operation {
		t.Errorf("Operation() = %q, want %q", err.Operation(), operation)
	}
}

// TestBufferError_ErrorMessage tests that error message format is correct
func TestBufferError_ErrorMessage(t *testing.T) {
	limit := 1024
	actualSize := 2048
	operation := "read_json"

	err := NewBufferSizeExceededError(limit, actualSize, operation)

	expectedMsg := "transport: buffer size exceeded: limit=1024 bytes, actual=2048 bytes, operation=read_json"
	if err.Error() != expectedMsg {
		t.Errorf("Error() = %q, want %q", err.Error(), expectedMsg)
	}
}

// TestBufferError_ErrorsIs tests that errors.Is() works with ErrBufferSizeExceeded sentinel
func TestBufferError_ErrorsIs(t *testing.T) {
	err := NewBufferSizeExceededError(1024, 2048, "test_op")

	if !errors.Is(err, ErrBufferSizeExceeded) {
		t.Error("errors.Is() should return true for ErrBufferSizeExceeded")
	}
}

// TestBufferError_Limit tests that Limit() method returns correct value
func TestBufferError_Limit(t *testing.T) {
	tests := []struct {
		name  string
		limit int
	}{
		{"small limit", 512},
		{"1KB limit", 1024},
		{"1MB limit", 1024 * 1024},
		{"10MB limit", 10 * 1024 * 1024},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewBufferSizeExceededError(tt.limit, tt.limit*2, "test")
			if err.Limit() != tt.limit {
				t.Errorf("Limit() = %d, want %d", err.Limit(), tt.limit)
			}
		})
	}
}

// TestBufferError_ActualSize tests that ActualSize() method returns correct value
func TestBufferError_ActualSize(t *testing.T) {
	tests := []struct {
		name       string
		actualSize int
	}{
		{"small size", 1536},
		{"2KB size", 2048},
		{"2MB size", 2 * 1024 * 1024},
		{"20MB size", 20 * 1024 * 1024},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewBufferSizeExceededError(1024, tt.actualSize, "test")
			if err.ActualSize() != tt.actualSize {
				t.Errorf("ActualSize() = %d, want %d", err.ActualSize(), tt.actualSize)
			}
		})
	}
}

// TestBufferError_Operation tests that Operation() method returns correct value
func TestBufferError_Operation(t *testing.T) {
	tests := []struct {
		name      string
		operation string
	}{
		{"read operation", "read_json"},
		{"write operation", "write_buffer"},
		{"parse operation", "parse_response"},
		{"accumulate operation", "accumulate_stdout"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewBufferSizeExceededError(1024, 2048, tt.operation)
			if err.Operation() != tt.operation {
				t.Errorf("Operation() = %q, want %q", err.Operation(), tt.operation)
			}
		})
	}
}

// TestBufferError_AllFields tests all accessor methods together
func TestBufferError_AllFields(t *testing.T) {
	limit := 5120
	actualSize := 10240
	operation := "complex_operation"

	err := NewBufferSizeExceededError(limit, actualSize, operation)

	if err.Limit() != limit {
		t.Errorf("Limit() = %d, want %d", err.Limit(), limit)
	}
	if err.ActualSize() != actualSize {
		t.Errorf("ActualSize() = %d, want %d", err.ActualSize(), actualSize)
	}
	if err.Operation() != operation {
		t.Errorf("Operation() = %q, want %q", err.Operation(), operation)
	}
}

// TestBufferError_Metadata tests that metadata is set correctly
func TestBufferError_Metadata(t *testing.T) {
	limit := 1024
	actualSize := 2048
	operation := "test_op"

	err := NewBufferSizeExceededError(limit, actualSize, operation)

	// Verify metadata is accessible through the error
	meta := err.Metadata()
	if meta == nil {
		t.Fatal("Metadata() should not return nil")
	}

	if limitMeta, ok := meta["limit"]; !ok {
		t.Error("metadata should contain 'limit' key")
	} else if limitMeta != limit {
		t.Errorf("metadata['limit'] = %v, want %d", limitMeta, limit)
	}

	if actualMeta, ok := meta["actual_size"]; !ok {
		t.Error("metadata should contain 'actual_size' key")
	} else if actualMeta != actualSize {
		t.Errorf("metadata['actual_size'] = %v, want %d", actualMeta, actualSize)
	}

	if opMeta, ok := meta["operation"]; !ok {
		t.Error("metadata should contain 'operation' key")
	} else if opMeta != operation {
		t.Errorf("metadata['operation'] = %v, want %q", opMeta, operation)
	}
}
