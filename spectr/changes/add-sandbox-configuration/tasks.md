# Implementation Tasks

## 1. Type Definitions
- [x] 1.1 Add `SandboxIgnoreViolations` struct to `pkg/claude/options.go`
  - Include `File []string` field with JSON tag `file,omitempty`
  - Include `Network []string` field with JSON tag `network,omitempty`
  - Add comprehensive godoc explaining violation ignore rules
- [x] 1.2 Add `SandboxNetworkConfig` struct to `pkg/claude/options.go`
  - Include `AllowUnixSockets []string` field with JSON tag `allowUnixSockets,omitempty`
  - Include `AllowAllUnixSockets bool` field with JSON tag `allowAllUnixSockets,omitempty`
  - Include `AllowLocalBinding bool` field with JSON tag `allowLocalBinding,omitempty` (note: macOS only)
  - Include `HttpProxyPort int` field with JSON tag `httpProxyPort,omitempty`
  - Include `SocksProxyPort int` field with JSON tag `socksProxyPort,omitempty`
  - Add comprehensive godoc explaining network configuration options
- [x] 1.3 Add `SandboxSettings` struct to `pkg/claude/options.go`
  - Include `Enabled bool` field with JSON tag `enabled,omitempty`
  - Include `AutoAllowBashIfSandboxed bool` field with JSON tag `autoAllowBashIfSandboxed,omitempty`
  - Include `ExcludedCommands []string` field with JSON tag `excludedCommands,omitempty`
  - Include `AllowUnsandboxedCommands bool` field with JSON tag `allowUnsandboxedCommands,omitempty`
  - Include `Network *SandboxNetworkConfig` field with JSON tag `network,omitempty`
  - Include `IgnoreViolations *SandboxIgnoreViolations` field with JSON tag `ignoreViolations,omitempty`
  - Include `EnableWeakerNestedSandbox bool` field with JSON tag `enableWeakerNestedSandbox,omitempty`
  - Add comprehensive godoc explaining sandbox configuration and platform support (macOS/Linux)

## 2. Options Integration
- [x] 2.1 Add `Sandbox *SandboxSettings` field to the `Options` struct in `pkg/claude/options.go`
  - Use JSON tag `sandbox,omitempty`
  - Add godoc comment explaining sandbox configuration usage
  - Position the field logically within the Options struct (near security/permission fields)

## 3. Documentation
- [x] 3.1 Add godoc examples for sandbox configuration
  - Example of enabling sandbox with auto-allow
  - Example of network configuration with Unix sockets
  - Example of excluded commands configuration
- [x] 3.2 Update any relevant documentation comments
  - Ensure platform constraints (macOS/Linux) are clearly documented
  - Note that sandbox is optional and nil means no sandboxing

## 4. Testing
- [x] 4.1 Add unit tests for SandboxSettings JSON serialization in `test/unit/`
  - Test marshaling with all fields populated
  - Test marshaling with minimal configuration
  - Test omitempty behavior for optional fields
- [x] 4.2 Add unit tests for SandboxNetworkConfig JSON serialization
  - Test proxy port configurations
  - Test Unix socket allowlist serialization
  - Test AllowAllUnixSockets flag
- [x] 4.3 Add unit tests for SandboxIgnoreViolations JSON serialization
  - Test file violation ignore list
  - Test network violation ignore list
  - Test combined file and network rules
- [x] 4.4 Add integration test demonstrating sandbox configuration usage (optional)
  - Test that Options with Sandbox field passes to Claude CLI correctly
  - Verify JSON structure matches Python SDK format

## 5. Validation
- [x] 5.1 Run `go build ./...` to ensure code compiles
- [x] 5.2 Run `go test ./...` to ensure all tests pass
- [x] 5.3 Run `golangci-lint run` to ensure linting passes
- [x] 5.4 Verify JSON serialization matches Python SDK format expectations
- [x] 5.5 Update tasks.md to mark all items as complete
