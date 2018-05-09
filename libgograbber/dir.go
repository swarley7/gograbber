package libgograbber

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/pmezard/go-difflib/difflib"
)

// checks to see whether host is http/s or other scheme.
// Returns error if endpoint is not a valid webserver. Prevents
func Prefetch(host Host, debug bool, jitter int, protocols StringSet) (h Host, err error) {
	var url string
	for scheme := range protocols.Set {
		ApplyJitter(jitter)
		url = fmt.Sprintf("%v://%v:%v", scheme, host.HostAddr, host.Port)
		if debug {
			Debug.Printf("Prefetch URL: %v\n", url)
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

func HTTPGetter(wg *sync.WaitGroup, host Host, debug bool, Jitter int, soft404Detection bool, statusCodesIgn IntSet, Ratio float64, path string, results chan Host, threads chan struct{}, ProjectName string, responseDirectory string, writeChan chan []byte) {
	defer func() {
		<-threads
		wg.Done()
	}()

	if strings.HasPrefix(path, "/") && len(path) > 0 {
		path = path[1:] // strip preceding '/' char
	}
	url := fmt.Sprintf("%v://%v:%v/%v", host.Protocol, host.HostAddr, host.Port, path)
	if debug {
		Debug.Printf("Trying URL: %v\n", url)
	}
	ApplyJitter(Jitter)

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
				Debug.Printf("[%v] is very similar to [%v] (%v match)\n", y.Sprintf("%s", url), y.Sprintf("%s", host.Soft404RandomURL), y.Sprintf("%.4f%%", (soft404Ratio*100)))
			}
			return
		}
	}

	Good.Printf("%v - %v\n", url, g.Sprintf("%d", host.HTTPResp.StatusCode))
	t := time.Now()
	currTime := fmt.Sprintf("%d%d%d%d%d%d", t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())
	var responseFilename string
	if ProjectName != "" {
		responseFilename = fmt.Sprintf("%v/%v_%v_%v_%v-%v_%v.html", responseDirectory, strings.ToLower(strings.Replace(ProjectName, " ", "_", -1)), host.Protocol, host.HostAddr, host.Port, currTime, rand.Int63())
	} else {
		responseFilename = fmt.Sprintf("%v/%v_%v_%v-%v_%v.html", responseDirectory, host.Protocol, host.HostAddr, host.Port, currTime, rand.Int63())
	}
	file, err := os.Create(responseFilename)
	if err != nil {
		Error.Printf("%v\n", err)
	}
	buf, err := ioutil.ReadAll(host.HTTPResp.Body)
	if err != nil {
		Error.Printf("%v\n", err)
	} else {
		if len(buf) > 0 {
			file.Write(buf)
			host.ResponseBodyFilename = responseFilename
		} else {
			_ = os.Remove(responseFilename)
		}
	}
	host.Path = path
	// writeChan <- []byte(fmt.Sprintf("%v\n", url))
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
