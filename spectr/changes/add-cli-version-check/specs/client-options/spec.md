## ADDED Requirements

### Requirement: CLI Version Validation
The SDK SHALL validate the Claude Code CLI version at connection time to ensure compatibility with the minimum supported version.

#### Scenario: CLI version meets minimum requirement
- **GIVEN** Claude Code CLI version 2.0.0 or higher is installed
- **WHEN** a new client connection is established
- **THEN** the version check passes
- **AND** the connection proceeds normally

#### Scenario: CLI version below minimum requirement
- **GIVEN** Claude Code CLI version 1.9.0 is installed
- **WHEN** a new client connection is attempted
- **THEN** a version mismatch error is returned
- **AND** the error message includes the current version and minimum required version
- **AND** the connection is not established

#### Scenario: CLI version cannot be determined
- **GIVEN** the `claude --version` command fails or returns invalid output
- **WHEN** a new client connection is attempted
- **THEN** a version check error is returned
- **AND** the error message indicates the version could not be determined

#### Scenario: Version check skipped via environment variable
- **GIVEN** the environment variable `CLAUDE_AGENT_SDK_SKIP_VERSION_CHECK` is set to `true`
- **WHEN** a new client connection is established
- **THEN** the version check is skipped
- **AND** the connection proceeds without version validation

---

### Requirement: Minimum CLI Version Constant
The SDK SHALL define a constant `MinimumClaudeCodeVersion` representing the minimum supported Claude Code CLI version.

#### Scenario: Minimum version defined
- **WHEN** inspecting the SDK constants
- **THEN** `MinimumClaudeCodeVersion` is set to "2.0.0"
- **AND** the constant is documented with its purpose

#### Scenario: Version comparison uses semantic versioning
- **GIVEN** the minimum version is "2.0.0"
- **WHEN** comparing CLI version "2.1.5"
- **THEN** the version check passes (2.1.5 >= 2.0.0)

#### Scenario: Pre-release versions are handled
- **GIVEN** the minimum version is "2.0.0"
- **WHEN** comparing CLI version "2.0.0-beta.1"
- **THEN** the version check behavior is well-defined in documentation

---

### Requirement: Version Mismatch Error Type
The SDK SHALL provide a dedicated error code `ErrCodeVersionMismatch` for CLI version compatibility failures.

#### Scenario: Version mismatch error contains details
- **GIVEN** CLI version 1.5.0 when minimum is 2.0.0
- **WHEN** version check fails
- **THEN** the error code is `ErrCodeVersionMismatch`
- **AND** the error message includes current version "1.5.0"
- **AND** the error message includes required minimum version "2.0.0"
- **AND** the error suggests upgrading the CLI

#### Scenario: Version check failure is distinguishable
- **WHEN** a version mismatch error occurs
- **THEN** it can be identified via error code inspection
- **AND** it is categorized as a client error
- **AND** it provides actionable guidance for resolution

---

### Requirement: Skip Version Check Environment Variable
The SDK SHALL respect the `CLAUDE_AGENT_SDK_SKIP_VERSION_CHECK` environment variable to bypass version validation.

#### Scenario: Environment variable set to true
- **GIVEN** `CLAUDE_AGENT_SDK_SKIP_VERSION_CHECK=true` is set
- **WHEN** initializing a client
- **THEN** version checking is completely skipped
- **AND** no `claude --version` command is executed
- **AND** connection proceeds regardless of CLI version

#### Scenario: Environment variable set to false
- **GIVEN** `CLAUDE_AGENT_SDK_SKIP_VERSION_CHECK=false` is set
- **WHEN** initializing a client
- **THEN** version checking is performed normally

#### Scenario: Environment variable not set
- **GIVEN** `CLAUDE_AGENT_SDK_SKIP_VERSION_CHECK` is not set
- **WHEN** initializing a client
- **THEN** version checking is performed normally (default behavior)

#### Scenario: Environment variable case insensitive
- **GIVEN** `CLAUDE_AGENT_SDK_SKIP_VERSION_CHECK=True` or `TRUE` is set
- **WHEN** initializing a client
- **THEN** version checking is skipped (case-insensitive match)

---
