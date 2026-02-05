# Configuration

## HostsConfig

The `HostsConfig` struct controls how txeh reads and writes the hosts file.

```go
type HostsConfig struct {
    ReadFilePath    string
    WriteFilePath   string
    RawText         *string
    MaxHostsPerLine int
}
```

### ReadFilePath / WriteFilePath

Override the default system hosts file location. If only `ReadFilePath` is set, `WriteFilePath` defaults to the same path.

```go
hosts, err := txeh.NewHosts(&txeh.HostsConfig{
    ReadFilePath:  "/etc/hosts",
    WriteFilePath: "/tmp/hosts.modified",
})
```

### RawText

Parse hosts data from a string instead of a file. When `RawText` is set, `Save()` will return an error since there is no file path. Use `RenderHostsFile()` or `SaveAs()` instead.

```go
data := "127.0.0.1 localhost\n"
hosts, err := txeh.NewHosts(&txeh.HostsConfig{
    RawText: &data,
})
```

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

## CLI Flags

The CLI exposes configuration through global flags:

| Flag | Config Equivalent |
|------|-------------------|
| `--read` / `-r` | `ReadFilePath` |
| `--write` / `-w` | `WriteFilePath` |
| `--max-hosts-per-line` / `-m` | `MaxHostsPerLine` |
| `--dryrun` / `-d` | Uses `RenderHostsFile()` instead of `Save()` |
| `--quiet` / `-q` | Suppresses informational output |
