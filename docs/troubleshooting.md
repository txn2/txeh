# Troubleshooting

## Common Issues

### Permission denied

```
Error: could not save /etc/hosts. Reason: open /etc/hosts: permission denied
```

The hosts file requires root/administrator privileges to modify. Use `sudo`:

```bash
sudo txeh add 127.0.0.1 myapp.local
```

### Command not found after go install

```
txeh: command not found
```

The Go binary directory may not be in your `$PATH`. Add it:

```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

Add this line to your shell profile (`~/.bashrc`, `~/.zshrc`, etc.) to make it permanent.

### Hosts not resolving

After adding entries, if the hostname doesn't resolve:

1. **Flush DNS cache:**

    === "macOS"

        ```bash
        sudo dscacheutil -flushcache && sudo killall -HUP mDNSResponder
        ```

    === "Linux (systemd)"

        ```bash
        sudo systemd-resolve --flush-caches
        ```

    === "Windows"

        ```cmd
        ipconfig /flushdns
        ```

2. **Verify the entry was added:**

    ```bash
    txeh show
    ```

3. **Check for conflicting entries.** Multiple entries for the same hostname may cause unexpected behavior.

### Windows hosts file location

txeh auto-detects the Windows hosts file via the `SystemRoot` environment variable. The default is:

```
C:\Windows\System32\drivers\etc\hosts
```

If your Windows installation uses a non-standard path, set `SystemRoot` or use the `--read`/`--write` flags.

## Getting Help

- [GitHub Issues](https://github.com/txn2/txeh/issues): Bug reports and feature requests.
- [GitHub Discussions](https://github.com/txn2/txeh/discussions): Questions and community support.
