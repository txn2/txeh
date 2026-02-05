# txeh

<div class="home-hero" markdown>
  <img src="images/logo.png" alt="txeh" class="hero-logo">
</div>

<div class="home-intro" markdown>

<span class="hero-highlight">/etc/hosts management as a Go library and CLI utility.</span>

Programmatic and command-line access to add, remove, and query hostname-to-IP mappings.
Thread-safe, cross-platform, and built to support tools like [kubefwd](https://github.com/txn2/kubefwd).

</div>

<div class="home-buttons">
  <a href="getting-started/" class="btn btn-primary">Get Started</a>
  <a href="https://github.com/txn2/txeh" class="btn btn-secondary">GitHub</a>
</div>

<div class="home-features">
  <div class="feature">
    <div class="feature-icon">:material-lock-outline:</div>
    <div><strong>Thread-Safe</strong><span>Mutex-protected operations for safe concurrent use</span></div>
  </div>
  <div class="feature">
    <div class="feature-icon">:material-ip-network:</div>
    <div><strong>IPv4 & IPv6</strong><span>Full support for both address families</span></div>
  </div>
  <div class="feature">
    <div class="feature-icon">:material-select-group:</div>
    <div><strong>CIDR Operations</strong><span>Bulk add/remove by CIDR range</span></div>
  </div>
  <div class="feature">
    <div class="feature-icon">:material-monitor:</div>
    <div><strong>Cross-Platform</strong><span>Linux, macOS, and Windows</span></div>
  </div>
</div>

<div class="home-install">

```bash
brew install txn2/tap/txeh
```

</div>

<div class="home-footer">

[![CI](https://github.com/txn2/txeh/actions/workflows/ci.yml/badge.svg)](https://github.com/txn2/txeh/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/txn2/txeh/branch/master/graph/badge.svg)](https://codecov.io/gh/txn2/txeh)
[![Go Report Card](https://goreportcard.com/badge/github.com/txn2/txeh)](https://goreportcard.com/report/github.com/txn2/txeh)
[![Go Reference](https://pkg.go.dev/badge/github.com/txn2/txeh.svg)](https://pkg.go.dev/github.com/txn2/txeh)

Apache 2.0 Licensed. Built by [Craig Johnston](https://github.com/cjimti).

</div>
