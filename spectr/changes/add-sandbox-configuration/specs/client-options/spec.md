# Client Options Specification - Sandbox Configuration Delta

## ADDED Requirements

### Requirement: Sandbox Configuration Support
The Options struct SHALL support a `Sandbox` field of type `*SandboxSettings` that enables configuration of bash command sandboxing for security control on macOS and Linux systems.

#### Scenario: Enable sandbox with auto-allow
- **GIVEN** an Options instance with:
  ```go
  Sandbox: &SandboxSettings{
    Enabled: true,
    AutoAllowBashIfSandboxed: true,
  }
  ```
- **WHEN** the client is initialized
- **THEN** bash sandboxing is enabled
- **AND** sandboxed bash commands are automatically approved without prompts

#### Scenario: Sandbox with excluded commands
- **GIVEN** an Options instance with:
  ```go
  Sandbox: &SandboxSettings{
    Enabled: true,
    ExcludedCommands: []string{"docker", "git"},
  }
  ```
- **WHEN** bash commands are executed
- **THEN** "docker" and "git" commands run outside the sandbox
- **AND** other commands run within the sandbox

#### Scenario: Disable sandbox
- **GIVEN** an Options instance with `Sandbox: nil` or `Sandbox.Enabled: false`
- **WHEN** the client executes bash commands
- **THEN** no sandboxing is applied
- **AND** commands run in normal execution mode

#### Scenario: Network configuration with Unix sockets
- **GIVEN** an Options instance with:
  ```go
  Sandbox: &SandboxSettings{
    Enabled: true,
    Network: &SandboxNetworkConfig{
      AllowUnixSockets: []string{"/var/run/docker.sock", "/tmp/custom.sock"},
    },
  }
  ```
- **WHEN** sandboxed commands attempt Unix socket access
- **THEN** access to "/var/run/docker.sock" and "/tmp/custom.sock" is allowed
- **AND** access to other Unix sockets is restricted

#### Scenario: HTTP proxy configuration
- **GIVEN** an Options instance with:
  ```go
  Sandbox: &SandboxSettings{
    Enabled: true,
    Network: &SandboxNetworkConfig{
      HttpProxyPort: 8080,
    },
  }
  ```
- **WHEN** sandboxed commands make HTTP requests
- **THEN** requests are routed through the HTTP proxy on port 8080
- **AND** proxy configuration is applied to outbound traffic

#### Scenario: Ignore file violations
- **GIVEN** an Options instance with:
  ```go
  Sandbox: &SandboxSettings{
    Enabled: true,
    IgnoreViolations: &SandboxIgnoreViolations{
      File: []string{"/tmp/safe-file.txt", "/var/log/app.log"},
    },
  }
  ```
- **WHEN** sandboxed commands access specified file paths
- **THEN** violations for "/tmp/safe-file.txt" and "/var/log/app.log" are ignored
- **AND** access is permitted without security alerts

#### Scenario: Nested sandbox for Docker
- **GIVEN** an Options instance with:
  ```go
  Sandbox: &SandboxSettings{
    Enabled: true,
    EnableWeakerNestedSandbox: true,
  }
  ```
- **WHEN** running within a Docker container or nested environment
- **THEN** a weaker sandbox configuration is used
- **AND** compatibility with container restrictions is maintained

---

### Requirement: SandboxSettings Type Definition
The system SHALL define a `SandboxSettings` struct type that encapsulates all sandbox configuration options with proper JSON serialization.

#### Scenario: SandboxSettings with all fields
- **GIVEN** a SandboxSettings struct with all fields populated
- **WHEN** marshaled to JSON
- **THEN** it produces a valid JSON object with camelCase field names
- **AND** the JSON includes: `enabled`, `autoAllowBashIfSandboxed`, `excludedCommands`, `allowUnsandboxedCommands`, `network`, `ignoreViolations`, `enableWeakerNestedSandbox`

#### Scenario: Minimal SandboxSettings
- **GIVEN** a SandboxSettings with only `Enabled: true`
- **WHEN** marshaled to JSON
- **THEN** only the `enabled: true` field is present
- **AND** optional fields are omitted or have appropriate zero values

#### Scenario: Nested configuration serialization
- **GIVEN** a SandboxSettings with nested `Network` and `IgnoreViolations`
- **WHEN** marshaled to JSON
- **THEN** nested objects are properly serialized
- **AND** the structure matches the Python SDK's JSON format

---

### Requirement: SandboxNetworkConfig Type Definition
The system SHALL define a `SandboxNetworkConfig` struct type for configuring network access within the sandbox.

#### Scenario: Unix socket allowlist
- **GIVEN** a SandboxNetworkConfig with `AllowUnixSockets: []string{"/var/run/docker.sock"}`
- **WHEN** marshaled to JSON
- **THEN** it produces `{"allowUnixSockets": ["/var/run/docker.sock"]}`

#### Scenario: Allow all Unix sockets
- **GIVEN** a SandboxNetworkConfig with `AllowAllUnixSockets: true`
- **WHEN** applied to sandbox configuration
- **THEN** all Unix socket access is permitted
- **AND** the `AllowUnixSockets` list is ignored

#### Scenario: Local port binding on macOS
- **GIVEN** a SandboxNetworkConfig with `AllowLocalBinding: true`
- **WHEN** running on macOS
- **THEN** sandboxed commands can bind to localhost ports
- **AND** the binding permission is platform-appropriate

#### Scenario: Proxy configuration
- **GIVEN** a SandboxNetworkConfig with:
  ```go
  HttpProxyPort: 8080
  SocksProxyPort: 1080
  ```
- **WHEN** marshaled to JSON
- **THEN** it produces `{"httpProxyPort": 8080, "socksProxyPort": 1080}`
- **AND** both proxy types can be configured simultaneously

---

### Requirement: SandboxIgnoreViolations Type Definition
The system SHALL define a `SandboxIgnoreViolations` struct type for specifying violation ignore rules.

#### Scenario: File violation ignore list
- **GIVEN** a SandboxIgnoreViolations with `File: []string{"/tmp/test.txt", "/var/log/app.log"}`
- **WHEN** marshaled to JSON
- **THEN** it produces `{"file": ["/tmp/test.txt", "/var/log/app.log"]}`

#### Scenario: Network violation ignore list
- **GIVEN** a SandboxIgnoreViolations with `Network: []string{"example.com", "api.service.io"}`
- **WHEN** marshaled to JSON
- **THEN** it produces `{"network": ["example.com", "api.service.io"]}`

#### Scenario: Both file and network ignore rules
- **GIVEN** a SandboxIgnoreViolations with both `File` and `Network` populated
- **WHEN** marshaled to JSON
- **THEN** both arrays are present in the JSON output
- **AND** the structure correctly represents all ignore rules

#### Scenario: Empty ignore rules
- **GIVEN** a SandboxIgnoreViolations with empty or nil slices
- **WHEN** marshaled to JSON
- **THEN** empty arrays or omitted fields are serialized appropriately
- **AND** no violations are ignored

---

### Requirement: Go Naming Conventions
All sandbox-related types SHALL follow Go naming conventions while maintaining JSON compatibility with the Python SDK.

#### Scenario: Struct field naming
- **GIVEN** any sandbox-related struct
- **WHEN** examining field names
- **THEN** all exported fields use PascalCase (e.g., `Enabled`, `AutoAllowBashIfSandboxed`)
- **AND** unexported fields use camelCase

#### Scenario: JSON tag mapping
- **GIVEN** any sandbox-related struct field
- **WHEN** marshaling to JSON
- **THEN** JSON field names use camelCase (e.g., `enabled`, `autoAllowBashIfSandboxed`)
- **AND** the JSON format matches the Python SDK's expected format

#### Scenario: Godoc documentation
- **GIVEN** any exported sandbox type or field
- **WHEN** examining the code
- **THEN** comprehensive godoc comments are present
- **AND** comments explain the purpose and platform constraints (e.g., "macOS/Linux only")

---
