# Library Quick Start

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

## Saving

```go
// Save to the configured write path
err := hosts.Save()

// Save to a specific path
err = hosts.SaveAs("/tmp/hosts.test")

// Render without saving
output := hosts.RenderHostsFile()
fmt.Println(output)
```
