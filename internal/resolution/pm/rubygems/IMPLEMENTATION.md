# RubyGems Package Manager Implementation

This implementation mirrors the Cargo, CocoaPods, and Composer package manager structure for RubyGems (Ruby/Bundler) dependency resolution.

## Files Created

### Core Implementation
- **pm.go**: Defines the RubyGems package manager with manifest pattern `Gemfile$`
- **cmd_factory.go**: Creates the `bundle install --quiet` command
- **job.go**: Handles job execution and error handling specific to Bundler/RubyGems
- **strategy.go**: Creates jobs for each Gemfile found

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
- **RubyGems**: `bundle install --quiet`
  - Uses `bundle install` to resolve dependencies and create/update the lock file
  - `--quiet` suppresses progress messages for cleaner output

### Manifest/Lock Files
- **Manifest**: `Gemfile`
- **Lock file**: `Gemfile.lock`

### Error Handling
RubyGems/Bundler-specific error patterns:

#### RubyGems/Bundler Errors:
- Bundler executable not found
- Gem not found in repository
- Version conflicts between gem requirements
- Network errors
- Ruby version requirements
- Gemfile parse errors (invalid Ruby syntax)

## Why Certain Functions Are Included

### pm.go
- **Name()**: Returns "rubygems" for package manager identification
- **Manifests()**: Returns regex pattern `Gemfile$` to identify Ruby manifest files

### cmd_factory.go
- **ICmdFactory interface**: Allows mocking in tests
- **IExecPath interface**: Allows mocking exec.LookPath in tests
- **MakeInstallCmd()**: Builds the command to run
  - Uses `bundle install` to resolve and install dependencies
  - Adds `--quiet` to suppress verbose output

### job.go
- **Install()**: Returns whether this job should run the install command
- **Run()**: Main execution method that runs the bundle install
- **runInstallCmd()**: Executes the bundle install command
- **handleError()**: Matches error patterns and adds helpful documentation
- **addDocumentation()**: Provides user-friendly error messages for common issues
- **Individual error documentation methods**: Each common error type gets specific guidance

### strategy.go
- **Invoke()**: Creates a job for each Gemfile found in the project
- **NewStrategy()**: Constructor that takes a list of files

## RubyGems-Specific Considerations

1. **Bundler**: The standard tool for managing Ruby gem dependencies
2. **Gemfile**: Written in Ruby DSL (Domain Specific Language)
3. **Lock file**: Bundler respects existing `Gemfile.lock` for consistent versions
4. **RubyGems.org**: Default gem repository
5. **Ruby version**: Some gems require specific Ruby versions
6. **Native extensions**: Some gems have C extensions that need compilation

## Differences from Other Ecosystems

- **Ruby DSL**: Gemfile uses Ruby syntax, not TOML/JSON/YAML
- **Bundler**: Separate tool from Ruby itself (like Cargo for Rust)
- **Version locking**: `Gemfile.lock` ensures exact versions across environments
- **Gem groups**: Supports development, test, production groups

## Testing

All test files follow the same pattern as Cargo, CocoaPods, and Composer:
- Unit tests for each component
- Mock command factory for testing without actual bundler execution
- Error handling tests with realistic Bundler error messages
- Strategy tests for single and multiple files
