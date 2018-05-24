package libgograbber

import (
	"crypto/tls"
	"flag"
	"fmt"
	"net/http"
	"os"
	"path"
	"strings"
	"sync"
	"time"

	phantomjs "github.com/swarley7/phantomjs"
)

// Initialise sets up the program's state
func Initialise(s *State, ports string, wordlist string, statusCodesIgn string, protocols string, timeout int, AdvancedUsage bool, easy bool) {

	if AdvancedUsage {

		var Usage = func() {
			fmt.Printf(LineSep())
			fmt.Fprintf(os.Stderr, "Advanced usage of %s:\n", os.Args[0])
			flag.PrintDefaults()
			fmt.Printf(LineSep())
			fmt.Printf("Examples for %s:\n", os.Args[0])
			fmt.Printf(">> Scan and dirbust the hosts from hosts.txt.\n")
			fmt.Printf("%v -i hosts.txt -w wordlist.txt -t 2000 -scan -dirbust\n", os.Args[0])
			fmt.Printf(">> Scan and dirbust the hosts from hosts.txt, and screenshot discovered web resources.\n")
			fmt.Printf("%v -i hosts.txt -w wordlist.txt -t 2000 -scan  -dirbust -screenshot\n", os.Args[0])
			fmt.Printf(">> Scan, dirbust, and screenshot the hosts from hosts.txt on common web application ports. Additionally, set the number of phantomjs processes to 3.\n")
			fmt.Printf("%v -i hosts.txt -w wordlist.txt -t 2000 -p_procs=3 -p top -scan -dirbust -screenshot\n", os.Args[0])
			fmt.Printf(">> Screenshot the URLs from urls.txt. Additionally, use a custom phantomjs path.\n")
			fmt.Printf("%v -U urls.txt -t 200 -j 400 -phantomjs /my/path/to/phantomjs -screenshot\n", os.Args[0])
			fmt.Printf(">> Screenshot the supplied URL. Additionally, use a custom phantomjs path.\n")
			fmt.Printf("%v -u http://example.com/test -t 200 -j 400 -phantomjs /my/path/to/phantomjs -screenshot\n", os.Args[0])
			fmt.Printf(">> EASY MODE/I DON'T WANT TO READ STUFF LEMME HACK OK?.\n")
			fmt.Printf("%v -i hosts.txt -w wordlist.txt -easy\n", os.Args[0])

			fmt.Printf(LineSep())
		}
		Usage()
		os.Exit(0)
	}
	if easy { // user wants to use easymode... lol?
		s.Timeout = 2
		s.Jitter = 25
		s.Scan = true
		s.Dirbust = true
		s.Screenshot = true
		s.Threads = 1000
		s.NumPhantomProcs = 7
		ports = "top"
	}

	s.Timeout = time.Duration(timeout) * time.Second

	tx = &http.Transport{
		DialContext:        (d).DialContext,
		DisableCompression: true,
		MaxIdleConns:       100,

		TLSClientConfig: &tls.Config{InsecureSkipVerify: s.IgnoreSSLErrors}}
	cl = http.Client{
		Transport: tx,
		Timeout:   s.Timeout,
	}
	// if s.FollowRedirect {
	// 	cl.CheckRedirect = func(req *http.Request, via []*http.Request) error {
	// 		return http.ErrUseLastResponse
	// 	}
	// }
	s.Targets = make(chan Host)

	s.URLProvided = false
	if s.URLFile != "" || s.SingleURL != "" {
		s.URLProvided = true
		s.Scan = false
	}

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

	if s.URLProvided { // A url and/or file full of urls was supplied - treat them as gospel

		go func() { // for reasons unknown this seems to work ok... other things dont. I don't understand computers
			defer close(s.Targets)
			if s.URLFile != "" {
				inputData, err := GetDataFromFile(s.URLFile)
				if err != nil {
					Error.Println(err)
					panic(err)
				}
				for _, item := range inputData {
					ParseURLToHost(item, s.Targets)
				}
			}
			if s.SingleURL != "" {
				s.URLProvided = true
				Info.Println(s.SingleURL)
				ParseURLToHost(s.SingleURL, s.Targets)
			}
		}()
		return
	}

	if ports != "" {
		if strings.ToLower(ports) == "full" {
			ports = full
		} else if strings.ToLower(ports) == "med" {
			ports = medium
		} else if strings.ToLower(ports) == "small" {
			ports = small
		} else if strings.ToLower(ports) == "large" {
			ports = large
		} else if strings.ToLower(ports) == "top" {
			ports = top
		}
		s.Ports = UnpackPortString(ports)

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
