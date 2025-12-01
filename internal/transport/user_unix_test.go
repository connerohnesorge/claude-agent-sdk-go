//go:build !windows

package transport

import (
	"errors"
	"os/exec"
	"os/user"
	"testing"
)

func TestResolveUserCredential_EmptyUsername(t *testing.T) {
	cred, err := resolveUserCredential("")
	if err != nil {
		t.Errorf("expected no error for empty username, got: %v", err)
	}
	if cred != nil {
		t.Error("expected nil credential for empty username")
	}
}

func TestResolveUserCredential_CurrentUser(t *testing.T) {
	// Get current user
	currentUser, err := user.Current()
	if err != nil {
		t.Skipf("could not get current user: %v", err)
	}

	cred, err := resolveUserCredential(currentUser.Username)
	if err != nil {
		t.Errorf("failed to resolve current user: %v", err)
	}
	if cred == nil {
		t.Fatal("expected non-nil credential")
	}

	// Verify UID matches
	if cred.Uid == 0 && currentUser.Uid != "0" {
		t.Error("expected non-root UID for non-root user")
	}
}

func TestResolveUserCredential_InvalidUser(t *testing.T) {
	// Use a username that definitely doesn't exist
	_, err := resolveUserCredential("nonexistent_user_abc123xyz")
	if err == nil {
		t.Error("expected error for nonexistent user")
	}

	// Verify error wraps ErrUserLookupFailed
	if !errors.Is(err, ErrUserLookupFailed) {
		t.Errorf("expected error to wrap ErrUserLookupFailed, got: %v", err)
	}
}

func TestResolveUserCredential_RootUser(t *testing.T) {
	// Try to resolve root user (should exist on all Unix systems)
	cred, err := resolveUserCredential("root")
	if err != nil {
		// On some systems root might not be accessible
		t.Skipf("could not resolve root user: %v", err)
	}
	if cred == nil {
		t.Fatal("expected non-nil credential for root")
	}
	if cred.Uid != 0 {
		t.Errorf("expected UID 0 for root, got: %d", cred.Uid)
	}
	if cred.Gid != 0 {
		t.Errorf("expected GID 0 for root, got: %d", cred.Gid)
	}
}

func TestConfigureUserCredential_EmptyUsername(t *testing.T) {
	cmd := exec.Command("echo", "test")

	err := configureUserCredential(cmd, "")
	if err != nil {
		t.Errorf("expected no error for empty username, got: %v", err)
	}

	// SysProcAttr should remain nil
	if cmd.SysProcAttr != nil && cmd.SysProcAttr.Credential != nil {
		t.Error("expected no credential to be set for empty username")
	}
}

func TestConfigureUserCredential_CurrentUser(t *testing.T) {
	// Get current user
	currentUser, err := user.Current()
	if err != nil {
		t.Skipf("could not get current user: %v", err)
	}

	cmd := exec.Command("echo", "test")

	err = configureUserCredential(cmd, currentUser.Username)
	if err != nil {
		t.Errorf("failed to configure current user credential: %v", err)
	}

	// SysProcAttr should be set
	if cmd.SysProcAttr == nil {
		t.Fatal("expected SysProcAttr to be set")
	}
	if cmd.SysProcAttr.Credential == nil {
		t.Fatal("expected Credential to be set")
	}
}

func TestConfigureUserCredential_InvalidUser(t *testing.T) {
	cmd := exec.Command("echo", "test")

	err := configureUserCredential(cmd, "nonexistent_user_abc123xyz")
	if err == nil {
		t.Error("expected error for nonexistent user")
	}

	// Verify error wraps ErrUserLookupFailed
	if !errors.Is(err, ErrUserLookupFailed) {
		t.Errorf("expected error to wrap ErrUserLookupFailed, got: %v", err)
	}
}
