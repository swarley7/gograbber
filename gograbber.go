package main

import (
	"flag"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"time"

	"github.com/swarley7/gograbber/libgograbber"
	// "./libgograbber"
)

func parseCMDLine() *libgograbber.State {
	s := libgograbber.State{Ports: libgograbber.IntSet{Set: map[int]bool{}}}
	s.Version = "Alpha (0.2a)"

	var ports string
	var wordlist string
	var statusCodesIgn string
	var statusCodes string
	var protocols string
	var timeout int
	var AdvancedUsage bool
	var easy bool
	var HostHeaderFile string
	libgograbber.InitLogger(os.Stdout, os.Stdout, os.Stdout, os.Stdout, os.Stderr)

	// Commandline arguments
	// Global
	flag.IntVar(&s.Threads, "t", 500, "Number of concurrent threads (actually goroutines, so not full OS threads). A typical system can support a couple of thousand with very little performance impact.")
	flag.IntVar(&timeout, "T", 2, "Timeout (seconds) for HTTP/TCP connections")

	flag.IntVar(&s.Jitter, "j", 0, "Be nice to serverz; introduce random delay (in ms) between requests")
	// flag.IntVar(&s.Sleep, "sleep", 0, "Minimum sleep (in ms) between requests")
	flag.BoolVar(&s.Debug, "debug", false, "Enable debug info (can be very noisy)")
	flag.IntVar(&s.VerbosityLevel, "v", 1, "Sets the logging/verbosity level. This isn't really working right now. I'm not sure I even want it to quite frankly.")

	// Scanner related
	flag.BoolVar(&s.Scan, "scan", false, "Enable TCP port scanner")

	flag.StringVar(&s.InputFile, "i", "", "Input filename of line seperated targets (hosts, IPs, CIDR ranges)")
	flag.StringVar(&ports, "p", "80,443", "Comma-separated ports/ranges to test with port scanner or directory bruteforce. Predefined port ranges are defined by 'top', 'small', 'med', 'large', 'full'")

	// I am very drunk right now

	// Dirbust related
	flag.BoolVar(&s.Dirbust, "dirbust", false, "Perform dirbust-like directory brute force of hosts using provided wordlist")
	flag.StringVar(&HostHeaderFile, "H", "", "Optional: Supply a file containing custom host headers that you would like to issue with each request (maybe for bypassing WAF/CDN/VHOST garbage?)")
	flag.StringVar(&protocols, "P", "http,https", "If provided, each host will be tested for the given protocol")
	flag.StringVar(&statusCodesIgn, "s", "400,401,403,404,407,502", "HTTP Status codes to ignore")
	flag.StringVar(&statusCodes, "S", "200,301,302,405,500", "HTTP Status codes to record. Currently does NOTHING (dw about it)")
	flag.StringVar(&s.URLFile, "U", "", "Input filename of line seperated complete URLs to test (overwrites -i, -p, -P, -scan)")
	flag.StringVar(&s.SingleURL, "u", "", "Single input URL to test (overwrites -i, -p, -P, -scan)")
	flag.StringVar(&wordlist, "w", "", "Wordlist file containing line separated endpoints to directory bruteforce")
	flag.Float64Var(&s.Ratio, "r", 0.95, "Soft 404 detection comparison ratio.")
	flag.BoolVar(&s.Soft404Detection, "soft404", true, "Perform soft 404 detection")

	// Reporting
	flag.StringVar(&s.OutputDirectory, "o", "gograbber_output", "Directory to store output in")
	flag.StringVar(&s.ProjectName, "project", "hack", "Name this project (if you want, otherwise... whatever?)")

	// screenshot related

	flag.BoolVar(&s.Screenshot, "screenshot", false, "Take pretty pictures of discovered URLs")
	flag.IntVar(&s.NumPhantomProcs, "p_procs", 5, "Number of phantomjs processes to spawn; helps when you're trying to screenshot a ton of stuff at once.")
	flag.StringVar(&s.Cookies, "C", "", "Optional cookies to supply with each request. Provide as a semicolon separated string, e.g. \"'cookie1=value1;cookie2=value2'\" (so just like, copy paste from Burp)")
	flag.StringVar(&s.UserAgent, "ua", fmt.Sprintf("gograbber - %v - yeeee", s.Version), "Set a custom user agent")

	flag.IntVar(&s.ImgX, "img_x", 1024, "The width of screenshot images in pixels")
	flag.IntVar(&s.ImgY, "img_y", 800, "The height of screenshot images in pixels")
	flag.IntVar(&s.ScreenshotQuality, "Q", 50, "Screenshot quality as a percentage (higher means more megatronz per screenshot).")
	flag.StringVar(&s.PhantomJSPath, "phantomjs", "phantomjs", "Path to phantomjs binary for rendering web pages")
	flag.BoolVar(&AdvancedUsage, "hh", false, "Print advanced usage details with examples!")
	flag.BoolVar(&s.FollowRedirects, "fr", false, "Follow redirects")
	flag.StringVar(&s.ScreenshotFileType, "screenshot_ext", "png", "Filetype for screenshots (valid are pdf, png, jpg)")

	flag.BoolVar(&s.IgnoreSSLErrors, "k", true, "Ignore SSL/TLS cert validation errors (super secure amirite?). Look, if you're using this app you probably know the risks, and let's face it, dgaf.")

	flag.BoolVar(&easy, "easy", false, "Enables common scan options: '-scan -dirbust -screenshot -p top -P http,https -t 2000 -j 25 -p_procs 7 -T 20'")

	flag.Parse()
	libgograbber.InitColours()
	libgograbber.PrintBanner(&s)
	s.StartTime = time.Now()
	if s.Debug {
		go func() {
			libgograbber.Debug.Println("Profiler running on: localhost:6060")
			http.ListenAndServe("localhost:6060", nil)
		}()
	}
	libgograbber.Initialise(&s, ports, wordlist, statusCodesIgn, protocols, timeout, AdvancedUsage, easy, HostHeaderFile)
	return &s
}

func main() {
	//profiling code - handy when dealing with concurrency and deadlocks ._.

	state := parseCMDLine()
	// lib.PrintOpts(state)
	if state != nil {
		// dothething awww ye
		libgograbber.Start(*state)
	}
}
