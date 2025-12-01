## ADDED Requirements

### Requirement: MaxBufferSize Option
The Options struct SHALL support a `MaxBufferSize` field that controls the maximum bytes allowed when buffering CLI stdout during JSON message accumulation, preventing unbounded memory growth.

#### Scenario: Default buffer size applied
- **GIVEN** an Options struct with `MaxBufferSize: 0` or omitted
- **WHEN** the query initializes the transport
- **THEN** a default buffer size of 1MB (1024 * 1024 bytes) is used
- **AND** the transport enforces this limit during JSON buffering

#### Scenario: Custom buffer size enforced
- **GIVEN** an Options struct with `MaxBufferSize: 512000` (500KB)
- **WHEN** the CLI outputs a large incomplete JSON message exceeding 500KB
- **THEN** the transport stops buffering when the limit is reached
- **AND** an error is returned indicating buffer size exceeded
- **AND** the error includes the limit value and actual buffer size

#### Scenario: Large output within buffer limit
- **GIVEN** an Options struct with `MaxBufferSize: 2097152` (2MB)
- **WHEN** the CLI outputs a 1.5MB complete JSON message
- **THEN** the message is successfully buffered and parsed
- **AND** no buffer size errors occur

#### Scenario: Buffer size protects against OOM
- **GIVEN** an Options struct with `MaxBufferSize: 1048576` (1MB)
- **WHEN** the CLI produces an extremely large incomplete JSON (10MB+)
- **THEN** buffering stops at the 1MB limit
- **AND** memory usage is bounded to the configured limit
- **AND** the process does not encounter out-of-memory errors

#### Scenario: Error provides diagnostic information
- **GIVEN** an Options struct with a MaxBufferSize limit
- **WHEN** the buffer size is exceeded during JSON accumulation
- **THEN** the error message includes:
  - The configured buffer size limit
  - The actual buffer size when exceeded
  - Clear indication of which operation failed
- **AND** the error is of type `ErrBufferSizeExceeded` for type-safe handling
