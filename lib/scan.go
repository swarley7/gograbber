package lib

import (
	"fmt"
	"math/rand"
	"net"
	"sync"
	"time"
)

// connectHost does the actual TCP connection
func ConnectHost(wg *sync.WaitGroup, timeout time.Duration, Jitter int, Debug bool, host Host, results chan Host, threads chan struct{}) {
	defer func() {
		<-threads
		wg.Done()
	}()
	if Debug {
		fmt.Printf("Port scanning: %v:%v\n", host.HostAddr, host.Port)
	}
	if Jitter > 0 {
		jitter := time.Duration(rand.Intn(Jitter)) * time.Millisecond
		if Debug {
			fmt.Printf("Jitter: %v\n", jitter)
		}
		time.Sleep(jitter)
	}
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%v:%v", host.HostAddr, host.Port), timeout)
	if err == nil {
		fmt.Printf("%v:%v OPEN\n", host.HostAddr, host.Port)
		conn.Close()
		results <- host
	}
}
