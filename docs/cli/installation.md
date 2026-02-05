# Installation

## Homebrew (macOS/Linux)

```bash
brew install txn2/tap/txeh
```

## Go Install

Requires Go 1.24 or later.

```bash
go install github.com/txn2/txeh/txeh@latest
```

The binary is installed to `$GOPATH/bin`. Ensure this directory is in your `$PATH`.

## Binary Downloads

Pre-built binaries are available for multiple platforms on the [GitHub Releases](https://github.com/txn2/txeh/releases) page:

| OS | Architectures |
|----|---------------|
| Linux | amd64, arm64, arm, 386 |
| macOS | amd64 (Intel), arm64 (Apple Silicon) |
| Windows | amd64, 386 |

### Linux packages

`.deb`, `.rpm`, and `.apk` packages are also available for each release.

## Build from Source

```bash
git clone https://github.com/txn2/txeh.git
cd txeh
go install ./txeh
```

## Verify Installation

```bash
txeh version
```
