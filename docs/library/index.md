# Go Library Overview

txeh provides a Go library for programmatic `/etc/hosts` file management. It was originally built to support [kubefwd](https://github.com/txn2/kubefwd) for Kubernetes service forwarding.

## Key Features

- **Thread-safe** -- All public methods use mutex locking for concurrent access
- **Cross-platform** -- Linux, macOS, Windows (auto-detects hosts file location)
- **Preserves formatting** -- Comments, empty lines, and structure are maintained
- **IPv4 & IPv6** -- Full support for both address families
- **CIDR operations** -- Bulk removal by network range
- **Inline comments** -- Group and track entries with comments
- **In-memory mode** -- Parse hosts data from strings without file I/O

## Installation

```bash
go get github.com/txn2/txeh
```

## Quick Example

```go
hosts, err := txeh.NewHostsDefault()
if err != nil {
    panic(err)
}

hosts.AddHost("127.0.0.1", "myapp.local")
hosts.Save()
```

See [Quick Start](quickstart.md) for detailed examples and [API Reference](api.md) for the complete API.
