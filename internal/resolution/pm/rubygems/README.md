# RubyGems resolution logic

The way resolution of RubyGems lock files works is as follows:

1. Run `bundle install --quiet` to install all dependencies and generate the lock file

The `--quiet` flag is used to suppress progress messages for cleaner output.

Generated `Gemfile.lock` file is then uploaded together with `Gemfile` for scanning.

## Notes

- RubyGems uses `bundle install` to resolve dependencies and create/update the lock file
- The manifest file is `Gemfile`
- The lock file is `Gemfile.lock`
- Bundler respects the existing `Gemfile.lock` if present to ensure consistent dependency versions
