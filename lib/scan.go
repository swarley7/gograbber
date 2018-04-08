package lib

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// ScanHosts performs a TCP Portscan of hosts. Currently uses complete handshake. May look into SYN scan later.
func ScanHosts(s *State, targets chan Host, results chan Host) {
	var wg sync.WaitGroup
	targetHost := make(TargetHost, s.Threads)
	var cnt int
	for urlComponent := range targets {
		wg.Add(1)
		routineId := Counter{cnt}
		targetHost <- routineId
		go targetHost.connectHost(s, urlComponent, results, &wg)
		cnt++
	}
}

// connectHost does the actual TCP connection
func (target TargetHost) connectHost(s *State, host Host, ch chan Host, wg *sync.WaitGroup) {
	defer wg.Done()
	if s.Debug {
		fmt.Printf("Port scanning: %v:%v\n", host.HostAddr, host.Port)
	}
	if s.Jitter > 0 {
		jitter := time.Duration(rand.Intn(s.Jitter)) * time.Millisecond
		if s.Debug {
			fmt.Printf("Jitter: %v\n", jitter)
		}
		time.Sleep(jitter)
	}
	conn, err := d.Dial("tcp", fmt.Sprintf("%v:%v", host.HostAddr, host.Port))
	if err == nil {
		fmt.Printf("%v:%v OPEN\n", host.HostAddr, host.Port)
		conn.Close()
		ch <- host
	}
	<-target
}
