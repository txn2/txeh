![txeh - /etc/hosts mangement](txeh.png)


# Etc Hosts Management Library

[![Go Report Card](https://goreportcard.com/badge/github.com/txn2/txeh)](https://goreportcard.com/report/github.com/txn2/txeh)
[![GoDoc](https://godoc.org/github.com/txn2/irsync/txeh?status.svg)](https://godoc.org/github.com/txn2/txeh)



This small go library was developed to encapsulate the complexity of
working with /etc/hosts by providing a simple interface to load, create, remove and save entries to an /etc/host file. No validation is done on the input data. Validation is considered out of scope for this project, use with caution.

Basic implemention:
```go

package main

import (
	"fmt"
	"strings"

	"github.com/txn2/txeh"
)

func main() {
	hosts, err := txeh.NewHostsDefault()
	if err != nil {
		panic(err)
	}

	hosts.AddHost("127.100.100.100", "test")
	hosts.AddHost("127.100.100.101", "logstash")
	hosts.AddHosts("127.100.100.102", []string{"a", "b", "c"})
	
	hosts.RemoveHosts([]string{"example", "example.machine", "example.machine.example.com"})
	hosts.RemoveHosts(strings.Fields("example2 example.machine2 example.machine.example.com2"))

	
	hosts.RemoveAddress("127.1.27.1")
	
	removeList := []string{
		"127.1.27.15",
		"127.1.27.14",
		"127.1.27.13",
	}
	
	hosts.RemoveAddresses(removeList)
	
	hfData := hosts.RenderHostsFile()

	fmt.Println(hfData)
	
	hosts.SaveAs("./test.hosts")
}

```
