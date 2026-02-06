---
hide:
  - toc
---

<div style="text-align: center; margin: 2rem 0;">
  <img src="images/logo.png" alt="txeh" style="max-width: 400px;">
</div>

# txn2/txeh

/etc/hosts management as a Go library and CLI utility. Programmatic and command-line access to add, remove, and query hostname-to-IP mappings. Thread-safe, cross-platform, and built to support tools like [kubefwd](https://github.com/txn2/kubefwd).

[Get Started](getting-started.md){ .md-button .md-button--primary }
[View on GitHub](https://github.com/txn2/txeh){ .md-button }

## Two Ways to Use

<div class="grid cards" markdown>

-   :material-console:{ .lg .middle } **Use the CLI**

    ---

    Manage `/etc/hosts` entries from the command line.

    - Add, remove, list hostnames
    - CIDR range operations
    - Dry run mode for previewing changes

    [:octicons-arrow-right-24: CLI Reference](cli.md)

-   :material-code-braces:{ .lg .middle } **Use the Go Library**

    ---

    Import into your Go application for programmatic hosts file management.

    - Thread-safe with mutex locking
    - In-memory mode from strings
    - Inline comment support

    [:octicons-arrow-right-24: Library docs](library.md)

</div>

## Core Features

<div class="grid cards" markdown>

-   :material-lock-outline:{ .lg .middle } **Thread-Safe**

    ---

    All public methods use mutex locking for safe concurrent access from multiple goroutines.

-   :material-ip-network:{ .lg .middle } **IPv4 & IPv6**

    ---

    Full support for both address families with family-specific lookups.

-   :material-select-group:{ .lg .middle } **CIDR Operations**

    ---

    Bulk add and remove hosts by CIDR range. List all entries in a network.

-   :material-monitor:{ .lg .middle } **Cross-Platform**

    ---

    Linux, macOS, and Windows. Auto-detects the system hosts file location.

</div>

## Quick Install

=== "Homebrew"

    ```bash
    brew install txn2/tap/txeh
    ```

=== "Go Install"

    ```bash
    go install github.com/txn2/txeh/txeh@latest
    ```

=== "Go Library"

    ```bash
    go get github.com/txn2/txeh
    ```

---

[![CI](https://github.com/txn2/txeh/actions/workflows/ci.yml/badge.svg)](https://github.com/txn2/txeh/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/txn2/txeh/branch/master/graph/badge.svg)](https://codecov.io/gh/txn2/txeh)
[![Go Report Card](https://goreportcard.com/badge/github.com/txn2/txeh)](https://goreportcard.com/report/github.com/txn2/txeh)
[![Go Reference](https://pkg.go.dev/badge/github.com/txn2/txeh.svg)](https://pkg.go.dev/github.com/txn2/txeh)
