package lib

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/pmezard/go-difflib/difflib"
)

// I want this function to do some intelligent analysis of GET requests. e.g. detect whether 200 is a soft-404, or whether 403 is just a directory access denied
// func SmartDetector(){

// }

// checks to see whether host is http/s or other scheme.
// Returns error if endpoint is not a valid webserver. Prevents
func Prefetch(host Host, s *State) (h Host, err error) {
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

// func DirbustHosts(s *State, targets chan Host, results chan Host, wgExt *sync.WaitGroup) {
// 	// defer close(results)
// 	defer wgExt.Done()
// 	wg := sync.WaitGroup{}
// 	targetHost := make(TargetHost, s.Threads)
// 	// endChan := make(chan struct{})
// 	var cnt int
// 	for {
// 		select {
// 		case host := <-targets:
// 			{
// 				var err error
// 				if !s.URLProvided {
// 					host, err = prefetch(host, s)
// 					if err != nil {
// 						continue
// 					}
// 					if host.Protocol == "" {
// 						continue
// 					}

// 				}
// 				if s.Soft404Detection {
// 					randURL := fmt.Sprintf("%v://%v:%v/%v", host.Protocol, host.HostAddr, host.Port, RandString(16))
// 					randResp, err := cl.Get(randURL)
// 					if err != nil {
// 						continue
// 						// panic(err)
// 					}
// 					data, err := ioutil.ReadAll(randResp.Body)
// 					if err != nil {
// 						// panic(err)
// 						continue
// 					}
// 					randResp.Body.Close()
// 					host.Soft404RandomURL = randURL
// 					host.Soft404RandomPageContents = strings.Split(string(data), " ")
// 				}

// 				for path := range host.Paths.Set {
// 					routineId := Counter{cnt}
// 					targetHost <- routineId
// 					wg.Add(1)
// 					go targetHost.HTTPGetter(host, s.Debug, s.Jitter, s.Soft404Detection, s.StatusCodesIgn, s.Ratio, path, results, &wg)
// 					cnt++
// 				}
// 			}
// 		}
// 		if targets == nil {
// 			return
// 		}
// 	}
// 	wg.Wait()
// }

func (target TargetHost) DirbustHost(s *State, host Host, path string, results chan Host) {
	defer func() {
		<-target
	}()
	var fuggoff bool = false
	if !s.URLProvided && !host.PrefetchDoneCheck(s.PrefetchedHosts) {
		host, err := Prefetch(host, s)
		if err != nil {
			fuggoff = true
		}
		if host.Protocol == "" {
			fuggoff = true
		}
		s.PrefetchedHosts[host.PrefetchHash()] = true
	}
	if s.Soft404Detection && !host.Soft404DoneCheck(s.Soft404edHosts) {
		randURL := fmt.Sprintf("%v://%v:%v/%v", host.Protocol, host.HostAddr, host.Port, RandString(16))
		randResp, err := cl.Get(randURL)
		if err != nil {
			fuggoff = true
			// panic(err)
		}
		data, err := ioutil.ReadAll(randResp.Body)
		if err != nil {
			// panic(err)
			fuggoff = true
		}
		randResp.Body.Close()
		host.Soft404RandomURL = randURL
		host.Soft404RandomPageContents = strings.Split(string(data), " ")
		s.Soft404edHosts[host.Soft404Hash()] = true
	}
	if !fuggoff {
		HTTPGetter(host, s.Debug, s.Jitter, s.Soft404Detection, s.StatusCodesIgn, s.Ratio, path, results)
	}
	return
}

func HTTPGetter(host Host, debug bool, jitter int, soft404Detection bool, statusCodesIgn IntSet, Ratio float64, path string, hostChan chan Host) {
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

	resp, err := cl.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if statusCodesIgn.Contains(resp.StatusCode) {
		return
	}
	if soft404Detection && path != "" {
		soft404Ratio := detectSoft404(resp, host.Soft404RandomPageContents)
		if soft404Ratio > Ratio {
			if debug {
				fmt.Printf("[%v] is very similar to [%v] (%.4f%% match)\n", url, host.Soft404RandomURL, (soft404Ratio * 100))
			}
			return
		}
	}

	fmt.Printf("%v - %v\n", url, resp.StatusCode)
	// host.Protocols = StringSet{map[string]bool{}}
	// host.Protocols.Add(protocol)
	host.Path = path
	hostChan <- host
}

func detectSoft404(resp *http.Response, randRespData []string) (ratio float64) {
	// defer resp.Body.Close()
	diff := difflib.SequenceMatcher{}
	responseData, _ := ioutil.ReadAll(resp.Body)
	diff.SetSeqs(strings.Split(string(responseData), " "), randRespData)
	ratio = diff.Ratio()
	return ratio
}
