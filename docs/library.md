# Go Library

txeh provides a Go library for programmatic `/etc/hosts` file management. It was originally built to support [kubefwd](https://github.com/txn2/kubefwd) for Kubernetes service forwarding.

## Installation

```bash
go get github.com/txn2/txeh
```

## Loading the Hosts File

### System default

```go
hosts, err := txeh.NewHostsDefault()
if err != nil {
    log.Fatal(err)
}
```

### Custom paths

```go
hosts, err := txeh.NewHosts(&txeh.HostsConfig{
    ReadFilePath:  "/custom/path/hosts",
    WriteFilePath: "/custom/path/hosts",
})
```

### From a string (in-memory)

```go
raw := "127.0.0.1 localhost\n::1 localhost\n"
hosts, err := txeh.NewHosts(&txeh.HostsConfig{
    RawText: &raw,
})
```

When `RawText` is set, `Save()` returns an error since there is no file path. Use `RenderHostsFile()` or `SaveAs()` instead.

## Adding Hosts

```go
// Single hostname
hosts.AddHost("127.0.0.1", "myapp.local")

// Multiple hostnames to one IP
hosts.AddHosts("127.0.0.1", []string{"app1", "app2", "app3"})

// With an inline comment
hosts.AddHostWithComment("127.0.0.1", "myapp", "dev environment")
hosts.AddHostsWithComment("127.0.0.1", []string{"svc1", "svc2"}, "kubefwd")
```

## Removing Hosts

```go
// By hostname
hosts.RemoveHost("myapp.local")
hosts.RemoveHosts([]string{"app1", "app2"})

// By IP address
hosts.RemoveAddress("127.0.0.1")
hosts.RemoveAddresses([]string{"10.0.0.1", "10.0.0.2"})

// By CIDR range
err := hosts.RemoveCIDRs([]string{"10.0.0.0/24"})

// By comment
hosts.RemoveByComment("dev environment")
hosts.RemoveByComments([]string{"old-env", "deprecated"})
```

## Querying

```go
// List hostnames at an IP
hostnames := hosts.ListHostsByIP("127.0.0.1")

// List IPs for a hostname (exact or partial match)
addresses := hosts.ListAddressesByHost("myapp", true)   // exact
addresses = hosts.ListAddressesByHost("myapp", false)    // contains

// List hostnames in a CIDR range
hostnames = hosts.ListHostsByCIDR("10.0.0.0/24")

// Lookup a specific hostname
found, ip, ipFamily := hosts.HostAddressLookup("myapp", txeh.IPFamilyV4)

// List by comment
entries := hosts.ListHostsByComment("kubefwd")
```

## Saving and Rendering

```go
// Save to the configured write path
err := hosts.Save()

// Save to a specific path
err = hosts.SaveAs("/tmp/hosts.test")

// Render without saving
output := hosts.RenderHostsFile()
fmt.Println(output)
```

## CIDR Operations

txeh supports CIDR (Classless Inter-Domain Routing) notation for bulk operations on IP address ranges.

### Removing by CIDR

Remove all host entries whose IP addresses fall within one or more CIDR ranges.

```go
err := hosts.RemoveCIDRs([]string{"10.0.0.0/24", "172.16.0.0/12"})
if err != nil {
    log.Fatal(err)
}
```

!!! note
    `RemoveCIDRs` is the only removal method that returns an error, since CIDR parsing can fail with invalid input.

### Listing by CIDR

```go
hostnames := hosts.ListHostsByCIDR("10.0.0.0/24")
for _, hn := range hostnames {
    fmt.Println(hn)
}
```

### Example: cleaning up development entries

```go
// Remove stale entries from a previous deployment
cidrs := []string{"10.244.0.0/16", "10.96.0.0/12"}
if err := hosts.RemoveCIDRs(cidrs); err != nil {
    log.Fatal(err)
}
hosts.Save()
```

## Inline Comments

txeh supports inline comments on host entries. Comments are useful for tracking which tool or purpose added specific entries.

Comments use the standard hosts file format: `IP HOSTNAME [HOSTNAME...] # comment`

**Key behavior:**

- Comments are used for grouping: hosts with the same IP and same comment are placed on the same line.
- Hosts added without a comment only match lines without a comment.
- Hosts added with a comment only match lines with the same comment.
- Comments are never modified on existing lines.
- If no matching line exists, a new line is created.

### Adding with comments

```go
hosts.AddHostWithComment("127.0.0.1", "myapp", "local development")
hosts.AddHostsWithComment("127.0.0.1", []string{"svc1", "svc2"}, "my project")
```

### Listing and removing by comment

```go
// List by comment
entries := hosts.ListHostsByComment("local development")

// Remove by comment
hosts.RemoveByComment("local development")
hosts.RemoveByComments([]string{"old-env", "deprecated"})
```

### Modifying comments

Comments are never modified on existing lines. To change a comment, remove and re-add:

```go
hosts.RemoveByComment("old comment")
hosts.AddHostsWithComment("127.0.0.1", []string{"app1", "app2"}, "new comment")
```

## Configuration

### MaxHostsPerLine

Controls how many hostnames can appear on a single line.

| Value | Behavior |
|-------|----------|
| `0` | Auto (9 on Windows, unlimited elsewhere) |
| `-1` | Unlimited |
| `>0` | Explicit limit |

When the limit is reached, additional hostnames are placed on a new line with the same IP address.

```go
hosts, err := txeh.NewHosts(&txeh.HostsConfig{
    MaxHostsPerLine: 5,
})
```

### ReadFilePath and WriteFilePath

Override the default system hosts file location. If only `ReadFilePath` is set, `WriteFilePath` defaults to the same path.

```go
hosts, err := txeh.NewHosts(&txeh.HostsConfig{
    ReadFilePath:  "/etc/hosts",
    WriteFilePath: "/tmp/hosts.modified",
})
```

## Thread Safety

All public methods on `Hosts` acquire a mutex before reading or modifying the internal state. This makes txeh safe for concurrent use from multiple goroutines.
