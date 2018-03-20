package lib

import (
	"fmt"
	"math/rand"
	"net"
	"sync"
	"time"
)

// ScanHosts performs a TCP Portscan of hosts. Currently uses complete handshake. May look into SYN scan later.
func ScanHosts(s *State) (h []Host) {
	ch := make(chan Host, s.Threads)
	var wg sync.WaitGroup

	for _, URLComponent := range s.URLComponents {
		wg.Add(1)
		go connectHost(s, URLComponent, ch, &wg)
	}

	go func() {
		for AliveHost := range ch {
			h = append(h, AliveHost)
		}
	}()
	wg.Wait()
	close(ch)
	return h
}

// connectHost does the actual TCP connection
func connectHost(s *State, host Host, ch chan Host, wg *sync.WaitGroup) {
	defer wg.Done()
	if s.Debug {
		fmt.Printf("Port scanning: %v:%v\n", host.HostAddr, host.Port)
	}
	if s.Jitter > 0 {
		jitter := time.Duration(rand.Intn(s.Jitter)) * time.Millisecond
		fmt.Printf("Jitter: %v\n", jitter)
		time.Sleep(jitter)
	}
	d := net.Dialer{Timeout: time.Duration(3 * time.Second)}
	conn, err := d.Dial("tcp", fmt.Sprintf("%v:%v", host.HostAddr, host.Port))
	if err == nil {
		fmt.Printf("%v:%v OPEN\n", host.HostAddr, host.Port)
		ch <- host
		conn.Close()
	}
}
