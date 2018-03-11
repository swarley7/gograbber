package lib

import (
	"fmt"
	"math"
	"strings"
	"sync"
	"time"

	multierror "github.com/hashicorp/go-multierror"
)

func Initialise(s *State, ports string) (errors *multierror.Error) {
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
	c2 := make(chan []Host)
	go GenerateURLs(s.Hosts, s.Protocol, s.Ports, s.Paths, c2)
	s.URLComponents = <-c2
	return
}

func Start(s State) {
	workers := s.Threads

	indexChan := make(chan string, 1000)
	configChan := make(chan string, 1000)
	writeChan := make(chan []byte, 1000)

	go writerWorker(writeChan, s.OutputFile)

	if s.Scan {
		fmt.Printf("Starting Port Scanner\n")
		if s.Debug {
			fmt.Printf("Testing %v host:port combinations\n", len(s.URLComponents))
		}
		fmt.Printf(LineSep())
		// AliveHosts := []Host{}

		s.URLComponents = ScanHosts(&s)
		// for AliveHost := range hostChan {
		// 	AliveHosts = append(AliveHosts, AliveHost)
		// }
		// s.URLComponents = AliveHosts
		for _, URLComponent := range s.URLComponents {
			fmt.Printf("%v:%v is alive\n", URLComponent.HostAddr, URLComponent.Port)
		}
	}
	// go statusUpdater()
	wg := sync.WaitGroup{}

	if false {
		wg.Add(1)
		finput := make(chan struct{})
		go routineManager(finput, workers, indexChan, configChan, writeChan, &wg)
		/*for x := 0; x < workers; x++ {
			go taskWorker(indexChan, configChan, writeChan, stopChan, stopped)
		}*/

		close(finput)
		filled = true
		wg.Wait()
		time.Sleep(time.Second * 5)
	}
}
