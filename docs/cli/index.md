# CLI Overview

txeh provides a command-line interface for managing `/etc/hosts` file entries.

```
 _            _
| |___  _____| |__
| __\ \/ / _ \ '_ \
| |_ >  <  __/ | | |
 \__/_/\_\___|_| |_|
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

| Command | Description |
|---------|-------------|
| `add` | Add hostnames to an IP address |
| `remove host` | Remove hostnames |
| `remove ip` | Remove IP addresses and their hosts |
| `remove cidr` | Remove CIDR ranges |
| `remove bycomment` | Remove entries by comment |
| `list ip` | List hosts for IP addresses |
| `list cidr` | List hosts in CIDR ranges |
| `list host` | List IPs for hostnames |
| `list bycomment` | List hosts by comment |
| `show` | Show the full hosts file |
| `version` | Print version |

See [Commands](commands.md) for detailed usage of each command.
