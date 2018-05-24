package libgograbber

import (
	"fmt"
	"net"
	"sync"
	"time"
)

// connectHost does the actual TCP connection
func ConnectHost(wg *sync.WaitGroup, timeout time.Duration, Jitter int, debug bool, host Host, results chan Host, threads chan struct{}, writeChan chan []byte) {
	defer func() {
		<-threads
		wg.Done()
	}()
	if debug {
		Info.Printf("Port scanning: %v:%v\n", host.HostAddr, host.Port)
	}
	ApplyJitter(Jitter)
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%v:%v", host.HostAddr, host.Port), timeout)
	if err == nil {
		defer conn.Close()
		Good.Printf("%v:%v - %v\n", host.HostAddr, host.Port, g.Sprintf("tcp/%v open", host.Port))
		writeChan <- []byte(fmt.Sprintf("%v,%v\n", host.HostAddr, host.Port))
		results <- host
	} else {
		if debug {
			Debug.Printf("Err: %v\n", err)
		}
	}
}
