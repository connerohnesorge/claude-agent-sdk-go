//go:build windows

package transport

import "os/exec"

// configureUserCredential is a no-op on Windows.
// User switching via SysProcAttr.Credential is not supported on Windows.
// Windows requires different APIs (CreateProcessAsUser, LogonUser) which are
// not implemented in this SDK.
//
// Returns nil for any input on Windows.
func configureUserCredential(cmd *exec.Cmd, username string) error {
	// Windows doesn't support syscall.Credential
	// Silently ignore user switching on Windows
	return nil
}
