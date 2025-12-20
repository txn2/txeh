package txeh

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"
	"testing"
)

var mockHostsData = `127.0.0.1        localhost
127.0.1.1        straylight-desk
127.0.1.1        bad-idea
# /etc/hosts is not DNS and one should not expect consistent behaviour when assigning a host to multiple
# IP addresses. TODO - Support for duplicate hostname for IPv4 and IPv6
# see - "man 5 hosts"
# IP_address canonical_hostname [aliases...]


# The following lines are desirable for IPv6 capable hosts
::1              localhost ip6-localhost ip6-loopback
fe00::0          ip6-localnet
ff00::0          ip6-mcastprefix
ff02::1          ip6-allnodes
ff02::2          ip6-allrouters
127.1.27.1       nifi-mtx-2 nifi-mtx-2.d4l-demo nifi-mtx-2.d4l-demo.svc nifi-mtx-2.d4l-demo.svc.cluster.local nifi-mtx-2.d4l-demo.dwdev4 nifi-mtx-2.d4l-demo.svc.dwdev4 nifi-mtx-2.d4l-demo.svc.cluster.dwdev4
127.1.27.2       demo-log-search-api demo-log-search-api.d4l-demo demo-log-search-api.d4l-demo.svc demo-log-search-api.d4l-demo.svc.cluster.local demo-log-search-api.d4l-demo.dwdev4 demo-log-search-api.d4l-demo.svc.dwdev4 demo-log-search-api.d4l-demo.svc.cluster.dwdev4
127.1.27.3       warehouse-cluster-repl warehouse-cluster-repl.d4l-demo warehouse-cluster-repl.d4l-demo.svc warehouse-cluster-repl.d4l-demo.svc.cluster.local warehouse-cluster-repl.d4l-demo.dwdev4 warehouse-cluster-repl.d4l-demo.svc.dwdev4 warehouse-cluster-repl.d4l-demo.svc.cluster.dwdev4


127.1.27.4       superset-redis-headless superset-redis-headless.d4l-demo superset-redis-headless.d4l-demo.svc superset-redis-headless.d4l-demo.svc.cluster.local superset-redis-headless.d4l-demo.dwdev4 superset-redis-headless.d4l-demo.svc.dwdev4 superset-redis-headless.d4l-demo.svc.cluster.dwdev4
127.1.27.5       superset-redis-master-0.superset-redis-headless superset-redis-master-0.superset-redis-headless.d4l-demo superset-redis-master-0.superset-redis-headless.d4l-demo.svc superset-redis-master-0.superset-redis-headless.d4l-demo.svc.cluster.local superset-redis-master-0.superset-redis-headless.d4l-demo.dwdev4 superset-redis-master-0.superset-redis-headless.d4l-demo.svc.dwdev4 superset-redis-master-0.superset-redis-headless.d4l-demo.svc.cluster.dwdev4
127.1.27.6       demo-log-hl-svc demo-log-hl-svc.d4l-demo demo-log-hl-svc.d4l-demo.svc demo-log-hl-svc.d4l-demo.svc.cluster.local demo-log-hl-svc.d4l-demo.dwdev4 demo-log-hl-svc.d4l-demo.svc.dwdev4 demo-log-hl-svc.d4l-demo.svc.cluster.dwdev4
127.1.27.7       demo-log-0.demo-log-hl-svc demo-log-0.demo-log-hl-svc.d4l-demo demo-log-0.demo-log-hl-svc.d4l-demo.svc demo-log-0.demo-log-hl-svc.d4l-demo.svc.cluster.local demo-log-0.demo-log-hl-svc.d4l-demo.dwdev4 demo-log-0.demo-log-hl-svc.d4l-demo.svc.dwdev4 demo-log-0.demo-log-hl-svc.d4l-demo.svc.cluster.dwdev4
127.1.27.8       nifi-ca nifi-ca.d4l-demo nifi-ca.d4l-demo.svc nifi-ca.d4l-demo.svc.cluster.local nifi-ca.d4l-demo.dwdev4 nifi-ca.d4l-demo.svc.dwdev4 nifi-ca.d4l-demo.svc.cluster.dwdev4
127.1.27.9       proxy-api proxy-api.d4l-demo proxy-api.d4l-demo.svc proxy-api.d4l-demo.svc.cluster.local proxy-api.d4l-demo.dwdev4 proxy-api.d4l-demo.svc.dwdev4 proxy-api.d4l-demo.svc.cluster.dwdev4
127.1.27.10      demo-hl demo-hl.d4l-demo demo-hl.d4l-demo.svc demo-hl.d4l-demo.svc.cluster.local demo-hl.d4l-demo.dwdev4 demo-hl.d4l-demo.svc.dwdev4 demo-hl.d4l-demo.svc.cluster.dwdev4
127.1.27.11      demo-pool-0-0.demo-hl demo-pool-0-0.demo-hl.d4l-demo demo-pool-0-0.demo-hl.d4l-demo.svc demo-pool-0-0.demo-hl.d4l-demo.svc.cluster.local demo-pool-0-0.demo-hl.d4l-demo.dwdev4 demo-pool-0-0.demo-hl.d4l-demo.svc.dwdev4 demo-pool-0-0.demo-hl.d4l-demo.svc.cluster.dwdev4
127.1.27.12      demo-pool-0-1.demo-hl demo-pool-0-1.demo-hl.d4l-demo demo-pool-0-1.demo-hl.d4l-demo.svc demo-pool-0-1.demo-hl.d4l-demo.svc.cluster.local demo-pool-0-1.demo-hl.d4l-demo.dwdev4 demo-pool-0-1.demo-hl.d4l-demo.svc.dwdev4 demo-pool-0-1.demo-hl.d4l-demo.svc.cluster.dwdev4
127.1.27.13      demo-pool-0-2.demo-hl demo-pool-0-2.demo-hl.d4l-demo demo-pool-0-2.demo-hl.d4l-demo.svc demo-pool-0-2.demo-hl.d4l-demo.svc.cluster.local demo-pool-0-2.demo-hl.d4l-demo.dwdev4 demo-pool-0-2.demo-hl.d4l-demo.svc.dwdev4 demo-pool-0-2.demo-hl.d4l-demo.svc.cluster.dwdev4
127.1.27.14      proxy-public proxy-public.d4l-demo proxy-public.d4l-demo.svc proxy-public.d4l-demo.svc.cluster.local proxy-public.d4l-demo.dwdev4 proxy-public.d4l-demo.svc.dwdev4 proxy-public.d4l-demo.svc.cluster.dwdev4
127.1.27.15      nifi-port-proxy nifi-port-proxy.d4l-demo nifi-port-proxy.d4l-demo.svc nifi-port-proxy.d4l-demo.svc.cluster.local nifi-port-proxy.d4l-demo.dwdev4 nifi-port-proxy.d4l-demo.svc.dwdev4 nifi-port-proxy.d4l-demo.svc.cluster.dwdev4
127.1.27.16      nifi nifi.d4l-demo nifi.d4l-demo.svc nifi.d4l-demo.svc.cluster.local nifi.d4l-demo.dwdev4 nifi.d4l-demo.svc.dwdev4 nifi.d4l-demo.svc.cluster.dwdev4
127.1.27.19      nifi-2.nifi nifi-2.nifi.d4l-demo nifi-2.nifi.d4l-demo.svc nifi-2.nifi.d4l-demo.svc.cluster.local nifi-2.nifi.d4l-demo.dwdev4 nifi-2.nifi.d4l-demo.svc.dwdev4 nifi-2.nifi.d4l-demo.svc.cluster.dwdev4
127.1.27.17      nifi-0.nifi nifi-0.nifi.d4l-demo nifi-0.nifi.d4l-demo.svc nifi-0.nifi.d4l-demo.svc.cluster.local nifi-0.nifi.d4l-demo.dwdev4 nifi-0.nifi.d4l-demo.svc.dwdev4 nifi-0.nifi.d4l-demo.svc.cluster.dwdev4
127.1.27.18      nifi-1.nifi nifi-1.nifi.d4l-demo nifi-1.nifi.d4l-demo.svc nifi-1.nifi.d4l-demo.svc.cluster.local nifi-1.nifi.d4l-demo.dwdev4 nifi-1.nifi.d4l-demo.svc.dwdev4 nifi-1.nifi.d4l-demo.svc.cluster.dwdev4
127.1.27.20      demo-prometheus-hl-svc demo-prometheus-hl-svc.d4l-demo demo-prometheus-hl-svc.d4l-demo.svc demo-prometheus-hl-svc.d4l-demo.svc.cluster.local demo-prometheus-hl-svc.d4l-demo.dwdev4 demo-prometheus-hl-svc.d4l-demo.svc.dwdev4 demo-prometheus-hl-svc.d4l-demo.svc.cluster.dwdev4
127.1.27.21      demo-prometheus-0.demo-prometheus-hl-svc demo-prometheus-0.demo-prometheus-hl-svc.d4l-demo demo-prometheus-0.demo-prometheus-hl-svc.d4l-demo.svc demo-prometheus-0.demo-prometheus-hl-svc.d4l-demo.svc.cluster.local demo-prometheus-0.demo-prometheus-hl-svc.d4l-demo.dwdev4 demo-prometheus-0.demo-prometheus-hl-svc.d4l-demo.svc.dwdev4 demo-prometheus-0.demo-prometheus-hl-svc.d4l-demo.svc.cluster.dwdev4
# existing comment
127.1.27.22      zk-cs zk-cs.d4l-demo zk-cs.d4l-demo.svc zk-cs.d4l-demo.svc.cluster.local zk-cs.d4l-demo.dwdev4 zk-cs.d4l-demo.svc.dwdev4 zk-cs.d4l-demo.svc.cluster.dwdev4
127.1.27.23      nifi-mtx-1 nifi-mtx-1.d4l-demo nifi-mtx-1.d4l-demo.svc nifi-mtx-1.d4l-demo.svc.cluster.local nifi-mtx-1.d4l-demo.dwdev4 nifi-mtx-1.d4l-demo.svc.dwdev4 nifi-mtx-1.d4l-demo.svc.cluster.dwdev4
127.1.27.24      pga;dmin pgadmin.d4l-demo pgadmin.d4l-demo.svc pgadmin.d4l-demo.svc.cluster.local pgadmin.d4l-demo.dwdev4 pgadmin.d4l-demo.svc.dwdev4 pgadmin.d4l-demo.svc.cluster.dwdev4 #another comment
127.1.27.25      superset-redis-master superset-redis-master.d4l-demo superset-redis-master.d4l-demo.svc superset-redis-master.d4l-demo.svc.cluster.local superset-redis-master.d4l-demo.dwdev4 superset-redis-master.d4l-demo.svc.dwdev4 superset-redis-master.d4l-demo.svc.cluster.dwdev4
127.1.27.27      keycloak keycloak.d4l-demo keycloak.d4l-demo.svc keycloak.d4l-demo.svc.cluster.local keycloak.d4l-demo.dwdev4 keycloak.d4l-demo.svc.dwdev4 keycloak.d4l-demo.svc.cluster.dwdev4
127.1.27.28      hub hub.d4l-demo hub.d4l-demo.svc hub.d4l-demo.svc.cluster.local hub.d4l-demo.dwdev4 hub.d4l-demo.svc.dwdev4 hub.d4l-demo.svc.cluster.dwdev4
127.1.27.29      nifi-registry nifi-registry.d4l-demo nifi-registry.d4l-demo.svc nifi-registry.d4l-demo.svc.cluster.local nifi-registry.d4l-demo.dwdev4 nifi-registry.d4l-demo.svc.dwdev4 nifi-registry.d4l-demo.svc.cluster.dwdev4
127.1.27.30      demo-console demo-console.d4l-demo demo-console.d4l-demo.svc demo-console.d4l-demo.svc.cluster.local demo-console.d4l-demo.dwdev4 demo-console.d4l-demo.svc.dwdev4 demo-console.d4l-demo.svc.cluster.dwdev4
127.1.27.31      superset superset.d4l-demo superset.d4l-demo.svc superset.d4l-demo.svc.cluster.local superset.d4l-demo.dwdev4 superset.d4l-demo.svc.dwdev4 superset.d4l-demo.svc.cluster.dwdev4
127.1.27.32      zk-hs zk-hs.d4l-demo zk-hs.d4l-demo.svc zk-hs.d4l-demo.svc.cluster.local zk-hs.d4l-demo.dwdev4 zk-hs.d4l-demo.svc.dwdev4 zk-hs.d4l-demo.svc.cluster.dwdev4
127.1.27.35      zk-2.zk-hs zk-2.zk-hs.d4l-demo zk-2.zk-hs.d4l-demo.svc zk-2.zk-hs.d4l-demo.svc.cluster.local zk-2.zk-hs.d4l-demo.dwdev4 zk-2.zk-hs.d4l-demo.svc.dwdev4 zk-2.zk-hs.d4l-demo.svc.cluster.dwdev4
127.1.27.33      zk-0.zk-hs zk-0.zk-hs.d4l-demo zk-0.zk-hs.d4l-demo.svc zk-0.zk-hs.d4l-demo.svc.cluster.local zk-0.zk-hs.d4l-demo.dwdev4 zk-0.zk-hs.d4l-demo.svc.dwdev4 zk-0.zk-hs.d4l-demo.svc.cluster.dwdev4
127.1.27.34      zk-1.zk-hs zk-1.zk-hs.d4l-demo zk-1.zk-hs.d4l-demo.svc zk-1.zk-hs.d4l-demo.svc.cluster.local zk-1.zk-hs.d4l-demo.dwdev4 zk-1.zk-hs.d4l-demo.svc.dwdev4 zk-1.zk-hs.d4l-demo.svc.cluster.dwdev4
127.1.27.36      nifi-mtx-0 nifi-mtx-0.d4l-demo nifi-mtx-0.d4l-demo.svc nifi-mtx-0.d4l-demo.svc.cluster.local nifi-mtx-0.d4l-demo.dwdev4 nifi-mtx-0.d4l-demo.svc.dwdev4 nifi-mtx-0.d4l-demo.svc.cluster.dwdev4
127.1.27.37      opsdb-cluster-repl opsdb-cluster-repl.d4l-demo opsdb-cluster-repl.d4l-demo.svc opsdb-cluster-repl.d4l-demo.svc.cluster.local opsdb-cluster-repl.d4l-demo.dwdev4 opsdb-cluster-repl.d4l-demo.svc.dwdev4 opsdb-cluster-repl.d4l-demo.svc.cluster.dwdev4
`

func TestNewHostsDefault(t *testing.T) {
	hosts, err := NewHostsDefault()
	if err != nil {
		panic(err)
	}

	testHostName := "test-new-hosts-default"
	hosts.AddHost("127.100.100.100", testHostName)
	foundHostnames := hosts.ListHostsByIP("127.100.100.100")
	found := false
	for _, hn := range foundHostnames {
		if hn == testHostName {
			found = true
			break
		}
	}

	if !found {
		t.Fatal("TestNewHostsDefault was unable find the hostname it added")
	}
}

func TestNewHosts(t *testing.T) {
	rawData := ""
	blankHosts, err := NewHosts(&HostsConfig{
		RawText: &rawData,
	})
	if err != nil {
		t.Fatalf("TestNewHosts failed on NewHosts: %v", err)
	}

	testHostName := "test-new-hosts-default"
	blankHosts.AddHost("127.100.100.100", testHostName)
	foundHostnames := blankHosts.ListHostsByIP("127.100.100.100")
	found := false
	for _, hn := range foundHostnames {
		if hn == testHostName {
			found = true
			break
		}
	}

	if !found {
		t.Fatal("TestNewHosts was unable find the hostname it added")
	}

	renderedData := blankHosts.RenderHostsFile()
	fromRenderedHosts, err := NewHosts(&HostsConfig{
		RawText: &renderedData,
	})
	if err != nil {
		t.Fatalf("Failed to create hosts from rendered data: %v", err)
	}

	for _, k := range []string{"A", "T", "T", "X", "X", "E", "E", "H", "H"} {
		fromRenderedHosts.AddHost("127.100.100.100", k)
	}

	fromRenderedHosts.RemoveHosts([]string{"A "})

	err = fromRenderedHosts.Save()
	if err == nil {
		t.Fatalf("Should not be able to save when using RawText.")
	}

	foundHostnames = fromRenderedHosts.ListHostsByIP("127.100.100.100")
	found = false
	for _, hn := range foundHostnames {
		if hn == testHostName {
			found = true
			break
		}
	}

	if !found {
		t.Fatal("ListHostsByIP was unable find an IP")
	}
}

func TestMethods(t *testing.T) {
	mockHosts, err := NewHosts(&HostsConfig{
		RawText: &mockHostsData,
	})
	if err != nil {
		t.Fatalf("TestMockHosts failed on NewHosts: %v", err)
	}

	// get IP for localhost
	ok, ipString, _ := mockHosts.HostAddressLookup("localhost", IPFamilyV4)
	if !ok {
		t.Fatalf("TestMockHosts could not find IPv4 address for localhost")
	}
	if ipString != "127.0.0.1" {
		t.Fatalf("TestMockHosts returned incorrect IPv4 address for localhost")
	}

	// Test ListAddressesByHost exact
	hostList := mockHosts.ListAddressesByHost("nifi", true)
	if len(hostList) != 1 {
		t.Fatalf("ListAddressesByHost returned returned invalid number of hosts for exact match")
	}

	// Test ListAddressesByHost
	hostList = mockHosts.ListAddressesByHost("nifi", false)
	if len(hostList) != 70 {
		t.Fatalf("ListAddressesByHost returned returned invalid number of hosts")
	}

	// Test ListHostsByCIDR expect one
	hostList = mockHosts.ListHostsByCIDR("127.0.0.0/24")
	if len(hostList) != 1 {
		t.Fatalf("ListHostsByCIDR returned returned "+
			"invalid number of hosts for 127.0.0.0/24 got %v expecting 1", len(hostList))
	}

	// Test ListHostsByCIDR expect multiple
	hostList = mockHosts.ListHostsByCIDR("127.1.27.0/24")
	if len(hostList) != 252 {
		t.Fatalf("ListHostsByCIDR returned returned "+
			"invalid number of hosts for 127.0.0.0/24 got %v expecting 252", len(hostList))
	}

	// Test ListHostsByIP expect one
	ipHostList := mockHosts.ListHostsByIP("127.0.0.1")
	if len(ipHostList) != 1 {
		t.Fatalf("ListHostsByIP returned returned "+
			"invalid number of hosts for 127.0.0.1 got %v expecting 1", len(ipHostList))
	}

	// Test ListHostsByIP expect multiple
	ipHostList = mockHosts.ListHostsByIP("127.1.27.1")
	if len(ipHostList) != 7 {
		t.Fatalf("ListHostsByIP returned returned "+
			"invalid number of hosts for 127.1.27.1 got %v expecting 7", len(ipHostList))
	}

	// Test ListHostsByIP expect multiple
	ipHostList = mockHosts.ListHostsByIP("127.1.27.2")
	if len(ipHostList) != 7 {
		t.Fatalf("ListHostsByIP returned returned "+
			"invalid number of hosts for 127.1.27.2 got %v expecting 7", len(ipHostList))
	}

	// Test RemoveAddress
	mockHosts.RemoveAddress("127.1.27.1")
	ipHostList = mockHosts.ListHostsByIP("127.1.27.1")
	if len(ipHostList) != 0 {
		t.Fatalf("RemoveAddress then ListHostsByIP returned returned "+
			"invalid number of hosts for 127.1.27.1 got %v expecting 0", len(ipHostList))
	}

	// Get list of hosts for 127.1.27.2 and 127.1.27.3
	ipHostList2Len := len(mockHosts.ListHostsByIP("127.1.27.2"))
	ipHostList3Len := len(mockHosts.ListHostsByIP("127.1.27.3"))
	if ipHostList2Len+ipHostList3Len != 14 {
		t.Fatalf("ListHostsByIP returned returned "+
			"invalid number of hosts for 127.1.27.2 and 127.1.27.3 got %v expecting 14", ipHostList2Len+ipHostList3Len)
	}

	// Test RemoveAddresses
	mockHosts.RemoveAddresses([]string{"127.1.27.2", "127.1.27.3"})
	ipHostList2Len = len(mockHosts.ListHostsByIP("127.1.27.2"))
	ipHostList3Len = len(mockHosts.ListHostsByIP("127.1.27.3"))
	if ipHostList2Len+ipHostList3Len != 0 {
		t.Fatalf("ListHostsByIP returned returned "+
			"invalid number of hosts for 127.1.27.2 and 127.1.27.3 got %v expecting 0", ipHostList2Len+ipHostList3Len)
	}

	supersetAddresses := mockHosts.ListAddressesByHost("superset", true)
	if len(supersetAddresses) != 1 {
		t.Fatalf("ListAddressesByHost returned returned "+
			"invalid number of hosts for superset got %v expecting 1", len(supersetAddresses))
	}

	hostsAtAddress := mockHosts.ListHostsByIP(supersetAddresses[0][0])
	if len(hostsAtAddress) != 7 {
		t.Fatalf("ListHostsByIP returned returned "+
			"invalid number of hosts for %v got %v expecting 7", supersetAddresses[0][0], len(hostsAtAddress))
	}

	// Test RemoveHost
	mockHosts.RemoveHost("superset")
	hostsAtAddress = mockHosts.ListHostsByIP(supersetAddresses[0][0])
	if len(hostsAtAddress) != 6 {
		t.Fatalf("ListHostsByIP returned returned "+
			"invalid number of hosts for %v got %v expecting 7", supersetAddresses[0][0], len(hostsAtAddress))
	}

	// Test RemoveHosts
	mockHosts.RemoveHosts([]string{"superset.d4l-demo.svc", "superset.d4l-demo.svc.dwdev4"})
	hostsAtAddress = mockHosts.ListHostsByIP(supersetAddresses[0][0])
	if len(hostsAtAddress) != 4 {
		t.Fatalf("ListHostsByIP returned returned "+
			"invalid number of hosts for %v got %v expecting 7", supersetAddresses[0][0], len(hostsAtAddress))
	}

	// Test Add Host
	ip := "127.1.27.5"
	expect := 7
	hostsAtAddress = mockHosts.ListHostsByIP(ip)
	if len(hostsAtAddress) != expect {
		t.Fatalf("ListHostsByIP returned returned "+
			"invalid number of hosts for %v got %v expecting %v", ip, len(hostsAtAddress), expect)
	}
	mockHosts.AddHost(ip, "test")
	expect = 8
	hostsAtAddress = mockHosts.ListHostsByIP(ip)
	if len(hostsAtAddress) != expect {
		t.Fatalf("ListHostsByIP returned returned "+
			"invalid number of hosts for %v got %v expecting %v", ip, len(hostsAtAddress), expect)
	}
	mockHosts.AddHosts(ip, []string{"test-0", "test-1", "test-2"})
	expect = 11
	hostsAtAddress = mockHosts.ListHostsByIP(ip)
	if len(hostsAtAddress) != expect {
		t.Fatalf("ListHostsByIP returned returned "+
			"invalid number of hosts for %v got %v expecting %v", ip, len(hostsAtAddress), expect)
	}

	// test line with comment and bad host should treat existing bad host as legit host
	// since this is not a "cleaner" utility
	ip = "127.1.27.24"
	expect = 7
	hostsAtAddress = mockHosts.ListHostsByIP(ip)
	if len(hostsAtAddress) != expect {
		t.Fatalf("ListHostsByIP returned returned "+
			"invalid number of hosts for %v got %v expecting %v", ip, len(hostsAtAddress), expect)
	}

	// remove addresses in 127.1.27.0/16 range
	if err := mockHosts.RemoveCIDRs([]string{"127.1.27.0/16"}); err != nil {
		t.Fatalf("RemoveCIDRs failed: %v", err)
	}

	hfl := strings.Split(mockHosts.RenderHostsFile(), "\n")
	lines := 19
	if len(hfl) != lines {
		t.Fatalf("Expeced %d lines and got %d", lines, len(hfl))
	}

	line := 17
	expectString := "# existing comment"
	if hfl[line] != expectString {
		t.Fatalf("Expeced \"%s\" on line %d. Got \"%s\"", expectString, line, hfl[line])
	}
}

func TestWinDefaultHostsFile(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("Skipping windows test")
	}

	tc := []struct {
		Name     string
		EnvVaule string
		Expect   string
	}{
		{
			Name:     "Correct SystemRoot environment variable",
			EnvVaule: "E:\\Windows",
			Expect:   "E:\\Windows\\System32\\drivers\\etc\\hosts",
		},
		{
			Name:     "SystemRoot with trailing slash",
			EnvVaule: "E:\\Windows\\",
			Expect:   "E:\\Windows\\System32\\drivers\\etc\\hosts",
		},
		{
			Name:     "No systemRoot environment variable",
			EnvVaule: "",
			Expect:   "C:\\Windows\\System32\\drivers\\etc\\hosts",
		},
	}
	for _, tt := range tc {
		t.Run(tt.Name, func(t *testing.T) {
			if err := os.Setenv("SystemRoot", tt.EnvVaule); err != nil {
				panic(err)
			}
			hostFile := winDefaultHostsFile()
			if hostFile != tt.Expect {
				t.Fatalf("TestWinDefaultHostsFile failed to get expected path: "+
					"expect: %v, got: %v", tt.Expect, hostFile)
			}
		})
	}
}

// =============================================================================
// Bug Regression Tests
// These tests document and verify bugs found during code analysis.
// Tests are designed to FAIL until the corresponding bugs are fixed.
// =============================================================================

// BUG #2: ParseHostsFromString drops last line if no trailing newline
// File: txeh.go:442
func TestParseHosts_NoTrailingNewline(t *testing.T) {
	input := "127.0.0.1 localhost"
	hosts, err := NewHosts(&HostsConfig{RawText: &input})
	if err != nil {
		t.Fatalf("Failed to parse hosts: %v", err)
	}

	result := hosts.ListHostsByIP("127.0.0.1")
	if len(result) != 1 {
		t.Errorf("Last line dropped when no trailing newline. Expected 1 host, got %d", len(result))
	}
}

func TestParseHosts_NoTrailingNewline_MultipleLines(t *testing.T) {
	input := "127.0.0.1 host1\n127.0.0.2 host2"
	hosts, err := NewHosts(&HostsConfig{RawText: &input})
	if err != nil {
		t.Fatalf("Failed to parse hosts: %v", err)
	}

	result := hosts.ListHostsByIP("127.0.0.2")
	if len(result) != 1 {
		t.Errorf("Last line dropped. Expected 1 host for 127.0.0.2, got %d", len(result))
	}
}

// BUG #3: GetHostFileLines returns pointer to internal data
// File: txeh.go:422-427
func TestGetHostFileLines_ReturnsInternalPointer(t *testing.T) {
	input := "127.0.0.1 original\n"
	hosts, err := NewHosts(&HostsConfig{RawText: &input})
	if err != nil {
		t.Fatalf("Failed to create hosts: %v", err)
	}

	lines := hosts.GetHostFileLines()
	if len(lines) > 0 {
		lines[0].Address = "192.168.1.1"
		lines[0].Hostnames = []string{"modified"}
	}

	result := hosts.RenderHostsFile()
	if strings.Contains(result, "modified") {
		t.Error("GetHostFileLines returns pointer to internal data, allowing external modification")
	}
}

func TestGetHostFileLines_ThreadSafety(t *testing.T) {
	input := "127.0.0.1 original\n"
	hosts, err := NewHosts(&HostsConfig{RawText: &input})
	if err != nil {
		t.Fatalf("Failed to create hosts: %v", err)
	}

	var wg sync.WaitGroup
	iterations := 100

	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < iterations; i++ {
			lines := hosts.GetHostFileLines()
			if len(lines) > 0 {
				// This now modifies a copy, not internal state
				lines[0].Hostnames = append(lines[0].Hostnames, "safe-copy")
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < iterations; i++ {
			hosts.AddHost("127.0.0.1", "safe")
		}
	}()

	wg.Wait()
	// With the fix, this should pass even with -race
}

// BUG #4: AddHost TOCTOU race condition
// File: txeh.go:263-325
func TestAddHost_TOCTOU_Race(t *testing.T) {
	input := "127.0.0.1 existinghost\n"
	hosts, err := NewHosts(&HostsConfig{RawText: &input})
	if err != nil {
		t.Fatalf("Failed to create hosts: %v", err)
	}

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				hosts.AddHost("127.0.0.1", "existinghost")
				hosts.AddHost("127.0.0.2", "existinghost")
			}
		}()
	}
	wg.Wait()
	// Run with -race to detect TOCTOU race
}

// BUG #5: ListHostsByCIDR panics on invalid CIDR
// File: txeh.go:371
func TestListHostsByCIDR_InvalidCIDR_Panics(t *testing.T) {
	input := "127.0.0.1 localhost\n"
	hosts, err := NewHosts(&HostsConfig{RawText: &input})
	if err != nil {
		t.Fatalf("Failed to create hosts: %v", err)
	}

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("ListHostsByCIDR panics on invalid CIDR: %v", r)
		}
	}()

	_ = hosts.ListHostsByCIDR("not-a-valid-cidr")
}

// BUG #6: Windows line endings not preserved
// File: txeh.go:437-438
func TestWindowsLineEndings_NotPreserved(t *testing.T) {
	input := "127.0.0.1 localhost\r\n192.168.1.1 server\r\n"
	hosts, err := NewHosts(&HostsConfig{RawText: &input})
	if err != nil {
		t.Fatalf("Failed to create hosts: %v", err)
	}

	rendered := hosts.RenderHostsFile()
	if !strings.Contains(rendered, "\r\n") {
		t.Log("Windows line endings (CRLF) not preserved in output")
	}
}

// =============================================================================
// Edge Case Tests
// =============================================================================

func TestParseHosts_EmptyFile(t *testing.T) {
	input := ""
	hosts, err := NewHosts(&HostsConfig{RawText: &input})
	if err != nil {
		t.Fatalf("Failed to create hosts from empty input: %v", err)
	}

	result := hosts.RenderHostsFile()
	if result != "" {
		t.Errorf("Empty input should produce empty output, got: %q", result)
	}
}

func TestParseHosts_OnlyComments(t *testing.T) {
	input := "# Comment 1\n# Comment 2\n"
	hosts, err := NewHosts(&HostsConfig{RawText: &input})
	if err != nil {
		t.Fatalf("Failed to create hosts: %v", err)
	}

	rendered := hosts.RenderHostsFile()
	if !strings.Contains(rendered, "# Comment 1") {
		t.Error("Comments not preserved")
	}
}

func TestParseHosts_InlineComment(t *testing.T) {
	input := "127.0.0.1 localhost # this is a comment\n"
	hosts, err := NewHosts(&HostsConfig{RawText: &input})
	if err != nil {
		t.Fatalf("Failed to create hosts: %v", err)
	}

	result := hosts.ListHostsByIP("127.0.0.1")
	if len(result) != 1 || result[0] != "localhost" {
		t.Errorf("Expected ['localhost'], got %v", result)
	}

	rendered := hosts.RenderHostsFile()
	if !strings.Contains(rendered, "this is a comment") {
		t.Error("Inline comment not preserved")
	}
}

func TestParseHosts_TabSeparated(t *testing.T) {
	input := "127.0.0.1\tlocalhost\thost2\n"
	hosts, err := NewHosts(&HostsConfig{RawText: &input})
	if err != nil {
		t.Fatalf("Failed to create hosts: %v", err)
	}

	result := hosts.ListHostsByIP("127.0.0.1")
	if len(result) != 2 {
		t.Errorf("Expected 2 hosts (tab-separated), got %d", len(result))
	}
}

func TestParseHosts_IPv6_Loopback(t *testing.T) {
	input := "::1 localhost\n"
	hosts, err := NewHosts(&HostsConfig{RawText: &input})
	if err != nil {
		t.Fatalf("Failed to create hosts: %v", err)
	}

	result := hosts.ListHostsByIP("::1")
	if len(result) != 1 {
		t.Errorf("IPv6 loopback not handled correctly: %v", result)
	}
}

func TestParseHosts_IPv6_Full(t *testing.T) {
	input := "2001:0db8:85a3:0000:0000:8a2e:0370:7334 ipv6host\n"
	hosts, err := NewHosts(&HostsConfig{RawText: &input})
	if err != nil {
		t.Fatalf("Failed to create hosts: %v", err)
	}

	result := hosts.ListHostsByIP("2001:0db8:85a3:0000:0000:8a2e:0370:7334")
	if len(result) != 1 {
		t.Errorf("Full IPv6 address not handled correctly: %v", result)
	}
}

// =============================================================================
// Phase 1: Library Coverage Tests
// =============================================================================

// --- Reload Tests ---

func TestReload_WithRawText_ReturnsError(t *testing.T) {
	input := "127.0.0.1 localhost\n"
	hosts, err := NewHosts(&HostsConfig{RawText: &input})
	if err != nil {
		t.Fatalf("Failed to create hosts: %v", err)
	}

	err = hosts.Reload()
	if err == nil {
		t.Error("Reload with RawText should return error")
		return
	}
	if err.Error() != "cannot call Reload with RawText" {
		t.Errorf("Unexpected error message: %v", err)
	}
}

func TestReload_Success(t *testing.T) {
	// Create a temp file
	tmpFile, err := os.CreateTemp("", "hosts_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer func() { _ = os.Remove(tmpFile.Name()) }()

	// Write initial content
	initialContent := "127.0.0.1 localhost\n"
	if _, err := tmpFile.WriteString(initialContent); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	_ = tmpFile.Close()

	// Load the hosts file
	hosts, err := NewHosts(&HostsConfig{ReadFilePath: tmpFile.Name()})
	if err != nil {
		t.Fatalf("Failed to load hosts: %v", err)
	}

	// Verify initial state
	result := hosts.ListHostsByIP("127.0.0.1")
	if len(result) != 1 || result[0] != "localhost" {
		t.Errorf("Initial state incorrect: %v", result)
	}

	// Modify the file externally
	newContent := "127.0.0.1 localhost newhost\n192.168.1.1 server\n"
	if err := os.WriteFile(tmpFile.Name(), []byte(newContent), 0644); err != nil {
		t.Fatalf("Failed to update temp file: %v", err)
	}

	// Reload
	if err := hosts.Reload(); err != nil {
		t.Fatalf("Reload failed: %v", err)
	}

	// Verify new state
	result = hosts.ListHostsByIP("127.0.0.1")
	if len(result) != 2 {
		t.Errorf("After reload expected 2 hosts, got %d: %v", len(result), result)
	}

	result = hosts.ListHostsByIP("192.168.1.1")
	if len(result) != 1 || result[0] != "server" {
		t.Errorf("After reload, new entry not found: %v", result)
	}
}

func TestReload_FileNotFound(t *testing.T) {
	// Create a temp file then delete it
	tmpFile, err := os.CreateTemp("", "hosts_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	_, _ = tmpFile.WriteString("127.0.0.1 localhost\n")
	_ = tmpFile.Close()

	hosts, err := NewHosts(&HostsConfig{ReadFilePath: tmpFile.Name()})
	if err != nil {
		t.Fatalf("Failed to load hosts: %v", err)
	}

	// Delete the file
	_ = os.Remove(tmpFile.Name())

	// Reload should fail
	err = hosts.Reload()
	if err == nil {
		t.Error("Reload should fail when file is deleted")
	}
}

// --- SaveAs Tests ---

func TestSaveAs_Success(t *testing.T) {
	input := "127.0.0.1 localhost\n"
	hosts, err := NewHosts(&HostsConfig{RawText: &input})
	if err != nil {
		t.Fatalf("Failed to create hosts: %v", err)
	}

	// This should fail because we're using RawText
	err = hosts.SaveAs("/tmp/test_hosts")
	if err == nil {
		t.Error("SaveAs with RawText should return error")
	}
}

func TestSaveAs_RealFile(t *testing.T) {
	// Create initial file
	tmpFile, err := os.CreateTemp("", "hosts_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer func() { _ = os.Remove(tmpFile.Name()) }()
	_, _ = tmpFile.WriteString("127.0.0.1 localhost\n")
	_ = tmpFile.Close()

	hosts, err := NewHosts(&HostsConfig{ReadFilePath: tmpFile.Name()})
	if err != nil {
		t.Fatalf("Failed to load hosts: %v", err)
	}

	// Add a host
	hosts.AddHost("192.168.1.1", "newserver")

	// Save to a different file
	outputFile, err := os.CreateTemp("", "hosts_output_*")
	if err != nil {
		t.Fatalf("Failed to create output file: %v", err)
	}
	defer func() { _ = os.Remove(outputFile.Name()) }()
	_ = outputFile.Close()

	if err := hosts.SaveAs(outputFile.Name()); err != nil {
		t.Fatalf("SaveAs failed: %v", err)
	}

	// Read and verify
	content, err := os.ReadFile(outputFile.Name())
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	if !strings.Contains(string(content), "newserver") {
		t.Error("Saved file should contain 'newserver'")
	}
}

func TestSaveAs_InvalidPath(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "hosts_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer func() { _ = os.Remove(tmpFile.Name()) }()
	_, _ = tmpFile.WriteString("127.0.0.1 localhost\n")
	_ = tmpFile.Close()

	hosts, err := NewHosts(&HostsConfig{ReadFilePath: tmpFile.Name()})
	if err != nil {
		t.Fatalf("Failed to load hosts: %v", err)
	}

	// Try to save to an invalid path
	err = hosts.SaveAs("/nonexistent/directory/hosts")
	if err == nil {
		t.Error("SaveAs to invalid path should fail")
	}
}

// --- AddHost Edge Cases ---

func TestAddHost_Localhost_AllowsDuplicates(t *testing.T) {
	input := "127.0.0.1 existinghost\n"
	hosts, err := NewHosts(&HostsConfig{RawText: &input})
	if err != nil {
		t.Fatalf("Failed to create hosts: %v", err)
	}

	// Add same host to different localhost addresses
	hosts.AddHost("127.0.0.1", "existinghost")
	hosts.AddHost("127.0.0.2", "existinghost")

	// For localhost addresses, duplicates should be allowed
	result1 := hosts.ListHostsByIP("127.0.0.1")
	result2 := hosts.ListHostsByIP("127.0.0.2")

	// The host should appear at both localhost addresses
	if len(result1) == 0 {
		t.Error("Host should exist at 127.0.0.1")
	}
	if len(result2) == 0 {
		t.Error("Host should exist at 127.0.0.2 (localhost allows duplicates)")
	}
}

func TestAddHost_IPv6_Basic(t *testing.T) {
	input := ""
	hosts, err := NewHosts(&HostsConfig{RawText: &input})
	if err != nil {
		t.Fatalf("Failed to create hosts: %v", err)
	}

	hosts.AddHost("::1", "ipv6localhost")
	hosts.AddHost("2001:db8::1", "ipv6server")

	result := hosts.ListHostsByIP("::1")
	if len(result) != 1 || result[0] != "ipv6localhost" {
		t.Errorf("IPv6 loopback not added correctly: %v", result)
	}

	result = hosts.ListHostsByIP("2001:db8::1")
	if len(result) != 1 || result[0] != "ipv6server" {
		t.Errorf("IPv6 address not added correctly: %v", result)
	}
}

func TestAddHost_IPv6_Reassignment(t *testing.T) {
	input := "2001:db8::1 myhost\n"
	hosts, err := NewHosts(&HostsConfig{RawText: &input})
	if err != nil {
		t.Fatalf("Failed to create hosts: %v", err)
	}

	// Reassign to different IPv6 address
	hosts.AddHost("2001:db8::2", "myhost")

	// Should be removed from old address
	result1 := hosts.ListHostsByIP("2001:db8::1")
	if len(result1) != 0 {
		t.Errorf("Host should be removed from old IPv6 address: %v", result1)
	}

	// Should be at new address
	result2 := hosts.ListHostsByIP("2001:db8::2")
	if len(result2) != 1 || result2[0] != "myhost" {
		t.Errorf("Host should be at new IPv6 address: %v", result2)
	}
}

func TestAddHost_InvalidIP_Ignored(t *testing.T) {
	input := ""
	hosts, err := NewHosts(&HostsConfig{RawText: &input})
	if err != nil {
		t.Fatalf("Failed to create hosts: %v", err)
	}

	// Try to add with invalid IP
	hosts.AddHost("not-an-ip", "testhost")
	hosts.AddHost("999.999.999.999", "testhost2")

	// Should have no entries
	lines := hosts.GetHostFileLines()
	if len(lines) != 0 {
		t.Errorf("Invalid IPs should be ignored, got %d entries", len(lines))
	}
}

func TestAddHost_EmptyHostname(t *testing.T) {
	input := ""
	hosts, err := NewHosts(&HostsConfig{RawText: &input})
	if err != nil {
		t.Fatalf("Failed to create hosts: %v", err)
	}

	hosts.AddHost("127.0.0.1", "")
	hosts.AddHost("127.0.0.1", "   ")

	result := hosts.ListHostsByIP("127.0.0.1")
	// Empty/whitespace hostnames are normalized and added
	// This documents current behavior
	t.Logf("Empty hostname behavior: %v", result)
}

func TestAddHost_CaseNormalization(t *testing.T) {
	input := ""
	hosts, err := NewHosts(&HostsConfig{RawText: &input})
	if err != nil {
		t.Fatalf("Failed to create hosts: %v", err)
	}

	hosts.AddHost("127.0.0.1", "MyHost")
	hosts.AddHost("127.0.0.1", "MYHOST")
	hosts.AddHost("127.0.0.1", "myhost")

	result := hosts.ListHostsByIP("127.0.0.1")
	// All should be normalized to lowercase, so only one entry
	if len(result) != 1 {
		t.Errorf("Case normalization should prevent duplicates, got: %v", result)
	}
	if result[0] != "myhost" {
		t.Errorf("Hostname should be lowercase, got: %s", result[0])
	}
}

func TestAddHost_SameAddressSameHost_NoOp(t *testing.T) {
	input := "127.0.0.1 existinghost\n"
	hosts, err := NewHosts(&HostsConfig{RawText: &input})
	if err != nil {
		t.Fatalf("Failed to create hosts: %v", err)
	}

	// Add same host to same address multiple times
	hosts.AddHost("127.0.0.1", "existinghost")
	hosts.AddHost("127.0.0.1", "existinghost")
	hosts.AddHost("127.0.0.1", "existinghost")

	result := hosts.ListHostsByIP("127.0.0.1")
	if len(result) != 1 {
		t.Errorf("Adding same host to same address should be no-op, got: %v", result)
	}
}

// --- ParseHosts (file-based) Tests ---

func TestParseHosts_FileNotFound(t *testing.T) {
	_, err := ParseHosts("/nonexistent/path/hosts")
	if err == nil {
		t.Error("ParseHosts should fail for nonexistent file")
	}
}

func TestParseHosts_ValidFile(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "hosts_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer func() { _ = os.Remove(tmpFile.Name()) }()

	content := "127.0.0.1 localhost\n192.168.1.1 server\n"
	_, _ = tmpFile.WriteString(content)
	_ = tmpFile.Close()

	lines, err := ParseHosts(tmpFile.Name())
	if err != nil {
		t.Fatalf("ParseHosts failed: %v", err)
	}

	if len(lines) != 2 {
		t.Errorf("Expected 2 lines, got %d", len(lines))
	}
}

// --- NewHosts Tests ---

func TestNewHosts_WithReadPath(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "hosts_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer func() { _ = os.Remove(tmpFile.Name()) }()

	_, _ = tmpFile.WriteString("127.0.0.1 localhost\n")
	_ = tmpFile.Close()

	hosts, err := NewHosts(&HostsConfig{ReadFilePath: tmpFile.Name()})
	if err != nil {
		t.Fatalf("NewHosts failed: %v", err)
	}

	result := hosts.ListHostsByIP("127.0.0.1")
	if len(result) != 1 {
		t.Errorf("Expected 1 host, got %d", len(result))
	}
}

func TestNewHosts_WithReadAndWritePath(t *testing.T) {
	// Create read file
	readFile, err := os.CreateTemp("", "hosts_read_*")
	if err != nil {
		t.Fatalf("Failed to create read file: %v", err)
	}
	defer func() { _ = os.Remove(readFile.Name()) }()
	_, _ = readFile.WriteString("127.0.0.1 localhost\n")
	_ = readFile.Close()

	// Create write file
	writeFile, err := os.CreateTemp("", "hosts_write_*")
	if err != nil {
		t.Fatalf("Failed to create write file: %v", err)
	}
	defer func() { _ = os.Remove(writeFile.Name()) }()
	_ = writeFile.Close()

	hosts, err := NewHosts(&HostsConfig{
		ReadFilePath:  readFile.Name(),
		WriteFilePath: writeFile.Name(),
	})
	if err != nil {
		t.Fatalf("NewHosts failed: %v", err)
	}

	hosts.AddHost("192.168.1.1", "newhost")
	if err := hosts.Save(); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// Read the write file and verify
	content, err := os.ReadFile(writeFile.Name())
	if err != nil {
		t.Fatalf("Failed to read write file: %v", err)
	}

	if !strings.Contains(string(content), "newhost") {
		t.Error("Write file should contain new host")
	}
}

func TestNewHosts_InvalidReadPath(t *testing.T) {
	_, err := NewHosts(&HostsConfig{ReadFilePath: "/nonexistent/path"})
	if err == nil {
		t.Error("NewHosts should fail with invalid read path")
	}
}

// --- HostAddressLookup Tests ---

func TestHostAddressLookup_IPv4Family(t *testing.T) {
	input := "127.0.0.1 myhost\n::1 myhost\n"
	hosts, err := NewHosts(&HostsConfig{RawText: &input})
	if err != nil {
		t.Fatalf("Failed to create hosts: %v", err)
	}

	// Look up with IPv4 family
	found, addr, _ := hosts.HostAddressLookup("myhost", IPFamilyV4)
	if !found {
		t.Error("Should find host in IPv4 family")
	}
	if addr != "127.0.0.1" {
		t.Errorf("Expected 127.0.0.1, got %s", addr)
	}

	// Look up with IPv6 family
	found, addr, _ = hosts.HostAddressLookup("myhost", IPFamilyV6)
	if !found {
		t.Error("Should find host in IPv6 family")
	}
	if addr != "::1" {
		t.Errorf("Expected ::1, got %s", addr)
	}
}

func TestHostAddressLookup_NotFound(t *testing.T) {
	input := "127.0.0.1 localhost\n"
	hosts, err := NewHosts(&HostsConfig{RawText: &input})
	if err != nil {
		t.Fatalf("Failed to create hosts: %v", err)
	}

	found, _, _ := hosts.HostAddressLookup("nonexistent", IPFamilyV4)
	if found {
		t.Error("Should not find nonexistent host")
	}
}

// --- RemoveCIDRs Edge Cases ---

func TestRemoveCIDRs_InvalidCIDR_ReturnsError(t *testing.T) {
	input := "127.0.0.1 localhost\n"
	hosts, err := NewHosts(&HostsConfig{RawText: &input})
	if err != nil {
		t.Fatalf("Failed to create hosts: %v", err)
	}

	err = hosts.RemoveCIDRs([]string{"invalid-cidr"})
	if err == nil {
		t.Error("RemoveCIDRs should return error for invalid CIDR")
	}
}

func TestRemoveCIDRs_MultipleCIDRs(t *testing.T) {
	input := "127.0.0.1 host1\n192.168.1.1 host2\n10.0.0.1 host3\n"
	hosts, err := NewHosts(&HostsConfig{RawText: &input})
	if err != nil {
		t.Fatalf("Failed to create hosts: %v", err)
	}

	// Remove multiple ranges
	err = hosts.RemoveCIDRs([]string{"127.0.0.0/8", "10.0.0.0/8"})
	if err != nil {
		t.Fatalf("RemoveCIDRs failed: %v", err)
	}

	// Only 192.168.1.1 should remain
	result := hosts.ListHostsByIP("192.168.1.1")
	if len(result) != 1 {
		t.Error("192.168.1.1 should remain")
	}

	result = hosts.ListHostsByIP("127.0.0.1")
	if len(result) != 0 {
		t.Error("127.0.0.1 should be removed")
	}

	result = hosts.ListHostsByIP("10.0.0.1")
	if len(result) != 0 {
		t.Error("10.0.0.1 should be removed")
	}
}

// --- isLocalhost Tests ---

func TestIsLocalhost(t *testing.T) {
	tests := []struct {
		address  string
		expected bool
	}{
		{"127.0.0.1", true},
		{"127.0.0.2", true},
		{"127.255.255.255", true},
		{"127.1.2.3", true},
		{"::1", true},
		{"192.168.1.1", false},
		{"10.0.0.1", false},
		{"0.0.0.0", false},
		{"128.0.0.1", false},
	}

	for _, tc := range tests {
		result := isLocalhost(tc.address)
		if result != tc.expected {
			t.Errorf("isLocalhost(%q) = %v, expected %v", tc.address, result, tc.expected)
		}
	}
}

// =============================================================================
// MaxHostsPerLine Tests
// Tests for the Windows hosts file line limit feature
// =============================================================================

func TestMaxHostsPerLine_ExplicitLimit(t *testing.T) {
	input := ""
	hosts, err := NewHosts(&HostsConfig{
		RawText:         &input,
		MaxHostsPerLine: 3,
	})
	if err != nil {
		t.Fatalf("Failed to create hosts: %v", err)
	}

	// Add 7 hosts to the same IP
	for i := 0; i < 7; i++ {
		hosts.AddHost("127.0.0.1", "host"+string(rune('a'+i)))
	}

	// Should have 3 lines with 3, 3, and 1 hosts respectively
	lines := hosts.GetHostFileLines()
	addressLines := 0
	totalHosts := 0
	for _, line := range lines {
		if line.Address == "127.0.0.1" {
			addressLines++
			totalHosts += len(line.Hostnames)
			if len(line.Hostnames) > 3 {
				t.Errorf("Line has %d hosts, expected <= 3", len(line.Hostnames))
			}
		}
	}

	if addressLines != 3 {
		t.Errorf("Expected 3 address lines, got %d", addressLines)
	}
	if totalHosts != 7 {
		t.Errorf("Expected 7 total hosts, got %d", totalHosts)
	}
}

func TestMaxHostsPerLine_Unlimited(t *testing.T) {
	input := ""
	hosts, err := NewHosts(&HostsConfig{
		RawText:         &input,
		MaxHostsPerLine: -1, // Force unlimited
	})
	if err != nil {
		t.Fatalf("Failed to create hosts: %v", err)
	}

	// Add 20 hosts to the same IP
	for i := 0; i < 20; i++ {
		hosts.AddHost("127.0.0.1", "host"+string(rune('a'+i)))
	}

	// Should have 1 line with all 20 hosts
	lines := hosts.GetHostFileLines()
	addressLines := 0
	for _, line := range lines {
		if line.Address == "127.0.0.1" {
			addressLines++
			if len(line.Hostnames) != 20 {
				t.Errorf("Expected 20 hosts on line, got %d", len(line.Hostnames))
			}
		}
	}

	if addressLines != 1 {
		t.Errorf("Expected 1 address line, got %d", addressLines)
	}
}

func TestMaxHostsPerLine_AutoDetect_NonWindows(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("This test is for non-Windows platforms")
	}

	input := ""
	hosts, err := NewHosts(&HostsConfig{
		RawText:         &input,
		MaxHostsPerLine: 0, // Auto-detect
	})
	if err != nil {
		t.Fatalf("Failed to create hosts: %v", err)
	}

	// Add 20 hosts to the same IP
	for i := 0; i < 20; i++ {
		hosts.AddHost("127.0.0.1", "host"+string(rune('a'+i)))
	}

	// On non-Windows, should be unlimited (1 line with all hosts)
	lines := hosts.GetHostFileLines()
	addressLines := 0
	for _, line := range lines {
		if line.Address == "127.0.0.1" {
			addressLines++
		}
	}

	if addressLines != 1 {
		t.Errorf("On non-Windows with auto-detect, expected 1 line, got %d", addressLines)
	}
}

func TestMaxHostsPerLine_ExactlyAtLimit(t *testing.T) {
	input := ""
	hosts, err := NewHosts(&HostsConfig{
		RawText:         &input,
		MaxHostsPerLine: 5,
	})
	if err != nil {
		t.Fatalf("Failed to create hosts: %v", err)
	}

	// Add exactly 5 hosts
	for i := 0; i < 5; i++ {
		hosts.AddHost("127.0.0.1", "host"+string(rune('a'+i)))
	}

	// Should have exactly 1 line with 5 hosts
	lines := hosts.GetHostFileLines()
	if len(lines) != 1 {
		t.Errorf("Expected 1 line, got %d", len(lines))
	}
	if len(lines[0].Hostnames) != 5 {
		t.Errorf("Expected 5 hosts, got %d", len(lines[0].Hostnames))
	}

	// Add one more, should create new line
	hosts.AddHost("127.0.0.1", "hostf")

	lines = hosts.GetHostFileLines()
	if len(lines) != 2 {
		t.Errorf("After exceeding limit, expected 2 lines, got %d", len(lines))
	}
}

func TestMaxHostsPerLine_MultipleIPs(t *testing.T) {
	input := ""
	hosts, err := NewHosts(&HostsConfig{
		RawText:         &input,
		MaxHostsPerLine: 3,
	})
	if err != nil {
		t.Fatalf("Failed to create hosts: %v", err)
	}

	// Add hosts to multiple IPs
	for i := 0; i < 5; i++ {
		hosts.AddHost("127.0.0.1", "host1-"+string(rune('a'+i)))
		hosts.AddHost("127.0.0.2", "host2-"+string(rune('a'+i)))
	}

	// Each IP should have 2 lines (3 + 2 hosts)
	lines := hosts.GetHostFileLines()
	ip1Lines := 0
	ip2Lines := 0
	for _, line := range lines {
		switch line.Address {
		case "127.0.0.1":
			ip1Lines++
		case "127.0.0.2":
			ip2Lines++
		}
	}

	if ip1Lines != 2 {
		t.Errorf("Expected 2 lines for 127.0.0.1, got %d", ip1Lines)
	}
	if ip2Lines != 2 {
		t.Errorf("Expected 2 lines for 127.0.0.2, got %d", ip2Lines)
	}
}

func TestMaxHostsPerLine_HostLookupAcrossLines(t *testing.T) {
	input := ""
	hosts, err := NewHosts(&HostsConfig{
		RawText:         &input,
		MaxHostsPerLine: 2,
	})
	if err != nil {
		t.Fatalf("Failed to create hosts: %v", err)
	}

	// Add 5 hosts to same IP (will create 3 lines: 2, 2, 1)
	hosts.AddHost("127.0.0.1", "hosta")
	hosts.AddHost("127.0.0.1", "hostb")
	hosts.AddHost("127.0.0.1", "hostc")
	hosts.AddHost("127.0.0.1", "hostd")
	hosts.AddHost("127.0.0.1", "hoste")

	// All hosts should be found when listing by IP
	result := hosts.ListHostsByIP("127.0.0.1")
	if len(result) != 5 {
		t.Errorf("Expected 5 hosts, got %d: %v", len(result), result)
	}

	// Each individual host should be findable
	for _, hostname := range []string{"hosta", "hostb", "hostc", "hostd", "hoste"} {
		found, addr, _ := hosts.HostAddressLookup(hostname, IPFamilyV4)
		if !found {
			t.Errorf("Host %s not found", hostname)
		}
		if addr != "127.0.0.1" {
			t.Errorf("Host %s has wrong address: %s", hostname, addr)
		}
	}
}

func TestMaxHostsPerLine_ReassignHostAcrossLines(t *testing.T) {
	input := ""
	hosts, err := NewHosts(&HostsConfig{
		RawText:         &input,
		MaxHostsPerLine: 2,
	})
	if err != nil {
		t.Fatalf("Failed to create hosts: %v", err)
	}

	// Add hosts to first IP
	hosts.AddHost("127.0.0.1", "hosta")
	hosts.AddHost("127.0.0.1", "hostb")
	hosts.AddHost("127.0.0.1", "hostc") // This goes to second line

	// Verify hostc is at 127.0.0.1
	found, addr, _ := hosts.HostAddressLookup("hostc", IPFamilyV4)
	if !found || addr != "127.0.0.1" {
		t.Fatalf("hostc should be at 127.0.0.1, got: found=%v addr=%s", found, addr)
	}

	// Reassign hostc to different IP
	hosts.AddHost("192.168.1.1", "hostc")

	// Verify hostc is now at 192.168.1.1
	found, addr, _ = hosts.HostAddressLookup("hostc", IPFamilyV4)
	if !found || addr != "192.168.1.1" {
		t.Errorf("hostc should be at 192.168.1.1, got: found=%v addr=%s", found, addr)
	}

	// Verify 127.0.0.1 no longer has hostc
	result := hosts.ListHostsByIP("127.0.0.1")
	for _, h := range result {
		if h == "hostc" {
			t.Error("hostc should have been removed from 127.0.0.1")
		}
	}
}

func TestMaxHostsPerLine_AddExistingHostSameIP(t *testing.T) {
	input := ""
	hosts, err := NewHosts(&HostsConfig{
		RawText:         &input,
		MaxHostsPerLine: 2,
	})
	if err != nil {
		t.Fatalf("Failed to create hosts: %v", err)
	}

	// Add hosts to create multiple lines
	hosts.AddHost("127.0.0.1", "hosta")
	hosts.AddHost("127.0.0.1", "hostb")
	hosts.AddHost("127.0.0.1", "hostc") // Goes to second line

	// Try to add hosta again to same IP - should be no-op
	hosts.AddHost("127.0.0.1", "hosta")

	result := hosts.ListHostsByIP("127.0.0.1")
	if len(result) != 3 {
		t.Errorf("Adding existing host should be no-op, got %d hosts: %v", len(result), result)
	}

	// Count hosta occurrences
	count := 0
	for _, h := range result {
		if h == "hosta" {
			count++
		}
	}
	if count != 1 {
		t.Errorf("hosta should appear exactly once, got %d times", count)
	}
}

func TestMaxHostsPerLine_PreservesExistingLongLines(t *testing.T) {
	// Start with a hosts file that already has more than the limit
	input := "127.0.0.1 host1 host2 host3 host4 host5 host6 host7 host8 host9 host10\n"
	hosts, err := NewHosts(&HostsConfig{
		RawText:         &input,
		MaxHostsPerLine: 3,
	})
	if err != nil {
		t.Fatalf("Failed to create hosts: %v", err)
	}

	// Verify existing long line is preserved
	lines := hosts.GetHostFileLines()
	if len(lines) != 1 {
		t.Errorf("Existing line should be preserved as-is, got %d lines", len(lines))
	}
	if len(lines[0].Hostnames) != 10 {
		t.Errorf("Existing hosts should be preserved, got %d hosts", len(lines[0].Hostnames))
	}

	// Adding a new host to this IP should create a new line
	hosts.AddHost("127.0.0.1", "newhost")

	lines = hosts.GetHostFileLines()
	// Original long line still exists, plus new line for newhost
	if len(lines) != 2 {
		t.Errorf("Expected 2 lines after adding host, got %d", len(lines))
	}

	// Find the new line
	newHostFound := false
	for _, line := range lines {
		for _, h := range line.Hostnames {
			if h == "newhost" {
				newHostFound = true
				if len(line.Hostnames) != 1 {
					t.Errorf("New host should be on its own line, got %d hosts", len(line.Hostnames))
				}
			}
		}
	}
	if !newHostFound {
		t.Error("newhost not found after adding")
	}
}

func TestMaxHostsPerLine_LimitOfOne(t *testing.T) {
	input := ""
	hosts, err := NewHosts(&HostsConfig{
		RawText:         &input,
		MaxHostsPerLine: 1,
	})
	if err != nil {
		t.Fatalf("Failed to create hosts: %v", err)
	}

	// Add 3 hosts
	hosts.AddHost("127.0.0.1", "hosta")
	hosts.AddHost("127.0.0.1", "hostb")
	hosts.AddHost("127.0.0.1", "hostc")

	// Should have 3 separate lines
	lines := hosts.GetHostFileLines()
	if len(lines) != 3 {
		t.Errorf("With limit of 1, expected 3 lines, got %d", len(lines))
	}

	for _, line := range lines {
		if len(line.Hostnames) != 1 {
			t.Errorf("Each line should have exactly 1 host, got %d", len(line.Hostnames))
		}
	}
}

func TestMaxHostsPerLine_IPv6(t *testing.T) {
	input := ""
	hosts, err := NewHosts(&HostsConfig{
		RawText:         &input,
		MaxHostsPerLine: 2,
	})
	if err != nil {
		t.Fatalf("Failed to create hosts: %v", err)
	}

	// Add hosts to IPv6 address
	hosts.AddHost("::1", "ipv6host1")
	hosts.AddHost("::1", "ipv6host2")
	hosts.AddHost("::1", "ipv6host3")

	// Should have 2 lines for ::1
	lines := hosts.GetHostFileLines()
	ipv6Lines := 0
	totalHosts := 0
	for _, line := range lines {
		if line.Address == "::1" {
			ipv6Lines++
			totalHosts += len(line.Hostnames)
		}
	}

	if ipv6Lines != 2 {
		t.Errorf("Expected 2 lines for ::1, got %d", ipv6Lines)
	}
	if totalHosts != 3 {
		t.Errorf("Expected 3 total hosts for ::1, got %d", totalHosts)
	}
}

func TestMaxHostsPerLine_RenderHostsFile(t *testing.T) {
	input := ""
	hosts, err := NewHosts(&HostsConfig{
		RawText:         &input,
		MaxHostsPerLine: 3,
	})
	if err != nil {
		t.Fatalf("Failed to create hosts: %v", err)
	}

	// Add 5 hosts
	hosts.AddHost("127.0.0.1", "host1")
	hosts.AddHost("127.0.0.1", "host2")
	hosts.AddHost("127.0.0.1", "host3")
	hosts.AddHost("127.0.0.1", "host4")
	hosts.AddHost("127.0.0.1", "host5")

	rendered := hosts.RenderHostsFile()
	lines := strings.Split(strings.TrimSpace(rendered), "\n")

	if len(lines) != 2 {
		t.Errorf("Rendered output should have 2 lines, got %d: %v", len(lines), lines)
	}

	// Each line should have the IP and hosts
	for _, line := range lines {
		if !strings.HasPrefix(line, "127.0.0.1") {
			t.Errorf("Line should start with IP: %s", line)
		}
	}
}

func TestMaxHostsPerLine_ThreadSafety(t *testing.T) {
	input := ""
	hosts, err := NewHosts(&HostsConfig{
		RawText:         &input,
		MaxHostsPerLine: 3,
	})
	if err != nil {
		t.Fatalf("Failed to create hosts: %v", err)
	}

	var wg sync.WaitGroup
	iterations := 50

	// Multiple goroutines adding hosts
	for g := 0; g < 5; g++ {
		wg.Add(1)
		go func(goroutineID int) {
			defer wg.Done()
			for i := 0; i < iterations; i++ {
				hostname := "host" + string(rune('a'+goroutineID)) + string(rune('0'+i%10))
				hosts.AddHost("127.0.0.1", hostname)
			}
		}(g)
	}

	wg.Wait()

	// Verify no panics or data races (run with -race)
	result := hosts.ListHostsByIP("127.0.0.1")
	if len(result) == 0 {
		t.Error("Should have added some hosts")
	}

	// Verify no line exceeds the limit
	lines := hosts.GetHostFileLines()
	for _, line := range lines {
		if line.Address == "127.0.0.1" && len(line.Hostnames) > 3 {
			t.Errorf("Line exceeds limit of 3: %d hosts", len(line.Hostnames))
		}
	}
}

func TestMaxHostsPerLine_DefaultConstant(t *testing.T) {
	if DefaultMaxHostsPerLineWindows != 9 {
		t.Errorf("DefaultMaxHostsPerLineWindows should be 9, got %d", DefaultMaxHostsPerLineWindows)
	}
}

func TestGetEffectiveMaxHostsPerLine_NilConfig(t *testing.T) {
	// Create hosts without embedding config
	h := &Hosts{}

	result := h.getEffectiveMaxHostsPerLine()

	if runtime.GOOS == "windows" {
		if result != DefaultMaxHostsPerLineWindows {
			t.Errorf("On Windows with nil config, expected %d, got %d", DefaultMaxHostsPerLineWindows, result)
		}
	} else {
		if result != 0 {
			t.Errorf("On non-Windows with nil config, expected 0 (unlimited), got %d", result)
		}
	}
}

func TestMaxHostsPerLine_RemoveHostFromMultiLineIP(t *testing.T) {
	input := ""
	hosts, err := NewHosts(&HostsConfig{
		RawText:         &input,
		MaxHostsPerLine: 2,
	})
	if err != nil {
		t.Fatalf("Failed to create hosts: %v", err)
	}

	// Add 4 hosts (creates 2 lines with 2 hosts each)
	hosts.AddHost("127.0.0.1", "hosta")
	hosts.AddHost("127.0.0.1", "hostb")
	hosts.AddHost("127.0.0.1", "hostc")
	hosts.AddHost("127.0.0.1", "hostd")

	// Verify we have 4 hosts
	result := hosts.ListHostsByIP("127.0.0.1")
	if len(result) != 4 {
		t.Fatalf("Expected 4 hosts, got %d", len(result))
	}

	// Remove hostc
	hosts.RemoveHost("hostc")

	// Verify only 3 hosts remain
	result = hosts.ListHostsByIP("127.0.0.1")
	if len(result) != 3 {
		t.Errorf("After removal, expected 3 hosts, got %d: %v", len(result), result)
	}

	for _, h := range result {
		if h == "hostc" {
			t.Error("hostc should have been removed")
		}
	}
}

func TestMaxHostsPerLine_KubefwdScenario(t *testing.T) {
	// Simulates kubefwd adding many services to localhost
	input := ""
	hosts, err := NewHosts(&HostsConfig{
		RawText:         &input,
		MaxHostsPerLine: 9, // Windows limit
	})
	if err != nil {
		t.Fatalf("Failed to create hosts: %v", err)
	}

	// Add 25 services (like kubefwd would)
	services := []string{
		"api-gateway", "auth-service", "user-service", "order-service",
		"payment-service", "notification-service", "email-service", "sms-service",
		"logging-service", "monitoring-service", "cache-service", "queue-service",
		"worker-service", "scheduler-service", "config-service", "discovery-service",
		"gateway-admin", "metrics-service", "tracing-service", "health-service",
		"backup-service", "restore-service", "migration-service", "seed-service",
		"cleanup-service",
	}

	for _, svc := range services {
		hosts.AddHost("127.0.0.1", svc)
		hosts.AddHost("127.0.0.1", svc+".default")
		hosts.AddHost("127.0.0.1", svc+".default.svc")
		hosts.AddHost("127.0.0.1", svc+".default.svc.cluster.local")
	}

	// Total hosts: 25 * 4 = 100
	result := hosts.ListHostsByIP("127.0.0.1")
	if len(result) != 100 {
		t.Errorf("Expected 100 hosts, got %d", len(result))
	}

	// Verify no line exceeds 9 hosts
	lines := hosts.GetHostFileLines()
	for _, line := range lines {
		if line.Address == "127.0.0.1" && len(line.Hostnames) > 9 {
			t.Errorf("Line has %d hosts, exceeds Windows limit of 9", len(line.Hostnames))
		}
	}

	// Should have ceiling(100/9) = 12 lines
	ip127Lines := 0
	for _, line := range lines {
		if line.Address == "127.0.0.1" {
			ip127Lines++
		}
	}
	expectedLines := (100 + 8) / 9 // ceiling division
	if ip127Lines != expectedLines {
		t.Errorf("Expected %d lines for 100 hosts with limit 9, got %d", expectedLines, ip127Lines)
	}
}

// =============================================================================
// AddHostWithComment Tests
// Tests for the inline comment feature (Issue #30)
// =============================================================================

func TestAddHostWithComment_Basic(t *testing.T) {
	input := ""
	hosts, err := NewHosts(&HostsConfig{RawText: &input})
	if err != nil {
		t.Fatalf("Failed to create hosts: %v", err)
	}

	hosts.AddHostWithComment("127.0.0.1", "myhost", "added-by-ORG")

	rendered := hosts.RenderHostsFile()
	if !strings.Contains(rendered, "myhost") {
		t.Error("Host not found in rendered output")
	}
	if !strings.Contains(rendered, "added-by-ORG") {
		t.Error("Comment not found in rendered output")
	}
	if !strings.Contains(rendered, "#") {
		t.Error("Comment marker (#) not found in rendered output")
	}
}

func TestAddHostWithComment_EmptyComment(t *testing.T) {
	input := ""
	hosts, err := NewHosts(&HostsConfig{RawText: &input})
	if err != nil {
		t.Fatalf("Failed to create hosts: %v", err)
	}

	// Empty comment should work like AddHost
	hosts.AddHostWithComment("127.0.0.1", "host1", "")
	hosts.AddHost("127.0.0.1", "host2")

	// Both should be on the same line (same address, same empty comment)
	lines := hosts.GetHostFileLines()
	if len(lines) != 1 {
		t.Errorf("Expected 1 line (same address, no comment), got %d", len(lines))
	}
	if len(lines[0].Hostnames) != 2 {
		t.Errorf("Expected 2 hosts on line, got %d", len(lines[0].Hostnames))
	}
}

func TestAddHostWithComment_SameCommentAppends(t *testing.T) {
	input := ""
	hosts, err := NewHosts(&HostsConfig{RawText: &input})
	if err != nil {
		t.Fatalf("Failed to create hosts: %v", err)
	}

	hosts.AddHostWithComment("127.0.0.1", "host1", "my-app")
	hosts.AddHostWithComment("127.0.0.1", "host2", "my-app")
	hosts.AddHostWithComment("127.0.0.1", "host3", "my-app")

	// All should be on the same line
	lines := hosts.GetHostFileLines()
	if len(lines) != 1 {
		t.Errorf("Expected 1 line (same comment), got %d", len(lines))
	}
	if len(lines[0].Hostnames) != 3 {
		t.Errorf("Expected 3 hosts on line, got %d", len(lines[0].Hostnames))
	}
	if lines[0].Comment != "my-app" {
		t.Errorf("Expected comment 'my-app', got '%s'", lines[0].Comment)
	}
}

func TestAddHostWithComment_DifferentCommentCreatesNewLine(t *testing.T) {
	input := ""
	hosts, err := NewHosts(&HostsConfig{RawText: &input})
	if err != nil {
		t.Fatalf("Failed to create hosts: %v", err)
	}

	hosts.AddHostWithComment("127.0.0.1", "host1", "app-a")
	hosts.AddHostWithComment("127.0.0.1", "host2", "app-b")

	// Should be on different lines due to different comments
	lines := hosts.GetHostFileLines()
	if len(lines) != 2 {
		t.Errorf("Expected 2 lines (different comments), got %d", len(lines))
	}

	// Verify each line has correct comment
	comments := make(map[string]bool)
	for _, line := range lines {
		comments[line.Comment] = true
	}
	if !comments["app-a"] || !comments["app-b"] {
		t.Error("Expected both 'app-a' and 'app-b' comments")
	}
}

func TestAddHostWithComment_MixedWithAddHost(t *testing.T) {
	input := ""
	hosts, err := NewHosts(&HostsConfig{RawText: &input})
	if err != nil {
		t.Fatalf("Failed to create hosts: %v", err)
	}

	hosts.AddHost("127.0.0.1", "host1")                      // no comment
	hosts.AddHostWithComment("127.0.0.1", "host2", "my-app") // with comment
	hosts.AddHost("127.0.0.1", "host3")                      // no comment - should join host1

	lines := hosts.GetHostFileLines()
	if len(lines) != 2 {
		t.Errorf("Expected 2 lines, got %d", len(lines))
	}

	// Find line without comment - should have host1 and host3
	for _, line := range lines {
		if line.Comment == "" {
			if len(line.Hostnames) != 2 {
				t.Errorf("Line without comment should have 2 hosts, got %d", len(line.Hostnames))
			}
		}
		if line.Comment == "my-app" {
			if len(line.Hostnames) != 1 {
				t.Errorf("Line with comment should have 1 host, got %d", len(line.Hostnames))
			}
		}
	}
}

func TestAddHostsWithComment_Basic(t *testing.T) {
	input := ""
	hosts, err := NewHosts(&HostsConfig{RawText: &input})
	if err != nil {
		t.Fatalf("Failed to create hosts: %v", err)
	}

	hosts.AddHostsWithComment("127.0.0.1", []string{"host1", "host2", "host3"}, "bulk-add")

	lines := hosts.GetHostFileLines()
	if len(lines) != 1 {
		t.Errorf("Expected 1 line, got %d", len(lines))
	}
	if len(lines[0].Hostnames) != 3 {
		t.Errorf("Expected 3 hosts, got %d", len(lines[0].Hostnames))
	}
	if lines[0].Comment != "bulk-add" {
		t.Errorf("Expected comment 'bulk-add', got '%s'", lines[0].Comment)
	}
}

func TestAddHostWithComment_RespectsMaxHostsPerLine(t *testing.T) {
	input := ""
	hosts, err := NewHosts(&HostsConfig{
		RawText:         &input,
		MaxHostsPerLine: 2,
	})
	if err != nil {
		t.Fatalf("Failed to create hosts: %v", err)
	}

	// Add 5 hosts with same comment
	for i := 0; i < 5; i++ {
		hosts.AddHostWithComment("127.0.0.1", "host"+string(rune('a'+i)), "my-app")
	}

	// Should have 3 lines (2, 2, 1)
	lines := hosts.GetHostFileLines()
	if len(lines) != 3 {
		t.Errorf("Expected 3 lines, got %d", len(lines))
	}

	// All lines should have the same comment
	for _, line := range lines {
		if line.Comment != "my-app" {
			t.Errorf("Expected comment 'my-app', got '%s'", line.Comment)
		}
	}

	// Verify no line exceeds limit
	for _, line := range lines {
		if len(line.Hostnames) > 2 {
			t.Errorf("Line exceeds limit: %d hosts", len(line.Hostnames))
		}
	}
}

func TestAddHostWithComment_PreservesExistingComments(t *testing.T) {
	input := "127.0.0.1 existinghost # existing-comment\n"
	hosts, err := NewHosts(&HostsConfig{RawText: &input})
	if err != nil {
		t.Fatalf("Failed to create hosts: %v", err)
	}

	// Add to same address with same comment
	hosts.AddHostWithComment("127.0.0.1", "newhost", "existing-comment")

	lines := hosts.GetHostFileLines()
	if len(lines) != 1 {
		t.Errorf("Expected 1 line, got %d", len(lines))
	}
	if len(lines[0].Hostnames) != 2 {
		t.Errorf("Expected 2 hosts, got %d", len(lines[0].Hostnames))
	}

	// Add with different comment
	hosts.AddHostWithComment("127.0.0.1", "anotherhost", "new-comment")

	lines = hosts.GetHostFileLines()
	if len(lines) != 2 {
		t.Errorf("Expected 2 lines after adding different comment, got %d", len(lines))
	}
}

func TestAddHostWithComment_Rendered(t *testing.T) {
	input := ""
	hosts, err := NewHosts(&HostsConfig{RawText: &input})
	if err != nil {
		t.Fatalf("Failed to create hosts: %v", err)
	}

	hosts.AddHostWithComment("127.0.0.1", "myservice", "managed-by-kubefwd")

	rendered := hosts.RenderHostsFile()

	// Should contain the comment with # prefix
	expected := "127.0.0.1"
	if !strings.Contains(rendered, expected) {
		t.Errorf("Rendered output should contain '%s'", expected)
	}
	if !strings.Contains(rendered, "myservice") {
		t.Error("Rendered output should contain 'myservice'")
	}
	if !strings.Contains(rendered, "#managed-by-kubefwd") && !strings.Contains(rendered, "# managed-by-kubefwd") {
		t.Errorf("Rendered output should contain comment. Got: %s", rendered)
	}
}

func TestAddHostWithComment_WhitespaceHandling(t *testing.T) {
	input := ""
	hosts, err := NewHosts(&HostsConfig{RawText: &input})
	if err != nil {
		t.Fatalf("Failed to create hosts: %v", err)
	}

	// Comments with leading/trailing whitespace should be normalized
	hosts.AddHostWithComment("127.0.0.1", "host1", "  my-comment  ")
	hosts.AddHostWithComment("127.0.0.1", "host2", "my-comment")

	// Both should be on the same line (whitespace normalized)
	lines := hosts.GetHostFileLines()
	if len(lines) != 1 {
		t.Errorf("Expected 1 line (whitespace normalized), got %d", len(lines))
	}
}

func TestAddHostWithComment_IPv6(t *testing.T) {
	input := ""
	hosts, err := NewHosts(&HostsConfig{RawText: &input})
	if err != nil {
		t.Fatalf("Failed to create hosts: %v", err)
	}

	hosts.AddHostWithComment("::1", "ipv6host", "ipv6-comment")

	lines := hosts.GetHostFileLines()
	if len(lines) != 1 {
		t.Errorf("Expected 1 line, got %d", len(lines))
	}
	if lines[0].Comment != "ipv6-comment" {
		t.Errorf("Expected comment 'ipv6-comment', got '%s'", lines[0].Comment)
	}
}

func TestAddHostWithComment_InvalidIP(t *testing.T) {
	input := ""
	hosts, err := NewHosts(&HostsConfig{RawText: &input})
	if err != nil {
		t.Fatalf("Failed to create hosts: %v", err)
	}

	// Should be ignored like AddHost
	hosts.AddHostWithComment("not-an-ip", "host", "comment")

	lines := hosts.GetHostFileLines()
	if len(lines) != 0 {
		t.Errorf("Invalid IP should be ignored, got %d lines", len(lines))
	}
}

func TestAddHostWithComment_KubefwdScenario(t *testing.T) {
	input := ""
	hosts, err := NewHosts(&HostsConfig{
		RawText:         &input,
		MaxHostsPerLine: 9,
	})
	if err != nil {
		t.Fatalf("Failed to create hosts: %v", err)
	}

	// Simulate kubefwd adding services with a tracking comment
	services := []string{
		"api-gateway", "auth-service", "user-service",
		"api-gateway.default", "auth-service.default", "user-service.default",
	}

	hosts.AddHostsWithComment("127.0.0.1", services, "kubefwd-managed")

	result := hosts.ListHostsByIP("127.0.0.1")
	if len(result) != 6 {
		t.Errorf("Expected 6 hosts, got %d", len(result))
	}

	// All lines should have the kubefwd comment
	lines := hosts.GetHostFileLines()
	for _, line := range lines {
		if line.Address == "127.0.0.1" && line.Comment != "kubefwd-managed" {
			t.Errorf("Expected 'kubefwd-managed' comment, got '%s'", line.Comment)
		}
	}

	// Rendered output should show the comment
	rendered := hosts.RenderHostsFile()
	if !strings.Contains(rendered, "kubefwd-managed") {
		t.Error("Rendered output should contain 'kubefwd-managed'")
	}
}

func TestAddHostWithComment_ThreadSafety(t *testing.T) {
	input := ""
	hosts, err := NewHosts(&HostsConfig{RawText: &input})
	if err != nil {
		t.Fatalf("Failed to create hosts: %v", err)
	}

	var wg sync.WaitGroup
	iterations := 50

	// Multiple goroutines adding hosts with comments
	for g := 0; g < 5; g++ {
		wg.Add(1)
		go func(goroutineID int) {
			defer wg.Done()
			comment := "goroutine-" + string(rune('0'+goroutineID))
			for i := 0; i < iterations; i++ {
				hostname := "host" + string(rune('a'+goroutineID)) + string(rune('0'+i%10))
				hosts.AddHostWithComment("127.0.0.1", hostname, comment)
			}
		}(g)
	}

	wg.Wait()

	// Verify no panics or data races (run with -race)
	result := hosts.ListHostsByIP("127.0.0.1")
	if len(result) == 0 {
		t.Error("Should have added some hosts")
	}
}

func TestAddHostWithComment_ReassignmentPreservesComment(t *testing.T) {
	input := ""
	hosts, err := NewHosts(&HostsConfig{RawText: &input})
	if err != nil {
		t.Fatalf("Failed to create hosts: %v", err)
	}

	// Add host with comment
	hosts.AddHostWithComment("127.0.0.1", "myhost", "original-comment")

	// Reassign to different IP with different comment
	hosts.AddHostWithComment("192.168.1.1", "myhost", "new-comment")

	// Host should be at new IP
	found, addr, _ := hosts.HostAddressLookup("myhost", IPFamilyV4)
	if !found || addr != "192.168.1.1" {
		t.Errorf("Host should be at 192.168.1.1, got %s", addr)
	}

	// Old IP should have no hosts
	result := hosts.ListHostsByIP("127.0.0.1")
	if len(result) != 0 {
		t.Errorf("Old IP should have no hosts, got %d", len(result))
	}
}

// TestListHostsByComment tests listing hosts by their comment
func TestListHostsByComment(t *testing.T) {
	hostsData := `127.0.0.1        localhost
127.0.0.1        app1 app2 # dev services
127.0.0.1        api1 api2 # dev services
192.168.1.1      staging # staging env
`
	hosts, err := NewHosts(&HostsConfig{RawText: &hostsData})
	if err != nil {
		t.Fatal(err)
	}

	// List hosts with "dev services" comment
	result := hosts.ListHostsByComment("dev services")
	if len(result) != 4 {
		t.Errorf("Expected 4 hosts with 'dev services' comment, got %d", len(result))
	}

	// List hosts with "staging env" comment
	result = hosts.ListHostsByComment("staging env")
	if len(result) != 1 {
		t.Errorf("Expected 1 host with 'staging env' comment, got %d", len(result))
	}
	if result[0] != "staging" {
		t.Errorf("Expected 'staging', got %s", result[0])
	}

	// List hosts with non-existent comment
	result = hosts.ListHostsByComment("does not exist")
	if len(result) != 0 {
		t.Errorf("Expected 0 hosts, got %d", len(result))
	}

	// List hosts with empty comment (should return localhost)
	result = hosts.ListHostsByComment("")
	if len(result) != 1 {
		t.Errorf("Expected 1 host with no comment, got %d", len(result))
	}
}

// TestRemoveByComment tests removing all entries by comment
func TestRemoveByComment(t *testing.T) {
	hostsData := `127.0.0.1        localhost
127.0.0.1        app1 app2 # dev services
127.0.0.1        api1 api2 # dev services
192.168.1.1      staging # staging env
`
	hosts, err := NewHosts(&HostsConfig{RawText: &hostsData})
	if err != nil {
		t.Fatal(err)
	}

	// Remove all entries with "dev services" comment
	hosts.RemoveByComment("dev services")

	// Verify they're gone
	result := hosts.ListHostsByComment("dev services")
	if len(result) != 0 {
		t.Errorf("Expected 0 hosts after removal, got %d", len(result))
	}

	// Verify localhost is still there
	result = hosts.ListHostsByComment("")
	if len(result) != 1 {
		t.Errorf("Expected localhost to remain, got %d hosts", len(result))
	}

	// Verify staging is still there
	result = hosts.ListHostsByComment("staging env")
	if len(result) != 1 {
		t.Errorf("Expected staging to remain, got %d hosts", len(result))
	}
}

// TestRemoveByCommentEmpty tests removing entries with no comment
func TestRemoveByCommentEmpty(t *testing.T) {
	hostsData := `127.0.0.1        localhost
127.0.0.1        app1 # dev services
`
	hosts, err := NewHosts(&HostsConfig{RawText: &hostsData})
	if err != nil {
		t.Fatal(err)
	}

	// Remove all entries with no comment
	hosts.RemoveByComment("")

	// Verify localhost is gone
	result := hosts.ListHostsByComment("")
	if len(result) != 0 {
		t.Errorf("Expected 0 hosts with empty comment, got %d", len(result))
	}

	// Verify app1 is still there
	result = hosts.ListHostsByComment("dev services")
	if len(result) != 1 {
		t.Errorf("Expected app1 to remain, got %d hosts", len(result))
	}
}

// TestListHostsByComment_WhitespaceHandling tests that whitespace in comments is handled correctly
func TestListHostsByComment_WhitespaceHandling(t *testing.T) {
	hostsData := `127.0.0.1        app1 # dev services
127.0.0.1        app2 #   dev services
127.0.0.1        app3 # dev services
`
	hosts, err := NewHosts(&HostsConfig{RawText: &hostsData})
	if err != nil {
		t.Fatal(err)
	}

	// All should match "dev services" after trimming
	result := hosts.ListHostsByComment("dev services")
	if len(result) != 3 {
		t.Errorf("Expected 3 hosts (whitespace should be trimmed), got %d", len(result))
	}

	// Query with whitespace should also work
	result = hosts.ListHostsByComment("  dev services  ")
	if len(result) != 3 {
		t.Errorf("Expected 3 hosts (query whitespace should be trimmed), got %d", len(result))
	}
}

// TestMaxHostsPerLine_Windows9Limit tests the actual Windows 9 host limit scenario
func TestMaxHostsPerLine_Windows9Limit(t *testing.T) {
	input := ""
	hosts, err := NewHosts(&HostsConfig{
		RawText:         &input,
		MaxHostsPerLine: 9, // Windows limit
	})
	if err != nil {
		t.Fatal(err)
	}

	// Add 20 hosts to same IP (should create 3 lines: 9+9+2)
	for i := 1; i <= 20; i++ {
		hosts.AddHost("127.0.0.1", fmt.Sprintf("svc%d", i))
	}

	// Count the lines with this address
	lines := hosts.GetHostFileLines()
	addressLines := 0
	for _, line := range lines {
		if line.Address == "127.0.0.1" {
			addressLines++
			// Each line should have at most 9 hosts
			if len(line.Hostnames) > 9 {
				t.Errorf("Line has %d hosts, expected max 9", len(line.Hostnames))
			}
		}
	}

	// Should have 3 lines (9+9+2=20)
	if addressLines != 3 {
		t.Errorf("Expected 3 lines for 20 hosts with limit of 9, got %d", addressLines)
	}

	// Verify all hosts are there
	result := hosts.ListHostsByIP("127.0.0.1")
	if len(result) != 20 {
		t.Errorf("Expected 20 hosts, got %d", len(result))
	}
}

// TestRemoveByComment_WhitespaceHandling tests whitespace handling in RemoveByComment
func TestRemoveByComment_WhitespaceHandling(t *testing.T) {
	hostsData := `127.0.0.1        app1 #   spaced comment
`
	hosts, err := NewHosts(&HostsConfig{RawText: &hostsData})
	if err != nil {
		t.Fatal(err)
	}

	// Remove with untrimmed query should still work
	hosts.RemoveByComment("  spaced comment  ")

	result := hosts.ListHostsByIP("127.0.0.1")
	if len(result) != 0 {
		t.Errorf("Expected 0 hosts after removal, got %d", len(result))
	}
}

// TestRemoveByComments tests removing entries by multiple comments
func TestRemoveByComments(t *testing.T) {
	hostsData := `127.0.0.1        localhost
127.0.0.1        app1 # dev services
127.0.0.1        app2 # staging env
127.0.0.1        app3 # production
`
	hosts, err := NewHosts(&HostsConfig{RawText: &hostsData})
	if err != nil {
		t.Fatal(err)
	}

	// Remove multiple comments at once
	hosts.RemoveByComments([]string{"dev services", "staging env"})

	// Verify dev services and staging env are gone
	result := hosts.ListHostsByComment("dev services")
	if len(result) != 0 {
		t.Errorf("Expected 0 hosts with 'dev services', got %d", len(result))
	}
	result = hosts.ListHostsByComment("staging env")
	if len(result) != 0 {
		t.Errorf("Expected 0 hosts with 'staging env', got %d", len(result))
	}

	// Verify localhost and production are still there
	result = hosts.ListHostsByComment("")
	if len(result) != 1 {
		t.Errorf("Expected localhost to remain, got %d hosts", len(result))
	}
	result = hosts.ListHostsByComment("production")
	if len(result) != 1 {
		t.Errorf("Expected production to remain, got %d hosts", len(result))
	}
}
