package txeh

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"
)

const UNKNOWN = 0
const EMPTY = 10
const COMMENT = 20
const ADDRESS = 30

type HostsConfig struct {
	FilePath string
}

type Hosts struct {
	HostsConfig
	hostFileLines HostFileLines
}

// AddressLocations the location of an address in the HFL
type AddressLocations map[string]int

// HostLocations maps a hostname
// to an original line number
type HostLocations map[string]int

type HostFileLines []HostFileLine

type HostFileLine struct {
	OriginalLineNum int
	LineType        int
	Address         string
	Parts           []string
	Hostnames       []string
	Raw             string
	Trimed          string
}

// NewHosts returns a new hosts object
func NewHosts(hc HostsConfig) *Hosts {
	h := &Hosts{HostsConfig: hc}

	if h.FilePath == "" {
		h.FilePath = "/etc/hosts"
	}

	h.hostFileLines = ParseHosts(h.FilePath)

	return h
}

// removeStringElement removed an element of a string slice
func removeStringElement(slice []string, s int) []string {
	return append(slice[:s], slice[s+1:]...)
}

// removeHFLElement removed an element of a HostFileLine slice
func removeHFLElement(slice []HostFileLine, s int) []HostFileLine {
	return append(slice[:s], slice[s+1:]...)
}

// AddHost adds a host to an address and removes the host
// from any existing address is may be associated with
func (h *Hosts) AddHost(address string, host string) {

	// does the host already exist
	if ok, exAdd, hflIdx := h.HostAddressLookup(host); ok {
		// if the address is the same we are done
		if address == exAdd {
			log.Print("Address / Host combination already exists. Do nothing.")
			return
		}

		// if the hostname is at a different address, go and remove it from the address
		for hidx, hst := range h.hostFileLines[hflIdx].Hostnames {
			if hst == host {
				h.hostFileLines[hflIdx].Hostnames = removeStringElement(h.hostFileLines[hflIdx].Hostnames, hidx)

				// remove the address line if empty
				if len(h.hostFileLines[hflIdx].Hostnames) < 1 {
					h.hostFileLines = removeHFLElement(h.hostFileLines, hflIdx)
				}
				break // unless we should continue because it could have duplicates
			}
		}
	}

	// if the address exists add it to the address line
	for i, hfl := range h.hostFileLines {
		if hfl.Address == address {
			h.hostFileLines[i].Hostnames = append(h.hostFileLines[i].Hostnames, host)
			return
		}
	}

	// the address and host do not already exist so go ahead and create them
	hfl := HostFileLine{
		LineType:  ADDRESS,
		Address:   address,
		Hostnames: []string{host},
	}

	h.hostFileLines = append(h.hostFileLines, hfl)
}

// HostAddressLookup returns true is the host is found, a string
// containing the address and the index of the hfl
func (h *Hosts) HostAddressLookup(host string) (bool, string, int) {
	for i, hfl := range h.hostFileLines {
		for _, hn := range hfl.Hostnames {
			if hn == strings.ToLower(host) {
				return true, hfl.Address, i
			}
		}
	}

	return false, "", 0
}

// addAddress adds and ip address and updates indexes
func (h *Hosts) addAddress(address string) int {
	h.hostFileLines = append(h.hostFileLines)
	return len(h.hostFileLines) + 1
}

func (h *Hosts) RenderHostsFile() string {
	hf := ""

	for _, hfl := range h.hostFileLines {
		hf = hf + fmt.Sprintln(lineFormatter(hfl))
	}

	return hf
}

func (h *Hosts) GetHostFileLines() *HostFileLines {
	return &h.hostFileLines
}

func ParseHosts(path string) []HostFileLine {
	input, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalln(err)
	}

	lines := strings.Split(string(input), "\n")

	hostFileLines := make([]HostFileLine, len(lines))

	// trim leading an trailing whitespace
	for i, l := range lines {
		curLine := &hostFileLines[i]
		curLine.OriginalLineNum = i
		curLine.Raw = l

		// trim line
		curLine.Trimed = strings.TrimSpace(l)

		// check for comment
		if strings.HasPrefix(curLine.Trimed, "#") {
			curLine.LineType = COMMENT
			continue
		}

		if curLine.Trimed == "" {
			curLine.LineType = EMPTY
			continue
		}

		curLine.Parts = strings.Fields(curLine.Trimed)

		if len(curLine.Parts) > 1 {
			curLine.LineType = ADDRESS
			curLine.Address = curLine.Parts[0]
			curLine.Hostnames = curLine.Parts[1:]
			continue
		}

		// if we can't figure out what this line is
		// at this point mark it as unknown
		curLine.LineType = UNKNOWN

	}

	return hostFileLines
}

func lineFormatter(hfl HostFileLine) string {

	if hfl.LineType < ADDRESS {
		return hfl.Raw
	}

	return fmt.Sprintf("%-15s %s", hfl.Address, strings.Join(hfl.Hostnames, " "))
}
