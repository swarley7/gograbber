package main

import (
	"flag"
	"fmt"
	"gograbber/lib"
	"log"
	"net/http"
	_ "net/http/pprof"
)

func parseCMDLine() *lib.State {
	s := lib.State{Ports: lib.IntSet{Set: map[int]bool{}}}
	var ports string
	var wordlist string
	var statusCodesIgn string
	var protocols string
	// Commandline arguments
	flag.IntVar(&s.Threads, "t", 20, "Number of concurrent threads")
	flag.StringVar(&ports, "p", "80", "Comma-separated ports to test with port scanner or directory bruteforce (defaults to http/80)")
	flag.StringVar(&wordlist, "w", "", "Wordlist file containing line separated endpoints to directory bruteforce")

	flag.BoolVar(&s.Debug, "debug", false, "Enable debug info")
	flag.BoolVar(&s.Scan, "scan", false, "Enable host discovery/TCP port scanner")
	flag.BoolVar(&s.Dirbust, "d", false, "Perform dirbust-like directory brute force of hosts using provided wordlist")
	flag.BoolVar(&s.Screenshot, "screenshot", false, "Take pretty pictures of discovered URLs")

	flag.StringVar(&s.InputFile, "i", "", "Input filename of line seperated targets (hosts, IPs, CIDR ranges)")
	flag.StringVar(&s.URLFile, "u", "", "Input filename of line seperated complete URLs to test (overwrites -i, -p, -P, -w, --scan)")
	// I am very drunk right now
	flag.StringVar(&s.OutputFile, "o", "", "Output filename")
	flag.StringVar(&protocols, "P", "http,https", "If provided, each host will be tested for the given protocol")
	flag.BoolVar(&s.Quiet, "q", false, "Don't print the banner and other noise")
	flag.StringVar(&statusCodesIgn, "s", "401,403", "Output filename")

	flag.Parse()

	lib.PrintBanner(&s)
	if err := lib.Initialise(&s, ports, wordlist, statusCodesIgn, protocols); err.ErrorOrNil() != nil {
		fmt.Printf("%s\n", err.Error())
		return nil
	}
	if s.Debug {
		go func() {
			log.Println(http.ListenAndServe("localhost:6061", nil))
		}()
	}

	return &s
}

func main() {
	//profiling code - handy when dealing with concurrency and deadlocks ._.

	state := parseCMDLine()
	lib.PrintOpts(state)
	if state != nil {
		// dothething awww ye
		lib.Start(*state)
	}
}
