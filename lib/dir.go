package lib

import (
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"
)

// I want this function to do some intelligent analysis of GET requests. e.g. detect whether 200 is a soft-404, or whether 403 is just a directory access denied
// func SmartDetector(){

// }

// checks to see whether host is http/s or other scheme.
// Returns error if endpoint is not a valid webserver. Prevents
func prefetch(host Host, s *State) (h Host, err error) {
	var url string
	for scheme := range s.Protocols.Set {
		if s.Jitter > 0 {
			jitter := time.Duration(rand.Intn(s.Jitter)) * time.Millisecond
			if s.Debug {
				fmt.Printf("Jitter: %v\n", jitter)
			}
			time.Sleep(jitter)
		}
		url = fmt.Sprintf("%v://%v:%v", scheme, host.HostAddr, host.Port)
		if s.Debug {
			fmt.Printf("Prefetch URL: %v\n", url)
		}
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
			resp.Body.Close()
			continue
		} else {
			host.Protocol = scheme
			resp.Body.Close()
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

	wg := sync.WaitGroup{}
	// go func() {
	// 	x := time.NewTicker(time.Second * 1)
	// 	for _ = range x.C {
	// 		fmt.Println(wg)
	// 	}
	// }()
	for _, URLComponent := range s.URLComponents {
		wg.Add(1)
		go distributeHTTPRequests(s, URLComponent, hostChan, respChan, &wg)
	}

	go func() {
		for url := range hostChan {
			h = append(h, url)
		}
	}()
	wg.Wait()
	close(hostChan)
	close(respChan)
	// write resps to file? return hosts for now
	return h
}

func distributeHTTPRequests(s *State, host Host, hostChan chan Host, respChan chan *http.Response, wg *sync.WaitGroup) {
	defer wg.Done()
	var err error
	if !s.URLProvided {
		host, err = prefetch(host, s)
		// fmt.Println(host, err)
		if err != nil {
			return
		}
		if host.Protocol == "" {
			return
		}

	}
	for path := range host.Paths.Set {
		wg.Add(1)
		go HTTPGetter(s, host, path, hostChan, respChan, wg)
	}
}

func HTTPGetter(s *State, host Host, path string, hostChan chan Host, respChan chan *http.Response, wg *sync.WaitGroup) {
	defer wg.Done()
	// debug
	if strings.HasPrefix(path, "/") {
		path = path[1:] // strip preceding '/' char
	}
	url := fmt.Sprintf("%v://%v:%v/%v", host.Protocol, host.HostAddr, host.Port, path)
	if s.Debug {
		fmt.Printf("Trying URL: %v\n", url)
	}
	if s.Jitter > 0 {
		jitter := time.Duration(rand.Intn(s.Jitter)) * time.Millisecond
		if s.Debug {
			fmt.Printf("Jitter: %v\n", jitter)
		}
		time.Sleep(jitter)
	}
	// client := &http.Client{
	// 	Transport: tx,
	// 	// CheckRedirect: func(req *http.Request, via []*http.Request) error { return errors.New("something bad happened") },
	// }
	// client := &http.Client{
	// 	Transport: tx}
	client := http.Client{Timeout: time.Duration(5 * time.Second)}
	resp, err := client.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if s.StatusCodesIgn.Contains(resp.StatusCode) {
		return
	}
	fmt.Printf("%v - %v\n", url, resp.StatusCode)
	// host.Protocols = StringSet{map[string]bool{}}
	// host.Protocols.Add(protocol)
	host.Paths = StringSet{map[string]bool{}}
	host.Paths.Add(path)
	respChan <- resp
	hostChan <- host
}
