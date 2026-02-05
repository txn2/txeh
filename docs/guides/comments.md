# Working with Comments

txeh supports inline comments on host entries, useful for tracking which tool or purpose added specific entries.

## How Comments Work

Comments use the standard hosts file format:

```
IP HOSTNAME [HOSTNAME...] # comment
```

**Key behavior:**

- Comments are used for **grouping**: hosts with the same IP AND same comment are placed on the same line
- Hosts added **without** a comment only match lines **without** a comment
- Hosts added **with** a comment only match lines **with the same comment**
- Comments are **never modified** on existing lines
- If no matching line exists, a new line is created

## CLI Usage

```bash
# Add host with a comment
sudo txeh add 127.0.0.1 myapp --comment "local development"
# Result: 127.0.0.1        myapp # local development

# Add another host with the same comment (groups together)
sudo txeh add 127.0.0.1 myapp2 --comment "local development"
# Result: 127.0.0.1        myapp myapp2 # local development

# Add host with a different comment (new line)
sudo txeh add 127.0.0.1 api --comment "staging environment"
# Result: 127.0.0.1        api # staging environment

# Add host without a comment (only matches lines without comments)
sudo txeh add 127.0.0.1 test
# Goes to first 127.0.0.1 line that has NO comment, or creates new line
```

## Library Usage

```go
// Add hosts with a comment
hosts.AddHostWithComment("127.0.0.1", "myapp", "local development")
hosts.AddHostsWithComment("127.0.0.1", []string{"svc1", "svc2"}, "my project")

// Add hosts without a comment (original behavior)
hosts.AddHost("127.0.0.1", "myhost")
hosts.AddHosts("127.0.0.1", []string{"a", "b", "c"})
```

## Listing by Comment

```bash
# CLI
txeh list bycomment "local development"
```

```go
// Library
entries := hosts.ListHostsByComment("local development")
```

## Removing by Comment

Removes entire lines that have the specified comment.

```bash
# CLI
sudo txeh remove bycomment "local development"
```

```go
// Library
hosts.RemoveByComment("local development")
hosts.RemoveByComments([]string{"old-env", "deprecated"})
```

## Modifying Comments

Comments are never modified on existing lines. To change a comment, remove and re-add:

```bash
sudo txeh remove bycomment "old comment"
sudo txeh add 127.0.0.1 app1 app2 --comment "new comment"
```

## Example Scenario

Starting hosts file:

```
127.0.0.1        localhost
127.0.0.1        app1 app2 # dev services
```

| Command | Result |
|---------|--------|
| `txeh add 127.0.0.1 app3 -c "dev services"` | app3 added to "dev services" line |
| `txeh add 127.0.0.1 api -c "staging env"` | New line: `127.0.0.1 api # staging env` |
| `txeh add 127.0.0.1 test` | test added to localhost line (no comment) |
