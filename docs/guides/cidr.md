# CIDR Operations

txeh supports CIDR (Classless Inter-Domain Routing) notation for bulk operations on IP address ranges.

## Removing by CIDR

Remove all host entries whose IP addresses fall within one or more CIDR ranges.

### CLI

```bash
# Remove all hosts in a /24 subnet
sudo txeh remove cidr 10.0.0.0/24

# Remove multiple CIDR ranges
sudo txeh remove cidr 192.168.1.0/24 172.16.0.0/12

# Preview with dry run
sudo txeh remove cidr 10.0.0.0/8 --dryrun
```

### Library

```go
err := hosts.RemoveCIDRs([]string{"10.0.0.0/24", "172.16.0.0/12"})
if err != nil {
    log.Fatal(err)
}
```

!!! note
    `RemoveCIDRs` is the only removal method that returns an error, since CIDR parsing can fail with invalid input.

## Listing by CIDR

List all hostnames whose IP addresses fall within a CIDR range.

### CLI

```bash
txeh list cidr 127.0.0.0/8
txeh list cidr 10.0.0.0/24 192.168.0.0/16
```

### Library

```go
hostnames := hosts.ListHostsByCIDR("10.0.0.0/24")
for _, hn := range hostnames {
    fmt.Println(hn)
}
```

## Use Cases

**Clean up development entries:**

```bash
# Remove all loopback entries added by kubefwd
sudo txeh remove cidr 127.1.0.0/16
```

**Audit network ranges:**

```bash
# List all hosts in your private network range
txeh list cidr 10.0.0.0/8
```

**DevOps automation:**

```go
// Remove stale entries from a previous deployment
cidrs := []string{"10.244.0.0/16", "10.96.0.0/12"}
if err := hosts.RemoveCIDRs(cidrs); err != nil {
    log.Fatal(err)
}
hosts.Save()
```
