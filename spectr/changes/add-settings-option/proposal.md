# Change: Add Settings Option to Go SDK

## Why
The Go SDK is missing the `Settings` field available in the Python SDK, preventing programmatic configuration of Claude Code settings. This limits the ability to override default settings, provide custom configurations, and control behavior without modifying global settings files.

## What Changes
- Add `Settings string` field to the `Options` struct in `pkg/claude/options.go`
- Update the `buildArgs()` method in `pkg/claude/query.go` to pass the settings value to the CLI via the `--settings` flag
- The `Settings` field accepts either:
  - A file path to a settings JSON file
  - An inline JSON string containing settings configuration

## Impact
- Affected specs: `client-options`
- Affected code:
  - `pkg/claude/options.go` - Add Settings field to Options struct
  - `pkg/claude/query.go` - Update buildArgs() to include --settings flag
- **BREAKING**: None - this is an additive change maintaining backward compatibility
- Brings Go SDK feature parity with Python SDK's settings option
