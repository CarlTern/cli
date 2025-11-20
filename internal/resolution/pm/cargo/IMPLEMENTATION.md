# Cargo Package Manager Implementation

This implementation mirrors the Composer and CocoaPods package manager structure for Cargo (Rust) dependency resolution.

## Files Created

### Core Implementation
- **pm.go**: Defines the Cargo package manager with manifest pattern `Cargo\.toml$`
- **cmd_factory.go**: Creates the `cargo generate-lockfile --offline` command
- **job.go**: Handles job execution and error handling specific to Cargo
- **strategy.go**: Creates jobs for each Cargo.toml found

### Tests
- **pm_test.go**: Tests for package manager name and manifest patterns
- **cmd_factory_test.go**: Tests for command generation
- **job_test.go**: Tests for job execution and error handling
- **strategy_test.go**: Tests for strategy invocation

### Documentation
- **README.md**: Explains the resolution process
- **IMPLEMENTATION.md**: This file

### Test Data
- **testdata/cmd_factory_mock.go**: Mock factory for testing

## Key Differences from Other Package Managers

### Command Comparison
- **Composer**: `composer update` with various flags
- **CocoaPods**: `pod install --no-repo-update`
- **Cargo**: `cargo generate-lockfile --offline`
  - Uses `generate-lockfile` specifically to create the lock file
  - `--offline` prevents network access and uses only local cache

### Manifest/Lock Files
- **Manifest**: `Cargo.toml`
- **Lock file**: `Cargo.lock`

### Error Handling
Cargo-specific error patterns:

#### Cargo Errors:
- Executable not found
- Network errors (failed to download, unable to update registry)
- Invalid crate names
- Version resolution conflicts
- Incompatible Rust version
- TOML parse errors

## Why Certain Functions Are Included

### pm.go
- **Name()**: Returns "cargo" for package manager identification
- **Manifests()**: Returns regex pattern `Cargo\.toml$` to identify Rust manifest files

### cmd_factory.go
- **ICmdFactory interface**: Allows mocking in tests
- **IExecPath interface**: Allows mocking exec.LookPath in tests
- **MakeInstallCmd()**: Builds the command to run
  - Uses `cargo generate-lockfile` (not `cargo build` or `cargo fetch`)
  - Adds `--offline` to prevent network access and use only local cache

### job.go
- **Install()**: Returns whether this job should run the install command
- **Run()**: Main execution method that runs the lockfile generation
- **runInstallCmd()**: Executes the cargo generate-lockfile command
- **handleError()**: Matches error patterns and adds helpful documentation
- **addDocumentation()**: Provides user-friendly error messages for common issues
- **Individual error documentation methods**: Each common error type gets specific guidance

### strategy.go
- **Invoke()**: Creates a job for each Cargo.toml found in the project
- **NewStrategy()**: Constructor that takes a list of files

## Cargo-Specific Considerations

1. **Offline flag**: Cargo can be slow when accessing the network, so we use `--offline` to rely on the local cache
2. **Generate-lockfile**: This is the specific command to create a lock file without building
3. **Registry cache**: Cargo maintains a local registry cache that can be populated with `cargo fetch`
4. **Rust version requirements**: Some crates require specific Rust versions
5. **TOML format**: Cargo.toml must be valid TOML syntax

## Differences from Other Ecosystems

- **Rust specific**: Unlike PHP (Composer) or Ruby (CocoaPods), Cargo is Rust-specific
- **No build**: We only generate the lock file, we don't compile code
- **Registry**: Uses crates.io as the default registry
- **Offline mode**: Cargo has robust offline support via local cache

## Testing

All test files follow the same pattern as Composer and CocoaPods:
- Unit tests for each component
- Mock command factory for testing without actual cargo execution
- Error handling tests with realistic Cargo error messages
- Strategy tests for single and multiple files
