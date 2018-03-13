package lib

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
)

// I want this function to do some intelligent analysis of GET requests. e.g. detect whether 200 is a soft-404, or whether 403 is just a directory access denied
// func SmartDetector(){

// }

// checks to see whether host is http/s or other scheme.
// Returns error if endpoint is not a valid webserver. Prevents
func prefetch(host Host, s *State) (h Host, err error) {
	var url string
	for scheme, _ := range s.Protocols.Set {
		url = fmt.Sprintf("%v://%v:%v", scheme, host.HostAddr, host.Port)
		client := &http.Client{
			Transport: tx}
		resp, err := client.Get(url)
		// resp.Body.Close()

		if err != nil {
			if strings.Contains(err.Error(), "http: server gave HTTP response to HTTPS client") {
				host.Protocol = "http" // we know it's a http port now
				return host, nil
			}
			continue
		} else if resp == nil {
			continue
		} else {
			host.Protocol = scheme
			return host, nil
		}
	}
	if err != nil {
		// We've tested all our schemes and it's still broken
		// probably not a http server?
		return Host{}, err
	}
	return host, nil
}

func DirbustHosts(s *State) (h []Host) {
	hostChan := make(chan Host, s.Threads)
	respChan := make(chan *http.Response, s.Threads)

	var wg sync.WaitGroup
	wg.Add(len(s.URLComponents) * len(s.Paths.Set))

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
	host, err := prefetch(host, s)
	// fmt.Println(host, err)
	if err != nil {
		wg.Add(len(host.Paths.Set) * -1) // Host is not going to be dirbusted - lets rm the ol' badboy
		return
	}
	if host.Protocol == "" {
		wg.Add(len(host.Paths.Set) * -1) // Host is not going to be dirbusted - lets rm the ol' badboy
		return
	}
	for path, _ := range host.Paths.Set {
		go HTTPGetter(s, host, path, hostChan, respChan, wg)
	}
}

func HTTPGetter(s *State, host Host, path string, hostChan chan Host, respChan chan *http.Response, wg *sync.WaitGroup) {
	defer wg.Done()
	if strings.HasPrefix(path, "/") {
		path = path[1:] // strip preceding '/' char
	}
	url := fmt.Sprintf("%v://%v:%v/%v", host.Protocol, host.HostAddr, host.Port, path)
	if s.Debug {
		fmt.Printf("Trying URL: %v\n", url)
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
