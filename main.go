package main

import (
	"flag"
	"fmt"

	"gograbber/lib"
)

func parseCMDLine() *lib.State {
	s := lib.State{Ports: lib.IntSet{Set: map[int]bool{}}}
	var ports string
	// Commandline arguments
	flag.IntVar(&s.Threads, "t", 20, "Number of concurrent threads")
	flag.StringVar(&ports, "p", "80", "Comma-separated ports to test with port scanner or directory bruteforce (defaults to http/80)")
	flag.StringVar(&s.InputFile, "i", "", "Input filename of line seperated targets (hosts, IPs, CIDR ranges)")
	// I am very drunk right now
	flag.StringVar(&s.Protocol, "P", "", "If provided, each host will be tested for the given protocol")
	flag.BoolVar(&s.Quiet, "q", false, "Don't print the banner and other noise")

	flag.Parse()

	lib.PrintBanner(&s)
	if err := lib.Initialise(&s, ports); err.ErrorOrNil() != nil {
		fmt.Printf("%s\n", err.Error())
		return nil
	}
	return &s
}

func main() {
	state := parseCMDLine()
	lib.PrintOpts(state)
	if state != nil {
		// dothething awww ye
		lib.Start(*state)
	}
}
