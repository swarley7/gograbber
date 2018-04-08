package lib

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/pmezard/go-difflib/difflib"
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
		resp, err := cl.Get(url)
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

func DirbustHosts(s *State, targets chan Host, results chan Host) {
	wg := sync.WaitGroup{}
	targetHost := make(TargetHost, s.Threads)
	var cnt int
	for host := range targets {
		var err error
		if !s.URLProvided {
			host, err = prefetch(host, s)
			if err != nil {
				continue
			}
			if host.Protocol == "" {
				continue
			}

		}
		if s.Soft404Detection {
			randURL := fmt.Sprintf("%v://%v:%v/%v", host.Protocol, host.HostAddr, host.Port, RandString(16))
			randResp, err := cl.Get(randURL)
			if err != nil {
				continue
				// panic(err)
			}
			data, err := ioutil.ReadAll(randResp.Body)
			if err != nil {
				// panic(err)
				continue
			}
			randResp.Body.Close()
			host.Soft404RandomURL = randURL
			host.Soft404RandomPageContents = strings.Split(string(data), " ")
		}
		for path := range host.Paths.Set {
			routineId := Counter{cnt}
			targetHost <- routineId
			wg.Add(1)
			go targetHost.HTTPGetter(host, s.Debug, s.Jitter, s.Soft404Detection, s.StatusCodesIgn, s.Ratio, path, results, &wg)
			cnt++
		}
	}
}

func (target TargetHost) HTTPGetter(host Host, debug bool, jitter int, soft404Detection bool, statusCodesIgn IntSet, Ratio float64, path string, hostChan chan<- Host, wg *sync.WaitGroup) {
	defer wg.Done()
	// debug
	if strings.HasPrefix(path, "/") && len(path) > 0 {
		path = path[1:] // strip preceding '/' char
	}
	url := fmt.Sprintf("%v://%v:%v/%v", host.Protocol, host.HostAddr, host.Port, path)
	if debug {
		fmt.Printf("Trying URL: %v\n", url)
	}
	if jitter > 0 {
		jitter := time.Duration(rand.Intn(jitter)) * time.Millisecond
		if debug {
			fmt.Printf("Jitter: %v\n", jitter)
		}
		time.Sleep(jitter)
	}
	// client := &http.Client{
	// 	Transport:     tx,
	// 	CheckRedirect: func(req *http.Request, via []*http.Request) error { return errors.New("something bad happened") },
	// }
	// client := &http.Client{
	// 	Transport: tx}
	resp, err := cl.Get(url)
	if err != nil {
		<-target
		return
	}
	defer resp.Body.Close()
	if statusCodesIgn.Contains(resp.StatusCode) {
		<-target
		return
	}
	if soft404Detection && path != "" {
		soft404Ratio := detectSoft404(resp, host.Soft404RandomPageContents)
		if soft404Ratio > Ratio {
			if debug {
				fmt.Printf("[%v] is very similar to [%v] (%.4f%% match)\n", url, host.Soft404RandomURL, (soft404Ratio * 100))
			}
			<-target
			return
		}
	}

	fmt.Printf("%v - %v\n", url, resp.StatusCode)
	// host.Protocols = StringSet{map[string]bool{}}
	// host.Protocols.Add(protocol)
	host.Paths = StringSet{map[string]bool{}}
	host.Paths.Add(path)
	hostChan <- host
	<-target
	return
}

// func HTTPHeader(host Host, debug bool, jitter int, soft404Detection bool, statusCodesIgn IntSet, Ratio float64, path string, hostChan chan<- Host, wg *sync.WaitGroup) {
// 	defer wg.Done()
// 	// debug
// 	if strings.HasPrefix(path, "/") {
// 		path = path[1:] // strip preceding '/' char
// 	}
// 	url := fmt.Sprintf("%v://%v:%v/%v", host.Protocol, host.HostAddr, host.Port, path)
// 	if debug {
// 		fmt.Printf("Trying URL: %v\n", url)
// 	}
// 	if jitter > 0 {
// 		jitter := time.Duration(rand.Intn(jitter)) * time.Millisecond
// 		if debug {
// 			fmt.Printf("Jitter: %v\n", jitter)
// 		}
// 		time.Sleep(jitter)
// 	}
// 	// client := &http.Client{
// 	// 	Transport:     tx,
// 	// 	CheckRedirect: func(req *http.Request, via []*http.Request) error { return errors.New("something bad happened") },
// 	// }
// 	// client := &http.Client{
// 	// 	Transport: tx}
// 	req, err := http.NewRequest("HEAD", url, nil) // Use head, to save data
// 	resp, err := cl.Do(req)
// 	if err != nil {
// 		return
// 	}
// 	defer resp.Body.Close()
// 	if statusCodesIgn.Contains(resp.StatusCode) {
// 		return
// 	}
// 	if soft404Detection && path != "/" {
// 		soft404Ratio := detectSoft404(resp, host.Soft404RandomPageContents)
// 		if soft404Ratio > Ratio {
// 			fmt.Printf("[%v] is very similar to [%v] (%.5f%% match)\n", url, host.Soft404RandomURL, (soft404Ratio * 100))
// 			return
// 		}

// 	}

// 	fmt.Printf("%v - %v\n", url, resp.StatusCode)
// 	// host.Protocols = StringSet{map[string]bool{}}
// 	// host.Protocols.Add(protocol)
// 	host.Paths = StringSet{map[string]bool{}}
// 	host.Paths.Add(path)
// 	hostChan <- host
// }
func detectSoft404(resp *http.Response, randRespData []string) (ratio float64) {
	// defer resp.Body.Close()
	diff := difflib.SequenceMatcher{}
	responseData, _ := ioutil.ReadAll(resp.Body)
	diff.SetSeqs(strings.Split(string(responseData), " "), randRespData)
	ratio = diff.Ratio()
	return ratio
}
