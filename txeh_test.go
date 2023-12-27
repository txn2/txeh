package txeh

import (
	"strings"
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
		t.Fatalf("TestNewHosts failed on NewHosts: %v", err)
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
	mockHosts.RemoveCIDRs([]string{"127.1.27.0/16"})

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

	mockHosts.UpdateHost("127.0.0.1", "127.0.0.2", "google.com", "Testing Testing")

	hfl = strings.Split(mockHosts.RenderHostsFile(), "\n")
	line = 0
	expectString = "127.0.0.2        google.com #Testing Testing"
	if hfl[line] != expectString {
		t.Fatalf("Expeced \"%s\" on line %d. Got \"%s\"", expectString, line, hfl[line])
	}
}
