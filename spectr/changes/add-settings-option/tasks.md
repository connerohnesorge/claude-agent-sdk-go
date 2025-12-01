# Implementation Tasks

## 1. Add Settings Field to Options Struct
- [ ] 1.1 Add `Settings string` field to `Options` struct in `pkg/claude/options.go`
- [ ] 1.2 Add godoc comment documenting the field purpose and accepted values (file path or JSON string)

## 2. Update CLI Argument Building
- [ ] 2.1 Update `buildArgs()` method in `pkg/claude/query.go` to include `--settings` flag
- [ ] 2.2 Add conditional logic to append `--settings` flag only when `q.opts.Settings != ""`

## 3. Testing
- [ ] 3.1 Add unit test for Options struct with Settings field
- [ ] 3.2 Add unit test for buildArgs() including Settings in arguments
- [ ] 3.3 Add integration test with file path settings
- [ ] 3.4 Add integration test with inline JSON settings (if feasible)

## 4. Documentation
- [ ] 4.1 Update relevant examples to show Settings usage (if applicable)
- [ ] 4.2 Verify godoc comments are comprehensive
