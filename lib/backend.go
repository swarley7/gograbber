package lib

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"

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
	s *State, ports string) (errors *multierror.Error) {

	return
}

// LeftPad2Len https://github.com/DaddyOh/golang-samples/blob/master/pad.go
func LeftPad2Len(s string, padStr string, overallLen int) string {
	var padCountInt int
	padCountInt = 1 + ((overallLen - len(padStr)) / len(padStr))
	var retStr = strings.Repeat(padStr, padCountInt) + s
	return retStr[(len(retStr) - overallLen):]
}

// RightPad2Len https://github.com/DaddyOh/golang-samples/blob/master/pad.go
func RightPad2Len(s string, padStr string, overallLen int) string {
	var padCountInt int
	padCountInt = 1 + ((overallLen - len(padStr)) / len(padStr))
	var retStr = s + strings.Repeat(padStr, padCountInt)
	return retStr[:overallLen]
}

// readLines reads a whole file into memory
// and returns a slice of its lines.
func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

// GetDataFromFile reads line separated data into a string array
func GetDataFromFile(fileName string) (data []string, err error) {
	if fileName != "" {
		data, err := readLines(fileName)
		if err != nil {
			fmt.Printf("File: %v does not exist, or you do not have permz (%v)", fileName, err)
			return nil, err
		}
		return data, err
	}
	return
}
