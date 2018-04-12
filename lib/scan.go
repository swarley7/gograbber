package lib

import (
	"fmt"
	"math/rand"
	"time"
)

// var finished = false

// // ScanHosts performs a TCP Portscan of hosts. Currently uses complete handshake. May look into SYN scan later.
// func ScanHosts(s *State, targets chan Host, results chan Host, wgExt *sync.WaitGroup) {
// 	// defer close(results)
// 	defer wgExt.Done()
// 	var wg sync.WaitGroup
// 	targetHost := make(TargetHost, s.Threads)
// 	var cnt int
// 	for {
// 		select {
// 		case urlComponent := <-targets:
// 			{
// 				wg.Add(1)
// 				routineId := Counter{cnt}
// 				targetHost <- routineId
// 				go targetHost.ConnectHost(s, urlComponent, results)
// 				cnt++
// 			}
// 		}
// 	}
// 	wg.Wait()

// }

// connectHost does the actual TCP connection
func (target TargetHost) ConnectHost(s *State, host Host, results chan Host) {
	defer func() {
		<-target
	}()
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
		results <- host
	}
}
