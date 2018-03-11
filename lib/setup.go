package lib

import (
	"fmt"
	"math"
	"strings"

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

	}

	return
}

func Start(s State) {
	if s.Scan {
		openPorts := ScanHosts(&s)
		for socketPair := range openPorts.Set {
			fmt.Printf("Host:Port %s is open\n", socketPair)
		}
	}
}
