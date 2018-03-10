package main

import (
	"flag"

	"gograbber/libgograbber"
)

func parseCMDLine() *libgograbber.State {
	argPath := flag.String("path", "", "the path(s) to test for each host")
	argPathFile := flag.String("path_file", "", "Provide the path to a file containing line separated list of endpoint paths to test per host")
	argUrls := flag.String("urls", "", "Provide space separated list of urls to grab eh")
	argUrlFile := flag.String("url_file", "", "Provide the path to a file containing line separated list of URLs")
	argHostFile := flag.String("host_file", "", "Provide the path to a file containing line separated list of IP addresses or CIDR networks to discover. WARNING: this option could cause mem leaks (e.g. don't add /8s pls)")
	flag.Parse()
}

func main() {
	state := parseCMDLine()
	if state != nil {
		libgograbber.Start(state)
	}
}
