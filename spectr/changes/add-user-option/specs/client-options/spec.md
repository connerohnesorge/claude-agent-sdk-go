# Client Options Delta

## ADDED Requirements

### Requirement: User Option for Process Isolation
The ClientOptions struct SHALL support a `User` field that allows running the Claude Code CLI subprocess as a different user for security isolation and privilege management.

#### Scenario: Run subprocess as specified user
- **GIVEN** a ClientOptions with `User: "appuser"`
- **WHEN** the client initializes and spawns the Claude Code CLI subprocess
- **THEN** the subprocess SHALL run with the UID and GID of user "appuser"
- **AND** the user credential SHALL be resolved via OS user lookup
- **AND** the process SHALL have the primary group of the specified user

#### Scenario: Run subprocess as current user when User is empty
- **GIVEN** a ClientOptions with `User: ""` or User field omitted
- **WHEN** the client initializes and spawns the subprocess
- **THEN** the subprocess SHALL run with the current process's UID and GID
- **AND** no credential switching SHALL occur

#### Scenario: Invalid username returns error
- **GIVEN** a ClientOptions with `User: "nonexistentuser"`
- **WHEN** the client attempts to initialize
- **THEN** an error SHALL be returned indicating user lookup failed
- **AND** the error message SHALL include the username that failed to resolve
- **AND** the subprocess SHALL NOT be started

#### Scenario: Insufficient permissions returns error
- **GIVEN** a ClientOptions with `User: "root"`
- **AND** the parent process does not have CAP_SETUID capability
- **WHEN** the client attempts to spawn the subprocess
- **THEN** an error SHALL be returned indicating permission denied
- **AND** the error SHALL be propagated to the caller

#### Scenario: Documentation explains security implications
- **WHEN** reviewing the `User` field documentation
- **THEN** it SHALL clearly state this is a Unix-specific feature
- **AND** it SHALL explain that parent process requires appropriate privileges (typically root or CAP_SETUID)
- **AND** it SHALL document use cases: security isolation, privilege reduction, multi-tenant environments
- **AND** it SHALL note that Windows is not supported

#### Scenario: User credential resolution includes primary group
- **GIVEN** a ClientOptions with `User: "appuser"`
- **AND** user "appuser" has primary GID 1001
- **WHEN** the subprocess is spawned
- **THEN** the process SHALL run with UID of "appuser"
- **AND** the process SHALL run with GID 1001
- **AND** the primary group SHALL be set correctly

#### Scenario: Python SDK parity achieved
- **GIVEN** the Python SDK's `ClaudeAgentOptions.user` field
- **WHEN** comparing functionality with Go SDK's `User` field
- **THEN** both SHALL provide equivalent user-switching capabilities
- **AND** both SHALL use OS-level user lookup
- **AND** both SHALL be Unix-specific
