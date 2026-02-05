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

### HostFileLine

Represents a single line in the hosts file.

```go
type HostFileLine struct {
    OriginalLineNum int
    LineType         int
    Address          string
    Parts            []string
    Hostnames        []string
    Raw              string
    Trimed           string
    Comment          string
}
```

### IPFamily

```go
type IPFamily int

const (
    IPFamilyV4 IPFamily = 4
    IPFamilyV6 IPFamily = 6
)
```

## Constructors

| Function | Description |
|----------|-------------|
| `NewHostsDefault()` | Load from system default hosts file |
| `NewHosts(config)` | Load with custom configuration |

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
| `ListHostsByCIDR(cidr)` | `[]string` | Hostnames in a CIDR range |
| `HostAddressLookup(host, family)` | `bool, string, IPFamily` | Lookup a hostname |
| `ListHostsByComment(comment)` | `[]string` | Hostnames with a comment |

### Output

| Method | Description |
|--------|-------------|
| `RenderHostsFile() string` | Render the hosts file as a string |
| `GetHostFileLines() *HostFileLines` | Get parsed line data |
| `Save() error` | Save to the configured write path |
| `SaveAs(path) error` | Save to a specific path |
