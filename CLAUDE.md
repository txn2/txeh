# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

txeh is a Go library and CLI utility for managing /etc/hosts file entries. It provides programmatic and command-line access to add, remove, and query hostname-to-IP mappings. Originally built to support [kubefwd](https://github.com/txn2/kubefwd) for Kubernetes port-forwarding.

**Key Features:**
- Thread-safe operations (mutex-protected for concurrent use)
- IPv4 and IPv6 support
- CIDR range operations for bulk add/remove
- Inline comment support for tracking entry sources
- Preserves comments and file formatting
- Cross-platform (Linux, macOS, Windows)

## Quick Reference

```bash
# Full verification (run before every commit)
make verify

# Run all tests including integration
make verify-full

# Format, lint, test individually
make fmt
make lint
make test-unit
make coverage
make security
```

## Architecture

```
┌─────────────────┐
│   txeh/txeh.go  │  CLI entry point (Cobra-based)
│    (cmd/)       │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│    txeh.go      │  Core library (public API)
│   (package)     │
└─────────────────┘
         │
         ▼
┌─────────────────┐
│  /etc/hosts     │  System hosts file
└─────────────────┘
```

### Key Components

1. **Library** (`txeh.go`): Core functionality for parsing, modifying, and rendering hosts files
   - `Hosts` struct: Thread-safe wrapper with mutex protection
   - `HostFileLine`: Represents a single line (ADDRESS, COMMENT, EMPTY, UNKNOWN)
   - `HostsConfig`: Configuration for read/write paths or raw text input

2. **CLI** (`txeh/` directory): Cobra-based command-line interface
   - Entry point: `txeh/txeh.go`
   - Commands in `txeh/cmd/`: add, remove, list, show, version
   - Global flags: `--dryrun`, `--quiet`, `--read`, `--write`, `--max-hosts-per-line`

## Code Standards

### Style
- Follow [Google Go Style Guide](https://google.github.io/styleguide/go/)
- All code must pass `golangci-lint` with project config
- Maximum cyclomatic complexity: 15
- Maximum cognitive complexity: 20

### Documentation
- All exported functions, types, and packages require doc comments
- Doc comments must start with the name being documented
- Comments should end with a period

### Error Handling
- Always check and handle errors (no `_` for errors)
- Wrap errors at boundaries: `fmt.Errorf("context: %w", err)`
- Use `errors.Is` and `errors.As`, never string comparison
- Define sentinel errors for expected conditions

```go
// Example error handling pattern
func (h *Hosts) Save() error {
    if err := h.writeFile(); err != nil {
        return fmt.Errorf("save hosts file: %w", err)
    }
    return nil
}
```

### Testing
- Test coverage minimum: 82%
- Use table-driven tests for multiple cases
- Use `t.Parallel()` for independent tests
- Race detection required: `go test -race ./...`

## Security Requirements

- Never log secrets, tokens, or passwords
- Validate all inputs at system boundaries
- File paths are sanitized (hosts file manager reads arbitrary paths by design)
- Run `gosec ./...` before commits

## Dependencies

- **Prefer stdlib** over external packages
- All dependencies pinned to exact versions
- Run `go mod verify` before committing
- Regular `govulncheck ./...` for vulnerability scanning

## Project-Specific Patterns

### Adding Hosts
```go
// Always use the thread-safe methods
hosts.AddHost("127.0.0.1", "myhost")
hosts.AddHosts("127.0.0.1", []string{"host1", "host2"})
hosts.AddHostWithComment("127.0.0.1", "myhost", "added by myapp")
```

### Removing Entries
```go
hosts.RemoveHost("myhost")
hosts.RemoveAddress("127.0.0.1")
hosts.RemoveCIDRs([]string{"10.0.0.0/8"})  // Returns error
hosts.RemoveByComment("added by myapp")
```

### Querying
```go
ips := hosts.ListAddressesByHost("myhost", true)  // true = lookup all
hostnames := hosts.ListHostsByIP("127.0.0.1")
entries := hosts.ListHostsByComment("myapp")
```

## Common Commands

```bash
# Development
make build          # Build CLI binary
make test-unit      # Run unit tests
make lint           # Run linter
make fmt            # Format code

# Before committing
make verify         # Format, lint, test, security scan

# Full CI verification
make verify-full    # Includes integration tests

# Debugging
go test -v -run TestSpecific ./...  # Run specific test

# Build release
goreleaser --skip-publish --clean --skip-validate
```

## File Naming Conventions

| Pattern | Purpose |
|---------|---------|
| `*_test.go` | Test files (same package) |
| `txeh/cmd/*.go` | CLI command implementations |

## Git Commit Guidelines

Use conventional commits:
- `feat:` New features
- `fix:` Bug fixes
- `docs:` Documentation changes
- `test:` Test additions/changes
- `refactor:` Code refactoring
- `perf:` Performance improvements
- `security:` Security fixes
- `chore:` Maintenance tasks

**Important:** A human must review and approve every line of code before commit.

## Troubleshooting

### Common Issues

**"package not found" errors:**
```bash
go mod tidy
go mod verify
```

**Linter failures:**
```bash
golangci-lint run --fix ./...  # Auto-fix where possible
```

**Permission denied on /etc/hosts:**
```bash
sudo txeh add 127.0.0.1 myhost  # Requires root/admin
```

**Windows hosts file location:**
The library auto-detects via `%SystemRoot%\System32\drivers\etc\hosts`
