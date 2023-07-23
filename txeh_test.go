package txeh

import (
	"testing"
)

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
