# Commands Reference

## add

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

## remove host

Remove one or more hostnames from the hosts file.

```bash
sudo txeh remove host [HOSTNAME] [HOSTNAME]...
```

**Examples:**

```bash
sudo txeh remove host myapp.local
sudo txeh remove host app1.local app2.local
```

## remove ip

Remove an IP address and all hostnames associated with it.

```bash
sudo txeh remove ip [IP] [IP]...
```

**Examples:**

```bash
sudo txeh remove ip 127.0.0.1
sudo txeh remove ip 10.0.0.1 10.0.0.2
```

## remove cidr

Remove all addresses within one or more CIDR ranges.

```bash
sudo txeh remove cidr [CIDR] [CIDR]...
```

**Examples:**

```bash
sudo txeh remove cidr 10.0.0.0/24
sudo txeh remove cidr 192.168.1.0/24 172.16.0.0/12
```

## remove bycomment

Remove all host entries that have a specific inline comment.

```bash
sudo txeh remove bycomment [COMMENT]
```

**Examples:**

```bash
sudo txeh remove bycomment "dev environment"
sudo txeh remove bycomment "kubefwd"
```

## list ip

List hostnames associated with one or more IP addresses.

```bash
txeh list ip [IP] [IP]...
```

## list cidr

List hostnames for all addresses within CIDR ranges.

```bash
txeh list cidr [CIDR] [CIDR]...
```

## list host

List IP addresses associated with one or more hostnames.

```bash
txeh list host [HOSTNAME] [HOSTNAME]...
```

## list bycomment

List all hosts that have a specific inline comment.

```bash
txeh list bycomment [COMMENT]
```

## show

Display the full rendered hosts file.

```bash
txeh show
```

## version

Print the txeh version.

```bash
txeh version
```
