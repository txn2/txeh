# Getting Started

Get up and running with txeh in under 5 minutes.

## Install the CLI

=== "Homebrew (macOS/Linux)"

    ```bash
    brew install txn2/tap/txeh
    ```

=== "Go Install"

    ```bash
    go install github.com/txn2/txeh/txeh@latest
    ```

=== "Binary Download"

    Download the latest release from [GitHub Releases](https://github.com/txn2/txeh/releases) for your platform.

## Your First Commands

Add a hostname:

```bash
sudo txeh add 127.0.0.1 myapp.local
```

List hosts for an IP:

```bash
txeh list ip 127.0.0.1
```

Remove a hostname:

```bash
sudo txeh remove host myapp.local
```

Preview changes without saving (dry run):

```bash
sudo txeh add 127.0.0.1 myapp.local --dryrun
```

## Use as a Go Library

Add the dependency:

```bash
go get github.com/txn2/txeh
```

Basic usage:

```go
package main

import (
    "fmt"
    "github.com/txn2/txeh"
)

func main() {
    hosts, err := txeh.NewHostsDefault()
    if err != nil {
        panic(err)
    }

    hosts.AddHost("127.0.0.1", "myapp.local")
    fmt.Println(hosts.RenderHostsFile())

    err = hosts.Save()
    if err != nil {
        panic(err)
    }
}
```

## Next Steps

- [CLI Commands Reference](cli/commands.md) -- Full CLI documentation
- [Go Library Quick Start](library/quickstart.md) -- Detailed library usage
- [Working with Comments](guides/comments.md) -- Organize entries with inline comments
