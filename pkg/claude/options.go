package claude

import "context"

// DefaultMaxBufferSize is the default maximum buffer size (1MB) for CLI stdout buffering
// during JSON message accumulation. This matches the Python SDK default.
const DefaultMaxBufferSize = 1024 * 1024

// Options configures the Claude SDK client.
type Options struct {
	// Cancellation and control
	Context context.Context

	// Directory and tool configuration
	AdditionalDirectories []string
	AllowedTools          []string
	DisallowedTools       []string
	Cwd                   string

	// System prompt customization
	// nil for vanilla, literal, or preset+append
	SystemPrompt SystemPromptConfig

	// Permission handling
	CanUseTool     CanUseToolFunc
	PermissionMode PermissionMode
	// Customize which tool is used for permission prompts
	PermissionPromptToolName string

	// Session management
	Continue        bool
	Resume          string
	ResumeSessionAt string
	ForkSession     bool

	// Environment and execution
	Env            map[string]string
	Executable     string // "node", "bun", "deno"
	ExecutableArgs []string
	ExtraArgs      map[string]*string

	// Model configuration
	Model             string
	FallbackModel     string
	MaxThinkingTokens int
	MaxTurns          int

	// Budget and output constraints
	// MaxBudgetUsd enforces a maximum spending limit in USD for API calls during the query session.
	// Precision is maintained to two decimal places (penny precision). A value of 0 or omission
	// means no budget enforcement.
	MaxBudgetUsd float64 `json:"maxBudgetUsd,omitempty"`

	// MaxBufferSize controls the maximum number of bytes for CLI stdout buffering during
	// JSON message accumulation. This prevents unbounded memory growth when processing
	// large outputs from the Claude CLI process.
	//
	// Default: When set to 0 (zero value), DefaultMaxBufferSize (1MB) is used.
	// This default matches the Python SDK and is suitable for most use cases.
	//
	// Customize this value when:
	//   - Processing very large outputs: Increase to accommodate bigger responses
	//   - Memory-constrained environments: Decrease to reduce memory footprint
	//   - Streaming large files or data: Adjust based on expected output sizes
	//
	// Example: MaxBufferSize: 2 * 1024 * 1024 // 2MB buffer for large outputs
	MaxBufferSize int `json:"maxBufferSize,omitempty"`

	// OutputFormat specifies the desired output format for structured outputs.
	// When set, the model's responses will conform to the specified JSON schema format.
	// A nil value uses the default text output format without schema constraints.
	OutputFormat *JsonSchemaOutputFormat `json:"outputFormat,omitempty"`

	// AllowDangerouslySkipPermissions bypasses permission checks when set to true.
	// WARNING: This is a security risk. When enabled, tools execute without user approval prompts.
	// Only use this in controlled environments where the implications of disabling permission
	// checks are fully understood. In production or untrusted environments, keep this false
	// to ensure user approval is required for tool execution.
	AllowDangerouslySkipPermissions bool `json:"allowDangerouslySkipPermissions,omitempty"`

	// Sandbox configures bash command sandboxing for security isolation.
	// When set, bash commands run in a restricted environment on macOS/Linux.
	// A nil value disables sandbox configuration (no sandboxing).
	// See SandboxSettings for detailed configuration options.
	Sandbox *SandboxSettings `json:"sandbox,omitempty"`

	// Plugins configures SDK plugins for extending functionality.
	// Plugins provide custom commands, agents, skills, and hooks that extend Claude Code's capabilities.
	// Currently only local plugins are supported via the 'local' type.
	Plugins []SdkPluginConfig `json:"plugins,omitempty"`

	// MCP servers
	McpServers      map[string]McpServerConfig
	StrictMcpConfig bool

	// Hooks and callbacks
	Hooks  map[HookEvent][]HookCallbackMatcher
	Stderr func(string)

	// Message handling
	IncludePartialMessages bool

	// SDK-specific
	PathToClaudeCodeExecutable string

	// Settings sources
	SettingSources []ConfigScope // validated scopes: local, user, project

	// Settings provides programmatic configuration of Claude Code settings.
	// Accepts either a file path to a settings JSON file or an inline JSON string.
	// This value is passed to the Claude Code CLI via the --settings flag.
	// When provided, these settings override default Claude Code behavior for the session.
	Settings string

	// Agents
	Agents map[string]AgentDefinition

	// User specifies the username to run the Claude Code CLI subprocess as.
	// When set, the subprocess will run with the credentials of the specified user.
	//
	// This feature is Unix-specific and requires appropriate permissions:
	//   - The parent process must have CAP_SETUID and CAP_SETGID capabilities
	//   - Typically requires running as root
	//   - On macOS, requires root or appropriate entitlements
	//
	// When empty (default), the subprocess runs as the current user.
	//
	// Use cases:
	//   - Security isolation: Run untrusted code with reduced privileges
	//   - Multi-tenant environments: Isolate different users' processes
	//   - Container deployments: Run as non-root users following security best practices
	//   - Principle of least privilege: Ensure processes only have necessary permissions
	//
	// Platform support:
	//   - Linux: Fully supported
	//   - macOS: Supported (requires root)
	//   - Windows: Not supported (will be ignored)
	//
	// Example:
	//
	//	opts := &claude.Options{
	//	    User: "nobody",  // Run as unprivileged user
	//	}
	User string
}

// AgentDefinition defines a custom agent.
//
// Tools and DisallowedTools control which tools the agent can use:
//   - Tools: Explicitly lists allowed tools. If set, only these tools are available.
//   - DisallowedTools: Lists tools to exclude from the agent's available tools.
//
// These fields are mutually exclusive in practice - use one or the other, not both.
// If both are specified, the CLI will respect both constraints (allow only Tools,
// but exclude DisallowedTools from that set).
type AgentDefinition struct {
	Description     string   `json:"description"`
	Prompt          string   `json:"prompt"`
	Tools           []string `json:"tools,omitempty"`
	DisallowedTools []string `json:"disallowedTools,omitempty"`
	Model           string   `json:"model,omitempty"`
}

// ModelInfo represents model information.
type ModelInfo struct {
	Value       string `json:"value"`
	DisplayName string `json:"displayName"`
	Description string `json:"description"`
}

// SlashCommand represents available slash commands.
type SlashCommand struct {
	Name         string `json:"name"`
	Description  string `json:"description"`
	ArgumentHint string `json:"argumentHint"`
}

// SandboxIgnoreViolations configures which security violations should be ignored
// within the sandbox environment. This allows specific files or network hosts
// to be accessed without triggering security alerts.
//
// Use this configuration carefully as it weakens sandbox security for the
// specified resources. Only add entries for known safe files or trusted hosts.
//
// Example:
//
//	ignore := &claude.SandboxIgnoreViolations{
//	    File:    []string{"/tmp/safe-file.txt", "/var/log/app.log"},
//	    Network: []string{"api.example.com"},
//	}
type SandboxIgnoreViolations struct {
	// File specifies file paths for which sandbox violations should be ignored.
	// When a sandboxed command accesses these paths, no security violation is raised.
	// Example: []string{"/tmp/safe-file.txt", "/var/log/app.log"}
	File []string `json:"file,omitempty"`

	// Network specifies network hosts for which sandbox violations should be ignored.
	// When sandboxed commands connect to these hosts, no security violation is raised.
	// Example: []string{"example.com", "api.service.io"}
	Network []string `json:"network,omitempty"`
}

// SandboxNetworkConfig configures network access restrictions within the sandbox.
// This allows fine-grained control over socket access, port binding, and proxy
// configuration for sandboxed commands.
//
// Example with Unix socket access and HTTP proxy:
//
//	network := &claude.SandboxNetworkConfig{
//	    AllowUnixSockets: []string{"/var/run/docker.sock"},
//	    HttpProxyPort:    8080,
//	}
type SandboxNetworkConfig struct {
	// AllowUnixSockets specifies Unix socket paths that are accessible within the sandbox.
	// Only the listed sockets can be accessed; all others are restricted.
	// Ignored if AllowAllUnixSockets is true.
	// Example: []string{"/var/run/docker.sock", "/tmp/custom.sock"}
	AllowUnixSockets []string `json:"allowUnixSockets,omitempty"`

	// AllowAllUnixSockets when true permits access to all Unix sockets.
	// This overrides the AllowUnixSockets allowlist.
	// Use with caution as it reduces sandbox security.
	AllowAllUnixSockets bool `json:"allowAllUnixSockets,omitempty"`

	// AllowLocalBinding when true permits sandboxed commands to bind to localhost ports.
	// Note: This option only has effect on macOS systems.
	AllowLocalBinding bool `json:"allowLocalBinding,omitempty"`

	// HttpProxyPort specifies the port number for an HTTP proxy to route outbound
	// HTTP requests from sandboxed commands. A value of 0 means no HTTP proxy.
	HttpProxyPort int `json:"httpProxyPort,omitempty"`

	// SocksProxyPort specifies the port number for a SOCKS5 proxy to route outbound
	// requests from sandboxed commands. A value of 0 means no SOCKS proxy.
	SocksProxyPort int `json:"socksProxyPort,omitempty"`
}

// SandboxSettings configures bash command sandboxing for security control.
// When enabled, bash commands execute in a restricted environment that limits
// file system access, network capabilities, and other system resources.
//
// Sandbox support is available on macOS and Linux systems only.
// On unsupported systems, sandbox configuration is ignored.
//
// Example usage:
//
//	opts := &claude.Options{
//	    Sandbox: &claude.SandboxSettings{
//	        Enabled:                  true,
//	        AutoAllowBashIfSandboxed: true,
//	        ExcludedCommands:         []string{"docker", "git"},
//	    },
//	}
type SandboxSettings struct {
	// Enabled activates bash command sandboxing.
	// When true, bash commands run in a restricted sandbox environment.
	// Supported on macOS and Linux only.
	Enabled bool `json:"enabled,omitempty"`

	// AutoAllowBashIfSandboxed when true automatically approves sandboxed bash
	// commands without user prompts. This reduces friction for sandboxed operations
	// since the sandbox provides security isolation.
	AutoAllowBashIfSandboxed bool `json:"autoAllowBashIfSandboxed,omitempty"`

	// ExcludedCommands lists commands that should run outside the sandbox.
	// Use this for commands that require elevated privileges or access to
	// resources that would be blocked by the sandbox.
	// Example: []string{"docker", "git", "npm"}
	ExcludedCommands []string `json:"excludedCommands,omitempty"`

	// AllowUnsandboxedCommands when true permits the dangerouslyDisableSandbox
	// flag to be used in tool calls, allowing specific commands to bypass
	// sandbox restrictions. Use with caution as this reduces security.
	AllowUnsandboxedCommands bool `json:"allowUnsandboxedCommands,omitempty"`

	// Network configures network access restrictions within the sandbox.
	// A nil value uses default network restrictions.
	Network *SandboxNetworkConfig `json:"network,omitempty"`

	// IgnoreViolations configures which security violations should be ignored.
	// A nil value means all violations are reported.
	IgnoreViolations *SandboxIgnoreViolations `json:"ignoreViolations,omitempty"`

	// EnableWeakerNestedSandbox when true uses a weaker sandbox configuration
	// suitable for Docker containers or other nested virtualization environments.
	// This allows the sandbox to function within environments that already have
	// their own security restrictions.
	EnableWeakerNestedSandbox bool `json:"enableWeakerNestedSandbox,omitempty"`
}
