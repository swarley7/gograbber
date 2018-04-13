package lib

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// connectHost does the actual TCP connection
func (target TargetHost) ConnectHost(wg *sync.WaitGroup, Jitter int, Debug bool, host Host, results chan Host) {
	defer func() {
		<-target
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
	conn, err := d.Dial("tcp", fmt.Sprintf("%v:%v", host.HostAddr, host.Port))
	if err == nil {
		fmt.Printf("%v:%v OPEN\n", host.HostAddr, host.Port)
		conn.Close()
		results <- host
	}
}
