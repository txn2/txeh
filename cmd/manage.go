package main

import (
	"fmt"

	"github.com/txn2/txeh"
)

func main() {

	hosts := txeh.NewHosts(txeh.HostsConfig{
		FilePath: "/etc/hosts",
	})

	hosts.AddHost("127.1.27.24", "api-asset")
	hosts.AddHost("127.1.27.24", "api-asset.dcp")
	hosts.AddHost("127.1.27.24", "api-auth.dcp")
	hosts.AddHost("127.1.27.24", "google4.com")
	hosts.AddHost("127.1.27.24", "api-auth.dcp.svc.cluster.local")
	hosts.AddHost("127.1.27.201", "api-asset.dcp.svc.cluster.local")

	hosts.AddHost("127.1.27.25", "api-auth")
	hosts.AddHost("127.1.27.200", "api-tcpql")

	hosts.AddHost("127.1.27.3", "craig")
	hosts.AddHost("127.1.27.3", "craig2")
	hosts.AddHost("127.2.27.3", "craig3")

	hf := hosts.RenderHostsFile()

	fmt.Print(hf)
}
