//go:build windows

package transport

// resolveUserCredential is a no-op on Windows.
// User switching via SysProcAttr.Credential is not supported on Windows.
// Windows requires different APIs (CreateProcessAsUser, LogonUser) which are
// not implemented in this SDK.
//
// Returns nil, nil for any input on Windows.
func resolveUserCredential(username string) (any, error) {
	// Windows doesn't support syscall.Credential
	// Return nil to indicate no credential switching
	return nil, nil
}
