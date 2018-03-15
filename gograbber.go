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
	var statusCodes string
	var protocols string
	// Commandline arguments
	flag.IntVar(&s.Threads, "t", 20, "Number of concurrent threads")
	flag.StringVar(&ports, "p", "80", "Comma-separated ports to test with port scanner or directory bruteforce")
	flag.StringVar(&wordlist, "w", "", "Wordlist file containing line separated endpoints to directory bruteforce")
	flag.IntVar(&s.Jitter, "j", 0, "Introduce random delay (in milliseconds) between requests")

	flag.BoolVar(&s.Debug, "debug", false, "Enable debug info")
	flag.BoolVar(&s.Scan, "scan", false, "Enable host discovery/TCP port scanner")
	flag.BoolVar(&s.Dirbust, "d", false, "Perform dirbust-like directory brute force of hosts using provided wordlist")
	flag.BoolVar(&s.Screenshot, "screenshot", false, "Take pretty pictures of discovered URLs")
	flag.IntVar(&s.ImgX, "img_x", 1024, "The width of screenshot images in pixels")
	flag.IntVar(&s.ImgY, "img_y", 800, "The height of screenshot images in pixels")

	flag.StringVar(&s.InputFile, "i", "", "Input filename of line seperated targets (hosts, IPs, CIDR ranges)")
	flag.StringVar(&s.URLFile, "U", "", "Input filename of line seperated complete URLs to test (overwrites -i, -p, -P, -w, --scan)")
	flag.StringVar(&s.SingleURL, "u", "", "Single input URL to test (overwrites -i, -p, -P, -w, --scan)")

	// I am very drunk right now
	flag.StringVar(&s.PhantomJSPath, "phantomjs", "phantomjs", "Path to phantomjs binary for rendering web pages")

	flag.StringVar(&s.OutputDirectory, "o", "", "Directory to store output in")
	flag.StringVar(&protocols, "P", "http,https", "If provided, each host will be tested for the given protocol")
	flag.BoolVar(&s.Quiet, "q", false, "Don't print the banner and other noise")
	flag.StringVar(&statusCodesIgn, "s", "401,403,404,407", "HTTP Status codes to ignore")
	flag.StringVar(&statusCodes, "S", "200,301,302,500", "HTTP Status codes to record")

	flag.Parse()

	lib.PrintBanner(&s)
	if err := lib.Initialise(&s, ports, wordlist, statusCodesIgn, protocols); err.ErrorOrNil() != nil {
		fmt.Printf("%s\n", err.Error())
		return nil
	}
	if s.Debug {
		go func() {
			log.Println(http.ListenAndServe("localhost:6060", nil))
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
