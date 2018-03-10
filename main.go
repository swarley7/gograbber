package main

import (
	"flag"
	"fmt"
	"gograbber/libgograbber"
)

func parseCMDLine() *libgograbber.State {
	var extensions string
	var codes string
	var proxy string

	s := libgograbber.InitState()

	// Set up the variables we're interested in parsing.
	flag.IntVar(&s.Threads, "t", 10, "Number of concurrent threads")
	// flag.StringVar(&s.Wordlist, "w", "", "Path to the wordlist")
	// flag.StringVar(&codes, "s", "200,204,301,302,307", "Positive status codes")
	// flag.StringVar(&s.OutputFileName, "o", "", "Output file to write results to (defaults to stdout)")
	// flag.StringVar(&s.Url, "u", "", "The target URL or Domain")
	// flag.StringVar(&s.Cookies, "c", "", "Cookies to use for the requests")
	// flag.StringVar(&s.Username, "U", "", "Username for Basic Auth")
	// flag.StringVar(&s.Password, "P", "", "Password for Basic Auth")
	flag.StringVar(&s.Ports, "p", "80", "Ports to test with port scanner or (defaults to http/80)")
	// flag.StringVar(&extensions, "x", "", "File extension(s) to search for")
	// flag.StringVar(&s.UserAgent, "a", "", "Set the User-Agent string")
	// flag.StringVar(&proxy, "proxy", "", "Proxy to use for requests [http(s)://host:port]")
	// flag.BoolVar(&s.Verbose, "v", false, "Verbose output (errors)")
	// flag.BoolVar(&s.Debug, "d", false, "Debug mode")
	// flag.BoolVar(&s.Scan, "scan", false, "Enable port scanner to find HTTP(S) services")
	// flag.BoolVar(&s.ShowIPs, "i", false, "Show IP addresses (dns mode only)")
	// flag.BoolVar(&s.FollowRedirect, "r", false, "Follow redirects")
	flag.BoolVar(&s.Quiet, "q", false, "Don't print the banner and other noise")
	// flag.BoolVar(&s.Expanded, "e", false, "Expanded mode, print full URLs")
	// flag.BoolVar(&s.NoStatus, "n", false, "Don't print status codes")
	// flag.BoolVar(&s.IncludeLength, "l", false, "Include the length of the body in the output")
	// flag.BoolVar(&s.UseSlash, "f", false, "Append a forward-slash to each directory request")
	// flag.BoolVar(&s.WildcardForced, "fw", false, "Force continued operation when wildcard found")
	// flag.BoolVar(&s.InsecureSSL, "k", false, "Skip SSL certificate verification")

	flag.Parse()

	libgograbber.PrintBanner(&s)

	if err := libgograbber.ValidateState(&s, extensions, codes, proxy); err.ErrorOrNil() != nil {
		fmt.Printf("%s\n", err.Error())
		return nil
	} else {
		libgograbber.Ruler(&s)
		return &s
	}
}

func main() {
	state := parseCMDLine()
	if state != nil {
		libgograbber.Start(state)
	}
}
