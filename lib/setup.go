package lib

import (
	"fmt"
	"math"
	"strings"

	multierror "github.com/hashicorp/go-multierror"
)

func Initialise(s *State, ports string, wordlist string, statusCodesIgn string, protocols string) (errors *multierror.Error) {
	if s.URLFile != "" {
		s.Scan = false
		inputData, err := GetDataFromFile(s.URLFile)
		if err != nil {
			panic(err)
		}

		for _, item := range inputData {
			// s.URLComponents
			h, err := ParseURLToHost(item)
			if err != nil {
				continue
			}
			s.URLComponents = append(s.URLComponents, h)
		}
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
	if s.InputFile != "" {
		inputData, err := GetDataFromFile(s.InputFile)
		if err != nil {
			panic(err)
		}
		c := make(chan StringSet)
		go ExpandHosts(inputData, c)
		targetList := <-c
		if s.Debug {
			for target := range targetList.Set {
				fmt.Printf("Target: %v\n", target)
			}
		}
		s.Hosts = targetList
	}
	s.StatusCodesIgn = IntSet{map[int]bool{}}
	for code, _ := range StrArrToInt(strings.Split(statusCodesIgn, ",")) {
		s.StatusCodesIgn.Add(code)
	}

	c2 := make(chan []Host)
	s.Protocols = StringSet{map[string]bool{}}
	for _, p := range strings.Split(protocols, ",") {
		s.Protocols.Add(p)
	}
	go GenerateURLs(s.Hosts, s.Ports, &s.Paths, c2)
	s.URLComponents = <-c2
	return
}

func Start(s State) {

	if s.Scan {
		fmt.Printf("Starting Port Scanner\n")
		if s.Debug {
			fmt.Printf("Testing %v host:port combinations\n", len(s.URLComponents))
		}
		fmt.Printf(LineSep())
		s.URLComponents = ScanHosts(&s)

		fmt.Printf(LineSep())
	}
	if s.Dirbust {
		fmt.Printf("Starting Dirbuster\n")
		if s.Debug {
			fmt.Printf("Testing %v URLs\n", len(s.URLComponents)*len(s.Paths.Set))
		}
		fmt.Printf(LineSep())

		s.URLComponents = DirbustHosts(&s)
		if s.Debug {
			fmt.Println(s.URLComponents)
		}
		for _, h := range s.URLComponents {
			fmt.Printf("%v: %v (%v) %v\n", h.HostAddr, h.Paths.Set, h.Port, h.Protocol)
		}
		fmt.Printf(LineSep())
	}
}
