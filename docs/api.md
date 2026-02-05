# API Reference

Full API documentation is available on [pkg.go.dev](https://pkg.go.dev/github.com/txn2/txeh).

## Types

### HostsConfig

Configuration for creating a new Hosts instance.

```go
type HostsConfig struct {
    ReadFilePath    string
    WriteFilePath   string
    RawText         *string
    MaxHostsPerLine int
}
```

| Field | Description |
|-------|-------------|
| `ReadFilePath` | Path to read the hosts file from |
| `WriteFilePath` | Path to write the hosts file to |
| `RawText` | Parse from string instead of file (disables Save) |
| `MaxHostsPerLine` | Max hostnames per line (0=auto, -1=unlimited) |

### Hosts

The main struct for hosts file operations. All methods are thread-safe.

```go
type Hosts struct {
    sync.Mutex
    Path          string
    ReadFilePath  string
    WriteFilePath string
    hostFileLines HostFileLines
}
```

### HostFileLine

Represents a single line in the hosts file.

```go
type HostFileLine struct {
    OriginalLineNum int
    LineType        int
    Address         string
    Parts           []string
    Hostnames       []string
    Raw             string
    Trimmed         string
    Comment         string
}
```

Line types:

| Type | Value | Description |
|------|-------|-------------|
| `UNKNOWN` | 0 | Unrecognized line |
| `EMPTY` | 10 | Empty/whitespace line |
| `COMMENT` | 20 | Comment line (`# ...`) |
| `ADDRESS` | 30 | Host entry (`IP hostname [hostname...]`) |

### IPFamily

```go
type IPFamily int64

const (
    IPFamilyV4 IPFamily = iota // 0
    IPFamilyV6                 // 1
)
```

### HostFileLines

```go
type HostFileLines []HostFileLine
```

## Constructors

| Function | Description |
|----------|-------------|
| `NewHostsDefault() (*Hosts, error)` | Load from system default hosts file |
| `NewHosts(config *HostsConfig) (*Hosts, error)` | Load with custom configuration |

## Methods

### Add

| Method | Description |
|--------|-------------|
| `AddHost(ip, hostname)` | Add a single hostname to an IP |
| `AddHosts(ip, hostnames)` | Add multiple hostnames to an IP |
| `AddHostWithComment(ip, hostname, comment)` | Add hostname with inline comment |
| `AddHostsWithComment(ip, hostnames, comment)` | Add hostnames with inline comment |

### Remove

| Method | Description |
|--------|-------------|
| `RemoveHost(hostname)` | Remove a hostname |
| `RemoveHosts(hostnames)` | Remove multiple hostnames |
| `RemoveAddress(ip)` | Remove an IP and all its hostnames |
| `RemoveAddresses(ips)` | Remove multiple IPs |
| `RemoveCIDRs(cidrs) error` | Remove all IPs in CIDR ranges |
| `RemoveByComment(comment)` | Remove entries with a comment |
| `RemoveByComments(comments)` | Remove entries matching any comment |

### Query

| Method | Returns | Description |
|--------|---------|-------------|
| `ListHostsByIP(ip)` | `[]string` | Hostnames at an IP |
| `ListAddressesByHost(host, exact)` | `[][]string` | IPs for a hostname |
| `ListHostsByCIDR(cidr)` | `[][]string` | IPs and hostnames in a CIDR range |
| `HostAddressLookup(host, family)` | `(bool, string, int)` | Lookup a hostname |
| `ListHostsByComment(comment)` | `[]string` | Hostnames with a comment |

### Output

| Method | Description |
|--------|-------------|
| `RenderHostsFile() string` | Render the hosts file as a string |
| `GetHostFileLines() HostFileLines` | Get a copy of all parsed host file lines |
| `Save() error` | Save to the configured write path |
| `SaveAs(path) error` | Save to a specific path |
