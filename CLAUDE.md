# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

txeh is a Go library and CLI utility for managing /etc/hosts file entries. It provides programmatic and command-line access to add, remove, and query hostname-to-IP mappings. Originally built to support [kubefwd](https://github.com/txn2/kubefwd) for Kubernetes port-forwarding.

## Common Commands

```bash
# Run tests
go test ./...

# Run a single test
go test -run TestMethodName

# Build and run CLI from source
go run ./txeh/txeh.go

# Install CLI locally
go install ./txeh

# Build test release (requires goreleaser)
goreleaser --skip-publish --clean --skip-validate
```

## Architecture

The codebase has two main components:

1. **Library** (`txeh.go`): Core functionality for parsing, modifying, and rendering hosts files
   - `Hosts` struct: Thread-safe wrapper around parsed hosts file with mutex protection
   - `HostFileLine`: Represents a single line (address entry, comment, or empty)
   - `HostsConfig`: Configuration for read/write paths or raw text input
   - Key methods: `AddHost`, `RemoveHost`, `RemoveAddress`, `RemoveCIDRs`, `ListHostsByIP`, `ListAddressesByHost`

2. **CLI** (`txeh/` directory): Cobra-based command-line interface
   - Entry point: `txeh/txeh.go`
   - Commands in `txeh/cmd/`: add, remove (host/ip/cidr), list (hosts/ip/cidr), show, version
   - Global flags: `--dryrun`, `--quiet`, `--read`, `--write`

## Key Design Details

- Hosts file parsing preserves comments and empty lines (line types: UNKNOWN, EMPTY, COMMENT, ADDRESS)
- All hostnames and addresses are normalized to lowercase
- Thread-safe: all public methods use mutex locking
- Supports both IPv4 and IPv6 (IPFamilyV4, IPFamilyV6)
- Cross-platform: auto-detects Windows hosts file location via SystemRoot env var
- `RawText` config option allows in-memory parsing without file I/O (useful for testing)
- CIDR range operations for bulk removal of address blocks
