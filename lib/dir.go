package lib

import (
	"errors"
	"fmt"
	"net/http"
	"sync"
)

// I want this function to do some intelligent analysis of GET requests. e.g. detect whether 200 is a soft-404, or whether 403 is just a directory access denied
// func SmartDetector(){

// }

func DirbustHosts(s *State) (h []Host) {
	hostChan := make(chan Host, s.Threads)
	respChan := make(chan *http.Response, s.Threads)

	var wg sync.WaitGroup
	wg.Add(len(s.URLComponents) * len(s.Paths.Set) * len(s.Protocols.Set))

	for _, URLComponent := range s.URLComponents {
		go distributeHTTPRequests(s, URLComponent, hostChan, respChan, &wg)
	}
	wg.Wait()
	close(hostChan)
	close(respChan)

	for url := range hostChan {
		h = append(h, url)
	}
	// write resps to file? return hosts for now
	return h
}

func distributeHTTPRequests(s *State, host Host, hostChan chan Host, respChan chan *http.Response, wg *sync.WaitGroup) {
	for path, _ := range host.Paths.Set {
		for protocol, _ := range host.Protocols.Set {
			go HTTPGetter(s, host, protocol, path, hostChan, respChan, wg)
		}
	}
}

func HTTPGetter(s *State, host Host, protocol string, path string, hostChan chan Host, respChan chan *http.Response, wg *sync.WaitGroup) {
	defer wg.Done()
	url := fmt.Sprintf("%v://%v:%v/%v", protocol, host.HostAddr, host.Port, path)
	if s.Debug {
		fmt.Printf("Trying URL: %v://%v:%v/%v\n", protocol, host.HostAddr, host.Port, path)
	}
	client := &http.Client{
		Transport:     tx,
		CheckRedirect: func(req *http.Request, via []*http.Request) error { return errors.New("something bad happened") },
	}
	resp, err := client.Get(url)
	if err != nil {
		return
	}
	if s.StatusCodesIgn.Contains(resp.StatusCode) {
		return
	}

	// host.Protocols = StringSet{map[string]bool{}}
	// host.Protocols.Add(protocol)
	respChan <- resp
	hostChan <- host
	resp.Body.Close()

}
