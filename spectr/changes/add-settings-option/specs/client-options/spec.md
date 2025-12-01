# Client Options - Delta Specification

## ADDED Requirements

### Requirement: Settings Option
The Options struct SHALL support a `Settings` field that allows programmatic configuration of Claude Code settings.

#### Scenario: Settings with file path
- **GIVEN** an Options struct with `Settings: "/path/to/settings.json"`
- **WHEN** the query is initialized
- **THEN** the CLI is invoked with `--settings /path/to/settings.json`
- **AND** Claude Code loads settings from the specified file

#### Scenario: Settings with inline JSON
- **GIVEN** an Options struct with `Settings: "{\"outputStyle\": \"compact\"}"`
- **WHEN** the query is initialized
- **THEN** the CLI is invoked with `--settings {"outputStyle": "compact"}`
- **AND** Claude Code applies the inline JSON settings

#### Scenario: No settings specified
- **GIVEN** an Options struct with `Settings: ""`
- **WHEN** the query is initialized
- **THEN** the `--settings` flag is not included in CLI arguments
- **AND** Claude Code uses default settings behavior

#### Scenario: Settings field is documented
- **WHEN** reviewing the `Settings` field godoc
- **THEN** it explains the field accepts either a file path or inline JSON string
- **AND** it describes the purpose of programmatic settings configuration
- **AND** it references the Python SDK's equivalent functionality for consistency
