# Cargo resolution logic

The way resolution of Cargo lock files works is as follows:

1. Run `cargo generate-lockfile --offline` to generate the lock file from the manifest

Generated `Cargo.lock` file is then uploaded together with `Cargo.toml` for scanning.