package main

import (
	"flag"
	"gograbber/lib"
	"net/http"
	_ "net/http/pprof"
	"os"
)

func parseCMDLine() *lib.State {
	s := lib.State{Ports: lib.IntSet{Set: map[int]bool{}}}
	var ports string
	var wordlist string
	var statusCodesIgn string
	var statusCodes string
	var protocols string
	var timeout int
	var AdvancedUsage bool
	lib.InitLogger(os.Stdout, os.Stdout, os.Stdout, os.Stdout, os.Stderr)

	// Commandline arguments
	// Global
	flag.IntVar(&s.Threads, "t", 20, "Number of concurrent threads")
	flag.IntVar(&timeout, "T", 3, "Timeout for HTTP connections")

	flag.IntVar(&s.Jitter, "j", 0, "Introduce random delay (in ms) between requests")
	// flag.IntVar(&s.Sleep, "sleep", 0, "Minimum sleep (in ms) between requests")
	flag.BoolVar(&s.Debug, "debug", false, "Enable debug info")
	flag.IntVar(&s.VerbosityLevel, "v", 1, "Sets the logging/verbosity level.")

	// Scanner related
	flag.BoolVar(&s.Scan, "scan", false, "Enable host discovery/TCP port scanner")

	flag.StringVar(&s.InputFile, "i", "", "Input filename of line seperated targets (hosts, IPs, CIDR ranges)")
	flag.StringVar(&ports, "p", "80,443", "Comma-separated ports to test with port scanner or directory bruteforce. Predefined port ranges are defined by 'top', 'small', 'med', 'large', 'full'")

	// I am very drunk right now

	// Dirbust related
	flag.BoolVar(&s.Dirbust, "dirbust", false, "Perform dirbust-like directory brute force of hosts using provided wordlist")

	flag.StringVar(&protocols, "P", "http,https", "If provided, each host will be tested for the given protocol")
	flag.StringVar(&statusCodesIgn, "s", "400,401,403,404,407,502", "HTTP Status codes to ignore")
	flag.StringVar(&statusCodes, "S", "200,301,302,405,500", "HTTP Status codes to record")
	flag.StringVar(&s.URLFile, "U", "", "Input filename of line seperated complete URLs to test (overwrites -i, -p, -P, -w, -scan)")
	flag.StringVar(&s.SingleURL, "u", "", "Single input URL to test (overwrites -i, -p, -P, -w, -scan)")
	flag.StringVar(&wordlist, "w", "", "Wordlist file containing line separated endpoints to directory bruteforce")
	flag.Float64Var(&s.Ratio, "r", 0.95, "Soft 404 detection comparison ratio.")
	flag.BoolVar(&s.Soft404Detection, "soft404", true, "Perform soft 404 detection")
	flag.IntVar(&s.Soft404Method, "m", 1, "1: simple (less overhead), 2: aggressive (a lot more requests!)")

	// Reporting
	flag.StringVar(&s.OutputDirectory, "o", "gograbber_output", "Directory to store output in")
	flag.StringVar(&s.ProjectName, "project", "hack", "Name this project (if you want, otherwise... whatever?)")

	// screenshot related

	flag.BoolVar(&s.Screenshot, "screenshot", false, "Take pretty pictures of discovered URLs")
	flag.IntVar(&s.NumPhantomProcs, "p_procs", 5, "Number of phantomjs processes to spawn; helps when you're trying to screenshot a ton of stuff at once.")

	flag.IntVar(&s.ImgX, "img_x", 1024, "The width of screenshot images in pixels")
	flag.IntVar(&s.ImgY, "img_y", 800, "The height of screenshot images in pixels")
	flag.IntVar(&s.ScreenshotQuality, "Q", 50, "Screenshot quality as a percentage (higher means more megatronz per screenshot).")
	flag.StringVar(&s.PhantomJSPath, "phantomjs", "phantomjs", "Path to phantomjs binary for rendering web pages")
	flag.BoolVar(&AdvancedUsage, "hh", false, "Print advanced usage details with examples!")
	flag.BoolVar(&s.IgnoreSSLErrors, "k", true, "Ignore SSL/TLS cert validation errors (super secure amirite?). Look, if you're using this app you probably know the risks, and let's face it, dgaf.")

	flag.Parse()
	lib.InitColours()
	lib.PrintBanner(&s)
	if err := lib.Initialise(&s, ports, wordlist, statusCodesIgn, protocols, timeout, AdvancedUsage); err.ErrorOrNil() != nil {
		lib.Error.Printf("%s\n", err.Error())
		return nil
	}
	if s.Debug {
		go func() {
			lib.Debug.Println(http.ListenAndServe("localhost:6060", nil))
		}()
	}
	return &s
}

func main() {
	//profiling code - handy when dealing with concurrency and deadlocks ._.

	state := parseCMDLine()
	// lib.PrintOpts(state)
	if state != nil {
		// dothething awww ye
		lib.Start(*state)
	}
}
