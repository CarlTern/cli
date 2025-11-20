# CocoaPods resolution logic

The way resolution of CocoaPods lock files works is as follows:

1. Run `pod install --no-repo-update` in order to install all dependencies

The `--no-repo-update` flag is used to prevent updating the local specs repository, which makes the process faster and avoids potential network issues.

Generated `Podfile.lock` file is then uploaded together with `Podfile` for scanning.