# Change: Add Sandbox Configuration Support

## Why

The Go SDK currently lacks sandbox configuration capabilities that are available in the Python SDK. This prevents Go developers from configuring bash command sandboxing for security, which is critical for controlling execution environments, network access, and file system isolation when running AI-generated commands.

Without sandbox configuration, users cannot:
- Enable bash command sandboxing on macOS/Linux
- Configure network access restrictions and proxy settings
- Specify which commands should run outside the sandbox
- Define file paths and network hosts to ignore violations for
- Control nested sandbox behavior for Docker environments

This creates a feature parity gap with the Python SDK and limits security control options for Go users.

## What Changes

Add three new struct types to `pkg/claude/options.go` matching the Python SDK's sandbox configuration:

1. **SandboxSettings** - Main sandbox configuration with fields:
   - `Enabled` (bool) - Enable bash sandboxing (macOS/Linux only)
   - `AutoAllowBashIfSandboxed` (bool) - Auto-approve sandboxed bash commands
   - `ExcludedCommands` ([]string) - Commands to run outside sandbox
   - `AllowUnsandboxedCommands` (bool) - Allow dangerouslyDisableSandbox flag
   - `Network` (*SandboxNetworkConfig) - Network access configuration
   - `IgnoreViolations` (*SandboxIgnoreViolations) - Violation ignore rules
   - `EnableWeakerNestedSandbox` (bool) - Use weaker sandbox for Docker/nested scenarios

2. **SandboxNetworkConfig** - Network configuration with fields:
   - `AllowUnixSockets` ([]string) - Unix socket paths accessible in sandbox
   - `AllowAllUnixSockets` (bool) - Allow all Unix sockets
   - `AllowLocalBinding` (bool) - Bind to localhost ports (macOS only)
   - `HttpProxyPort` (int) - HTTP proxy port for outbound requests
   - `SocksProxyPort` (int) - SOCKS5 proxy port for outbound requests

3. **SandboxIgnoreViolations** - Violation ignore rules with fields:
   - `File` ([]string) - File paths to ignore violations for
   - `Network` ([]string) - Network hosts to ignore violations for

Add `Sandbox *SandboxSettings` field to the `Options` struct to enable sandbox configuration.

All struct fields will use appropriate Go naming conventions (PascalCase) and JSON tags for proper serialization when passed to the Claude CLI.

## Impact

- **Affected specs**: `client-options` (ADDED Requirements)
- **Affected code**:
  - `/home/connerohnesorge/Documents/001Repos/claude-agent-sdk-go/pkg/claude/options.go` - Add new types and Options field
  - Future examples demonstrating sandbox configuration usage (optional)
- **Breaking changes**: None - this is an additive change
- **Python SDK parity**: Achieves feature parity with Python SDK v0.1.0+ sandbox configuration
- **TypeScript SDK**: Should verify alignment (TypeScript SDK may have similar structures)
