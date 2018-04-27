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

func HTTPGetter(wg *sync.WaitGroup, host Host, debug bool, jitter int, soft404Detection bool, statusCodesIgn IntSet, Ratio float64, path string, results chan Host, threads chan struct{}) {
	defer func() {
		<-threads
		wg.Done()
	}()
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
	var err error
	host.HTTPReq, err = http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}
	host.HTTPResp, err = cl.Do(host.HTTPReq)
	if err != nil {
		return
	}
	defer host.HTTPResp.Body.Close()
	if statusCodesIgn.Contains(host.HTTPResp.StatusCode) {
		return
	}
	if soft404Detection && path != "" {
		soft404Ratio := detectSoft404(host.HTTPResp, host.Soft404RandomPageContents)
		if soft404Ratio > Ratio {
			if debug {
				fmt.Printf("[%v] is very similar to [%v] (%.4f%% match)\n", url, host.Soft404RandomURL, (soft404Ratio * 100))
			}
			return
		}
	}

	fmt.Printf("%v - %v\n", url, host.HTTPResp.StatusCode)
	host.Path = path
	results <- host
}

func detectSoft404(resp *http.Response, randRespData []string) (ratio float64) {
	// defer resp.Body.Close()
	diff := difflib.SequenceMatcher{}
	responseData, _ := ioutil.ReadAll(resp.Body)
	diff.SetSeqs(strings.Split(string(responseData), " "), randRespData)
	ratio = diff.Ratio()
	return ratio
}
