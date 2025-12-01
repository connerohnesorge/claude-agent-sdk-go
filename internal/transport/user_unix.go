//go:build !windows

package transport

import (
	"fmt"
	"os/exec"
	"os/user"
	"strconv"
	"syscall"
)

// resolveUserCredential resolves a username to syscall.Credential for process credentials.
// This function looks up the user by username and returns the UID and GID needed
// to set process credentials via SysProcAttr.Credential.
//
// Returns nil, nil if username is empty (no user switching requested).
// Returns an error if the user cannot be found or the UID/GID cannot be parsed.
func resolveUserCredential(username string) (*syscall.Credential, error) {
	if username == "" {
		return nil, nil
	}

	u, err := user.Lookup(username)
	if err != nil {
		return nil, fmt.Errorf("%w: user '%s': %v", ErrUserLookupFailed, username, err)
	}

	uid, err := strconv.ParseUint(u.Uid, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid UID '%s' for user '%s': %v", ErrInvalidUserID, u.Uid, username, err)
	}

	gid, err := strconv.ParseUint(u.Gid, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid GID '%s' for user '%s': %v", ErrInvalidUserID, u.Gid, username, err)
	}

	return &syscall.Credential{
		Uid: uint32(uid),
		Gid: uint32(gid),
	}, nil
}

// configureUserCredential configures the command to run as the specified user.
// If username is empty, this is a no-op.
// Returns an error if user resolution fails.
func configureUserCredential(cmd *exec.Cmd, username string) error {
	if username == "" {
		return nil
	}

	cred, err := resolveUserCredential(username)
	if err != nil {
		return err
	}

	if cred != nil {
		if cmd.SysProcAttr == nil {
			cmd.SysProcAttr = &syscall.SysProcAttr{}
		}
		cmd.SysProcAttr.Credential = cred
	}

	return nil
}
