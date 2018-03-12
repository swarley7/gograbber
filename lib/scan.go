package lib

import (
	"fmt"
	"net"
	"sync"
	"time"
)

// ScanHosts performs a TCP Portscan of hosts. Currently uses complete handshake. May look into SYN scan later.
func ScanHosts(s *State) (h []Host) {
	ch := make(chan Host, s.Threads)
	var wg sync.WaitGroup
	wg.Add(len(s.URLComponents))

	for _, URLComponent := range s.URLComponents {
		go connectHost(s, URLComponent, ch, &wg)
	}
	wg.Wait()
	close(ch)

	for AliveHost := range ch {
		h = append(h, AliveHost)
	}
	return h
}

// connectHost does the actual TCP connection
func connectHost(s *State, host Host, ch chan Host, wg *sync.WaitGroup) {
	defer wg.Done()
	if s.Debug {
		fmt.Printf("Port scanning: %v:%v\n", host.HostAddr, host.Port)
	}
	d := net.Dialer{Timeout: 5 * time.Second}
	conn, err := d.Dial("tcp", fmt.Sprintf("%v:%v", host.HostAddr, host.Port))
	if err == nil {
		fmt.Printf("%v:%v OPEN\n", host.HostAddr, host.Port)
		ch <- host
		conn.Close()
	}
}
