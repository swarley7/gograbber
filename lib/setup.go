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
)

// Initialise sets up the program's state
func Initialise(s *State, ports string, wordlist string, statusCodesIgn string, protocols string, timeout int, AdvancedUsage bool) (errors *multierror.Error) {
	s.Targets = make(chan Host, s.Threads)
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
				fmt.Printf("Port: (%v) is invalid!\n", port)
				continue
			}
			s.Ports.Add(port)
		}
		// for p, _ := range s.Ports.Set {
		// 	fmt.Printf("%v\n", p)
		// }
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
				fmt.Printf("Target: %v\n", target)
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
	return
}

// Start does the thing
func Start(s State) {
	os.Mkdir(path.Join(s.OutputDirectory), 0755) // drwxr-xr-x
	cl.Timeout = s.Timeout * time.Second
	d.Timeout = s.Timeout * time.Second
	ScanChan := make(chan Host, s.Threads)
	DirbChan := make(chan Host, s.Threads)
	ScreenshotChan := make(chan Host, s.Threads)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go RoutineManager(&s, ScanChan, DirbChan, ScreenshotChan, &wg)
	wg.Wait()
	// go func() {
	// 	for _, host := range s.URLComponents {
	// 		s.Targets <- host
	// 	}
	// 	close(s.Targets)
	// }()
	// if s.Scan {
	// 	fmt.Printf(LineSep())

	// 	fmt.Printf("Starting Port Scanner\n")
	// 	if s.Debug {
	// 		fmt.Printf("Testing %v host:port combinations\n", len(s.URLComponents))
	// 	}
	// 	fmt.Printf(LineSep())
	// 	wg.Add(1)
	// 	go ScanHosts(&s, s.Targets, ScanChan, &wg)

	// 	fmt.Printf(LineSep())
	// } else {
	// 	ScanChan = s.Targets
	// }
	// if s.Dirbust {
	// 	fmt.Printf("Starting Dirbuster\n")
	// 	// if s.Debug {
	// 	// 	var numURLs int
	// 	// 	if len(s.Paths.Set) != 0 {
	// 	// 		numURLs = len(s.URLComponents) * len(s.Paths.Set)
	// 	// 	} else {
	// 	// 		numURLs = len(s.URLComponents)
	// 	// 	}
	// 	// 	fmt.Printf("Testing %v URLs\n", numURLs)
	// 	// }
	// 	fmt.Printf(LineSep())
	// 	wg.Add(1)
	// 	go DirbustHosts(&s, ScanChan, DirbChan, &wg)
	// 	if s.Debug {
	// 		fmt.Println(s.URLComponents)
	// 	}
	// 	fmt.Printf(LineSep())
	// } else {
	// 	s.Targets = ScanChan
	// 	DirbChan = s.Targets
	// }
	// if s.Screenshot {
	// 	fmt.Printf("Starting Screenshotter\n")
	// 	// Allocate phantom processes sensibly
	// 	numTargets := (len(s.URLComponents) * len(s.Paths.Set)) / 10
	// 	var numProcs = 10
	// 	if numTargets <= 10 {
	// 		numProcs = 1
	// 	} else if x := s.Threads / 10; numTargets > x {
	// 		numProcs = x
	// 	}
	// 	procs := make([]phantomjs.Process, numProcs)
	// 	if s.Debug {
	// 		fmt.Printf("Creating [%v] PhantomJS processes... This could take a second\n", numProcs)
	// 	}
	// 	for i := 0; i < numProcs; i++ {
	// 		procs[i] = phantomjs.Process{BinPath: s.PhantomJSPath,
	// 			Port:   phantomjs.DefaultPort + i,
	// 			Stdout: os.Stdout,
	// 			Stderr: os.Stderr}
	// 		if err := procs[i].Open(); err != nil {
	// 			panic(err)
	// 		}
	// 		fmt.Printf("-> Process: #[%v] of [%v] created on localhost:%v\n", i, numProcs, phantomjs.DefaultPort+i)
	// 		defer procs[i].Close()
	// 	}
	// 	s.PhantomProcesses = procs

	// 	s.ScreenshotDirectory = path.Join(s.OutputDirectory, "screenshots")
	// 	os.Mkdir(s.ScreenshotDirectory, 0755) // drwxr-xr-x
	// 	if s.Debug {
	// 		fmt.Printf("Testing %v URLs\n", len(s.URLComponents)*len(s.Paths.Set))
	// 	}
	// 	fmt.Printf(LineSep())
	// 	wg.Add(1)
	// 	go Screenshot(&s, DirbChan, ScreenshotChan, &wg)
	// 	fmt.Printf(LineSep())
	// } else {
	// 	s.Targets = DirbChan
	// 	ScreenshotChan = s.Targets
	// }
	// wg.Wait()
	// fmt.Printf("Starting Reporter\n")
	// if s.Debug {
	// 	fmt.Printf("Reporting on %v URLs\n", len(s.URLComponents)*len(s.Paths.Set))
	// }
	fmt.Printf(LineSep())
	s.ReportDirectory = path.Join(s.OutputDirectory, "report")
	os.Mkdir(s.ReportDirectory, 0755) // drwxr-xr-x
	reportFile := MarkdownReport(&s, ScreenshotChan)
	fmt.Printf("Report written to: [%v]\n", reportFile)
	fmt.Printf(LineSep())
}
