# CLI Reference

txeh provides a command-line interface for managing `/etc/hosts` file entries.

## Installation

=== "Homebrew (macOS/Linux)"

    ```bash
    brew install txn2/tap/txeh
    ```

=== "Go Install"

    Requires Go 1.24 or later.

    ```bash
    go install github.com/txn2/txeh/txeh@latest
    ```

    The binary is installed to `$GOPATH/bin`. Ensure this directory is in your `$PATH`.

=== "Binary Download"

    Pre-built binaries are available on the [GitHub Releases](https://github.com/txn2/txeh/releases) page:

    | OS | Architectures |
    |----|---------------|
    | Linux | amd64, arm64, arm, 386 |
    | macOS | amd64 (Intel), arm64 (Apple Silicon) |
    | Windows | amd64, 386 |

    `.deb`, `.rpm`, and `.apk` packages are also available for each release.

=== "Build from Source"

    ```bash
    git clone https://github.com/txn2/txeh.git
    cd txeh
    go install ./txeh
    ```

Verify installation:

```bash
txeh version
```

## Quick Start

```bash
# Add a hostname
sudo txeh add 127.0.0.1 myapp.local

# List hosts for an IP
txeh list ip 127.0.0.1

# Remove a hostname
sudo txeh remove host myapp.local

# Preview changes without saving
sudo txeh add 127.0.0.1 myapp.local --dryrun
```

## Global Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--dryrun` | `-d` | Output to stdout without saving |
| `--quiet` | `-q` | Suppress output |
| `--read` | `-r` | Override path to read hosts file |
| `--write` | `-w` | Override path to write hosts file |
| `--max-hosts-per-line` | `-m` | Max hostnames per line (0=auto, -1=unlimited) |

## Commands

### add

Add one or more hostnames to an IP address.

```bash
sudo txeh add [IP] [HOSTNAME] [HOSTNAME]...
```

**Flags:**

| Flag | Short | Description |
|------|-------|-------------|
| `--comment` | `-c` | Add an inline comment to the entry |

**Examples:**

```bash
# Add a single hostname
sudo txeh add 127.0.0.1 myapp.local

# Add multiple hostnames
sudo txeh add 127.0.0.1 app1.local app2.local app3.local

# Add with a comment for organization
sudo txeh add 127.0.0.1 myapp.local --comment "dev environment"

# Preview without saving
sudo txeh add 127.0.0.1 myapp.local --dryrun
```

### remove host

Remove one or more hostnames from the hosts file.

```bash
sudo txeh remove host [HOSTNAME] [HOSTNAME]...
```

```bash
sudo txeh remove host myapp.local
sudo txeh remove host app1.local app2.local
```

### remove ip

Remove an IP address and all hostnames associated with it.

```bash
sudo txeh remove ip [IP] [IP]...
```

```bash
sudo txeh remove ip 127.0.0.1
sudo txeh remove ip 10.0.0.1 10.0.0.2
```

### remove cidr

Remove all addresses within one or more CIDR ranges.

```bash
sudo txeh remove cidr [CIDR] [CIDR]...
```

```bash
sudo txeh remove cidr 10.0.0.0/24
sudo txeh remove cidr 192.168.1.0/24 172.16.0.0/12

# Preview with dry run
sudo txeh remove cidr 10.0.0.0/8 --dryrun
```

### remove bycomment

Remove all host entries that have a specific inline comment.

```bash
sudo txeh remove bycomment [COMMENT]
```

```bash
sudo txeh remove bycomment "dev environment"
sudo txeh remove bycomment "kubefwd"
```

### list ip

List hostnames associated with one or more IP addresses.

```bash
txeh list ip [IP] [IP]...
```

### list cidr

List hostnames for all addresses within CIDR ranges.

```bash
txeh list cidr [CIDR] [CIDR]...
```

### list host

List IP addresses associated with one or more hostnames.

```bash
txeh list host [HOSTNAME] [HOSTNAME]...
```

### list bycomment

List all hosts that have a specific inline comment.

```bash
txeh list bycomment [COMMENT]
```

### show

Display the full rendered hosts file.

```bash
txeh show
```

### version

Print the txeh version.

```bash
txeh version
```
