package lib

import (
	"fmt"
	"net"

	multierror "github.com/hashicorp/go-multierror"
)

// A single result which comes from an individual web
// request.
type Result struct {
	Entity string
	Status int
	Extra  string
	Size   *int64
}

type PrintResultFunc func(s *State, r *Result)
type ProcessorFunc func(s *State, entity string, resultChan chan<- Result)
type SetupFunc func(s *State) bool

// Shim type for "set" containing ints
type IntSet struct {
	Set map[int]bool
}

// Shim type for "set" containing strings
type StringSet struct {
	Set map[string]bool
}

func Hosts(cidr string) ([]string, error) {
	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}

	var ips []string
	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); inc(ip) {
		ips = append(ips, ip.String())
	}
	return ips, nil
}

//  http://play.golang.org/p/m8TNTtygK0
func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

// ExpandHosts takes a string array of IP addresses/CIDR masks and converts into a string array of pure IP addresses
func ExpandHosts(hosts []string) (allHosts []string) {
	for _, host := range hosts {
		ips, err := Hosts(host)
		if err != nil { // Not a CIDR... Might be a straight IP
			ip := net.ParseIP(host)
			if ip == nil {
				continue
			}
			allHosts = append(allHosts, ip.String())
		}
		allHosts = append(allHosts, ips...)
	}
	return allHosts
}

func Start(s State) {
	if s.Scan {
		openPorts := ScanHosts(&s)
		for socketPair := range openPorts.Set {
			fmt.Printf("Host:Port %s is open", socketPair)
		}
	}
}

func Initialise(
	s *State) (errors *multierror.Error) {

	return
}