package lib

import (
	"flag"
	"fmt"
	"math"
	"os"
	"path"
	"strings"
	"sync"
	"time"

	multierror "github.com/hashicorp/go-multierror"
	phantomjs "github.com/swarley7/phantomjs"
)

// Initialise sets up the program's state
func Initialise(s *State, ports string, wordlist string, statusCodesIgn string, protocols string, timeout int, AdvancedUsage bool) (errors *multierror.Error) {

	s.Targets = make(chan Host)

	if AdvancedUsage {

		var Usage = func() {
			fmt.Printf(LineSep())
			fmt.Fprintf(os.Stderr, "Advanced usage of %s:\n", os.Args[0])
			flag.PrintDefaults()
			fmt.Printf(LineSep())
			fmt.Printf("Examples for %s:\n", os.Args[0])
			fmt.Println()
		}
		Usage()
		os.Exit(0)
	}
	s.PrefetchedHosts = map[string]bool{}
	s.Soft404edHosts = map[string]bool{}
	s.Timeout = time.Duration(timeout) * time.Second
	cl.Timeout = s.Timeout

	s.URLProvided = false
	if wordlist != "" {
		pathData, err := GetDataFromFile(wordlist)
		if err != nil {
			panic(err)
		}
		s.Paths = StringSet{Set: map[string]bool{}}
		for _, path := range pathData {
			s.Paths.Add(path)
		}
	}
	if s.URLFile != "" || s.SingleURL != "" { // A url and/or file full of urls was supplied - treat them as gospel
		s.URLProvided = true
		if s.URLFile != "" {
			s.Scan = false
			inputData, err := GetDataFromFile(s.URLFile)
			if err != nil {
				panic(err)
			}

			for _, item := range inputData {
				// s.URLComponents
				ParseURLToHost(item, s.Targets)

			}
		}
		if s.SingleURL != "" {
			ParseURLToHost(s.SingleURL, s.Targets)
		}
		// close(s.Targets)
		s.Scan = false
		return
	}
	if ports != "" {

		for _, port := range StrArrToInt(strings.Split(ports, ",")) {
			if v := int(math.Pow(2, 16.0)); 0 > port || port >= v {
				Error.Printf("Port: (%v) is invalid!\n", port)
				continue
			}
			s.Ports.Add(port)
		}

	}
	if s.InputFile != "" {
		inputData, err := GetDataFromFile(s.InputFile)
		if err != nil {
			panic(err)
		}
		// c := make(chan StringSet)
		targetList := ExpandHosts(inputData)
		//  := <-c
		if s.Debug {
			for target := range targetList.Set {
				Debug.Printf("Target: %v\n", target)
			}
		}
		s.Hosts = targetList
	}
	s.StatusCodesIgn = IntSet{map[int]bool{}}
	for _, code := range StrArrToInt(strings.Split(statusCodesIgn, ",")) {
		s.StatusCodesIgn.Add(code)
	}

	// c2 := make(chan []Host)
	s.Protocols = StringSet{map[string]bool{}}
	for _, p := range strings.Split(protocols, ",") {
		s.Protocols.Add(p)
	}
	go GenerateURLs(s.Hosts, s.Ports, &s.Paths, s.Targets)
	// fmt.Println(s)
	if !s.Dirbust && !s.Scan && !s.Screenshot && !s.URLProvided {
		flag.Usage()
		os.Exit(1)
	}
	// close(s.Targets)
	return
}

// Start does the thing
func Start(s State) {
	fmt.Printf(LineSep())

	os.Mkdir(path.Join(s.OutputDirectory), 0755) // drwxr-xr-x
	// cl.Timeout = s.Timeout * time.Second
	// d.Timeout = 1 * time.Second
	// d.DisableKeepAlives()

	ScanChan := make(chan Host)
	DirbChan := make(chan Host)
	ScreenshotChan := make(chan Host)
	if s.Scan {
		s.ScanOutputDirectory = path.Join(s.OutputDirectory, "portscan")
		os.Mkdir(s.ScanOutputDirectory, 0755) // drwxr-xr-x
	}
	if s.Dirbust {
		s.HTTPResponseDirectory = path.Join(s.OutputDirectory, "raw_http_response")
		os.Mkdir(s.HTTPResponseDirectory, 0755) // drwxr-xr-x
		s.DirbustOutputDirectory = path.Join(s.OutputDirectory, "dirbust")
		os.Mkdir(s.DirbustOutputDirectory, 0755) // drwxr-xr-x
	}
	if s.Screenshot {
		procs := make([]phantomjs.Process, s.NumPhantomProcs)
		if s.Debug {
			Debug.Printf("Creating [%v] PhantomJS processes... This could take a second\n", s.NumPhantomProcs)
		}
		for i := 0; i < s.NumPhantomProcs; i++ {
			procs[i] = phantomjs.Process{BinPath: s.PhantomJSPath,
				Port:            phantomjs.DefaultPort + i,
				Stdout:          os.Stdout,
				Stderr:          os.Stderr,
				IgnoreSslErrors: s.IgnoreSSLErrors,
			}
			if err := procs[i].Open(); err != nil {
				panic(err)
			}
			Info.Printf("-> Phantomjs process: #[%v] (%v of %v) created on localhost:%v\n", i, i+1, s.NumPhantomProcs, phantomjs.DefaultPort+i)
			defer procs[i].Close()
		}
		s.PhantomProcesses = procs
		s.ScreenshotDirectory = path.Join(s.OutputDirectory, "screenshots")
		os.Mkdir(s.ScreenshotDirectory, 0755) // drwxr-xr-x
		if s.Debug {
			Debug.Printf("Testing %v URLs\n", len(s.URLComponents)*len(s.Paths.Set))
		}
	}
	wg := sync.WaitGroup{}
	wg.Add(1)
	go RoutineManager(&s, ScanChan, DirbChan, ScreenshotChan, &wg)

	s.ReportDirectory = path.Join(s.OutputDirectory, "report")
	os.Mkdir(s.ReportDirectory, 0755) // drwxr-xr-x
	reportFile := MarkdownReport(&s, ScreenshotChan)
	wg.Wait()

	Info.Printf("Report written to: [%v]\n", g.Sprintf("%s", reportFile))
	fmt.Printf(LineSep())
}
