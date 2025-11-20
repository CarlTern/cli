# CocoaPods Package Manager Implementation

This implementation mirrors the Composer package manager structure for CocoaPods dependency resolution.

## Files Created

### Core Implementation
- **pm.go**: Defines the CocoaPods package manager with manifest pattern `Podfile$`
- **cmd_factory.go**: Creates the `pod install --no-repo-update` command
- **job.go**: Handles job execution and error handling specific to CocoaPods
- **strategy.go**: Creates jobs for each Podfile found

### Tests
- **pm_test.go**: Tests for package manager name and manifest patterns
- **cmd_factory_test.go**: Tests for command generation
- **job_test.go**: Tests for job execution and error handling
- **strategy_test.go**: Tests for strategy invocation

### Documentation
- **README.md**: Explains the resolution process

### Test Data
- **testdata/cmd_factory_mock.go**: Mock factory for testing

## Key Differences from Composer

### Command Differences
- **Composer**: `composer update` with flags like `--no-interaction`, `--no-scripts`, `--ignore-platform-reqs`, `--no-autoloader`, `--no-install`, `--no-plugins`, `--no-audit`
- **CocoaPods**: `pod install --no-repo-update`
  - Uses `install` instead of `update` because it respects the existing lock file
  - `--no-repo-update` prevents updating the local specs repository for faster execution

### Manifest/Lock Files
- **Composer**: `composer.json` (manifest), `composer.lock` (lock)
- **CocoaPods**: `Podfile` (manifest), `Podfile.lock` (lock)

### Error Handling
Both implementations handle similar error types but with package manager-specific patterns:

#### Composer Errors:
- Executable not found
- Missing PHP extensions (phar)
- Invalid requirement format
- Network issues
- Invalid versions
- Dependency resolution failures

#### CocoaPods Errors:
- Executable not found
- Specs repository not found
- Network errors
- Invalid pod names
- Version conflicts
- Xcode deployment target issues

## Why Certain Functions Are Included

### pm.go
- **Name()**: Returns the package manager name for identification
- **Manifests()**: Returns regex patterns to identify CocoaPods manifest files (Podfile)

### cmd_factory.go
- **ICmdFactory interface**: Allows mocking in tests
- **IExecPath interface**: Allows mocking exec.LookPath in tests
- **MakeInstallCmd()**: Builds the actual shell command to run
  - Uses `pod install` (not `update`) to respect existing lock files
  - Adds `--no-repo-update` to avoid unnecessary network calls

### job.go
- **Install()**: Returns whether this job should run the install command
- **Run()**: Main execution method that runs the install command
- **runInstallCmd()**: Executes the pod install command
- **handleError()**: Matches error patterns and adds helpful documentation
- **addDocumentation()**: Provides user-friendly error messages for common issues
- **Individual error documentation methods**: Each common error type gets specific guidance

### strategy.go
- **Invoke()**: Creates a job for each Podfile found in the project
- **NewStrategy()**: Constructor that takes a list of files

## CocoaPods-Specific Considerations

1. **No repo update flag**: CocoaPods can be slow when updating the specs repository, so we skip it
2. **Install vs Update**: `pod install` is the correct command for CI/dependency resolution as it respects the lock file
3. **Platform requirements**: Unlike Composer's PHP platform requirements, CocoaPods has Xcode/iOS deployment target requirements
4. **Specs repository**: CocoaPods maintains a local repository of pod specifications that can become out of sync

## Testing

All test files follow the same pattern as Composer:
- Unit tests for each component
- Mock command factory for testing without actual pod execution
- Error handling tests with realistic error messages
- Strategy tests for single and multiple files
