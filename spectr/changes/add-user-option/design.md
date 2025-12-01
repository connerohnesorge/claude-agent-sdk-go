# Design: Add User Option for Process Isolation

## Context
The Python SDK provides a `user` field in `ClaudeAgentOptions` that allows running the Claude Code CLI subprocess as a different user via the `anyio.open_process()` function's `user` parameter. This is a critical security feature for:
- **Security isolation**: Running untrusted code with reduced privileges
- **Multi-tenant environments**: Isolating different users' processes
- **Container deployments**: Running as non-root users following security best practices
- **Principle of least privilege**: Ensuring processes only have necessary permissions

The Go SDK currently lacks this capability, creating a gap in feature parity with the Python SDK.

## Goals
- Add `User` field to `Options` struct for specifying the username
- Implement user credential resolution and process attribute configuration
- Achieve parity with Python SDK's `user` option
- Maintain Unix-specific behavior with clear documentation
- Handle errors gracefully (invalid users, permission issues)

## Non-Goals
- Windows support (syscall.Credential is Unix-specific)
- Group management beyond primary group (can be added later if needed)
- User creation or validation beyond OS-level lookup
- Support for numeric UIDs (username strings only for simplicity)

## Decisions

### Decision 1: Use syscall.SysProcAttr.Credential
**Choice**: Use Go's `syscall.SysProcAttr.Credential` field to set process UID/GID.

**Rationale**:
- This is the standard Go approach for setting process credentials
- Provides access to `setuid`/`setgid` functionality at the OS level
- Integrates cleanly with `exec.Cmd` structure
- Type-safe and idiomatic Go

**Alternatives considered**:
1. **Shell wrapper with `sudo` or `su`**: Rejected because it adds external dependencies, complicates error handling, and bypasses Go's type system
2. **CGO with setuid**: Rejected due to increased complexity and CGO requirements
3. **Numeric UID in options**: Rejected for usability; usernames are more user-friendly

### Decision 2: User Credential Resolution Helper
**Choice**: Create a `resolveUserCredential(username string) (*syscall.Credential, error)` helper function.

**Rationale**:
- Encapsulates user lookup logic using `user.Lookup(username)`
- Converts string UID/GID to `uint32` for `syscall.Credential`
- Provides single point of error handling for user resolution
- Makes testing easier (can mock in tests)

**Implementation**:
```go
func resolveUserCredential(username string) (*syscall.Credential, error) {
    u, err := user.Lookup(username)
    if err != nil {
        return nil, fmt.Errorf("user lookup failed: %w", err)
    }

    uid, err := strconv.ParseUint(u.Uid, 10, 32)
    if err != nil {
        return nil, fmt.Errorf("invalid UID: %w", err)
    }

    gid, err := strconv.ParseUint(u.Gid, 10, 32)
    if err != nil {
        return nil, fmt.Errorf("invalid GID: %w", err)
    }

    return &syscall.Credential{
        Uid: uint32(uid),
        Gid: uint32(gid),
    }, nil
}
```

### Decision 3: ProcessConfig Extension
**Choice**: Add `User string` field to `ProcessConfig` struct.

**Rationale**:
- Maintains clean separation between high-level `Options` and low-level `ProcessConfig`
- Transport layer handles platform-specific details
- Keeps `options.go` focused on SDK-level configuration

### Decision 4: Unix-Only Support
**Choice**: Document as Unix-only feature; no Windows implementation initially.

**Rationale**:
- Python SDK's `anyio.open_process(user=...)` is also Unix-specific
- Windows has fundamentally different security model (no direct setuid equivalent)
- Implementing Windows would require `CreateProcessAsUser` or `LogonUser` APIs, significantly increasing complexity
- Can be added later if there's demand

### Decision 5: Error Handling Strategy
**Choice**: Return errors for invalid users rather than silently failing or falling back.

**Rationale**:
- Fail-fast approach prevents security misconfigurations
- User expects explicit error if they specify a user that doesn't exist
- Consistent with Go error handling conventions
- Aligns with Python SDK behavior

## Risks / Trade-offs

### Risk: Permission Requirements
**Risk**: Setting credentials requires the parent process to have `CAP_SETUID` and `CAP_SETGID` capabilities (typically requires root).

**Mitigation**:
- Document in godoc that this requires appropriate permissions
- Provide clear error messages when permission is denied
- Include example showing typical usage patterns

### Risk: Platform-Specific Behavior
**Risk**: Code will only work on Unix-like systems (Linux, macOS, BSD).

**Mitigation**:
- Clear documentation stating Unix-only support
- Could add build tags (`// +build !windows`) if needed
- Windows users will see graceful error if attempted

### Risk: Security Group Handling
**Risk**: Only setting primary group (GID), not supplementary groups.

**Mitigation**:
- Document current limitation in godoc
- `syscall.Credential.Groups` field exists for future enhancement
- Primary group is sufficient for most isolation use cases

### Trade-off: Username vs UID
**Trade-off**: Accepting username strings is more user-friendly but requires OS user lookup.

**Decision**: Prioritize usability. Username strings are easier to use and read in configuration.

## Migration Plan
This is an additive change with no breaking changes:
1. Add `User` field to `Options` struct (optional, defaults to empty string)
2. Update transport layer to handle user switching when specified
3. Existing code continues to work unchanged (empty `User` = current behavior)
4. Users can opt-in by setting `User` field in their options

### Rollback
If issues arise, the feature can be temporarily disabled by:
1. Ignoring the `User` field in transport layer
2. Returning an error if `User` is set
3. Removing the field in a future version (breaking change)

## Open Questions
1. **Should we support numeric UIDs?** - Current design uses usernames only. Could add UID support if needed.
2. **Should we support supplementary groups?** - Current design only sets primary GID. Could extend `Credential.Groups` if required.
3. **Should we validate the user has necessary permissions before attempting switch?** - Would require additional privilege checking logic.

## References
- Python SDK: `.claude/contexts/claude-agent-sdk-python-v0.1.0.md` line 3525 (`user=self._options.user`)
- Python SDK options: line 4357 (`user: str | None = None`)
- Go syscall package: `syscall.SysProcAttr.Credential`
- Go user package: `user.Lookup(username)` for user resolution
