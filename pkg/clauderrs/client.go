package clauderrs

// ClientError represents client-related errors.
type ClientError struct {
	*BaseError
}

// NewClientError creates a new client error.
func NewClientError(code ErrorCode, message string, cause error) *ClientError {
	return &ClientError{
		BaseError: NewBaseError(CategoryClient, code, message, cause),
	}
}

// WithSessionID adds session ID metadata to the error.
func (e *ClientError) WithSessionID(sessionID string) *ClientError {
	_ = e.WithMetadata(MetadataKeySessionID, sessionID)

	return e
}

// NewVersionMismatchError creates a client error for version mismatches.
func NewVersionMismatchError(currentVersion string, minimumVersion string) *ClientError {
	message := "Claude Code CLI version " + currentVersion + " is below minimum required version " + minimumVersion + ". Please upgrade the CLI to continue."
	err := NewClientError(ErrCodeVersionMismatch, message, nil)
	_ = err.WithMetadata(MetadataKeyCurrentVersion, currentVersion)
	_ = err.WithMetadata(MetadataKeyMinimumVersion, minimumVersion)

	return err
}
