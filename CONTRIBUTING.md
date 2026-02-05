# Contributing to txeh

Thank you for your interest in contributing to txeh!

## Development Setup

1. **Go 1.21+** required
2. Install tools:
   ```bash
   go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
   go install golang.org/x/vuln/cmd/govulncheck@latest
   go install github.com/securego/gosec/v2/cmd/gosec@latest
   ```
3. Clone and build:
   ```bash
   git clone https://github.com/txn2/txeh.git
   cd txeh
   make build
   ```

## Pull Request Process

1. Fork the repository
2. Create a feature branch: `git checkout -b feat/your-feature`
3. Make changes with tests
4. Run `make verify` before committing
5. Submit PR with description

## Code Style

- Run `gofmt -s -w .` before committing
- All exported functions need doc comments
- Error messages should be lowercase, no punctuation
- Follow [Effective Go](https://go.dev/doc/effective_go) guidelines

## Commit Messages

Use conventional commits:
- `feat:` New features
- `fix:` Bug fixes
- `docs:` Documentation changes
- `test:` Test additions/changes
- `refactor:` Code refactoring
- `chore:` Maintenance tasks

Examples:
```
feat: add support for IPv6 CIDR ranges
fix: handle empty hostname in AddHost
docs: update README with new CLI flags
test: add coverage for RemoveAddress edge cases
```

## Testing

- Add tests for new functionality
- Maintain 80%+ coverage
- Run `go test -race ./...` to check for race conditions

```bash
# Run all tests
make test

# Run tests with coverage
make coverage

# Run linting
make lint

# Run full verification suite
make verify
```

## Reporting Issues

- Check existing issues before creating a new one
- Use the issue templates provided
- Include version, OS, and reproduction steps

## Security Vulnerabilities

Please report security vulnerabilities privately. See [SECURITY.md](SECURITY.md) for details.

## License

By contributing, you agree that your contributions will be licensed under the Apache License 2.0.
