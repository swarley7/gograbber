package lib

import (
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/benbjohnson/phantomjs"
)

func Screenshot(s *State) (h []Host) {
	// Start the process once.
	// Open a URL.
	hostChan := make(chan Host, s.Threads)
	respChan := make(chan *http.Response, s.Threads)
	if err := phantomjs.DefaultProcess.Open(); err != nil {
		panic(err)
	}
	defer phantomjs.DefaultProcess.Close()
	var wg sync.WaitGroup
	for _, URLComponent := range s.URLComponents {
		go distributeScreenshotWorkers(s, URLComponent, hostChan, respChan, &wg)
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

func distributeScreenshotWorkers(s *State, host Host, hostChan chan Host, respChan chan *http.Response, wg *sync.WaitGroup) {
	if !s.URLProvided {
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
	}
	for path := range host.Paths.Set {
		go ScreenshotAURL(s, host, path, hostChan, respChan, wg)
	}
}

func ScreenshotAURL(s *State, host Host, path string, hostChan chan Host, respChan chan *http.Response, wg *sync.WaitGroup) (err error) {
	p := phantomjs.NewProcess()
	page, err := p.CreateWebPage()
	if err != nil {
		return err
	}
	defer wg.Done()
	if strings.HasPrefix(path, "/") {
		path = path[1:] // strip preceding '/' char
	}
	url := fmt.Sprintf("%v://%v:%v/%v", host.Protocol, host.HostAddr, host.Port, path)
	if s.Debug {
		fmt.Printf("Trying to screenshot URL: %v\n", url)
	}
	if s.Jitter > 0 {
		jitter := time.Duration(rand.Intn(s.Jitter)) * time.Millisecond
		fmt.Printf("Jitter: %v\n", jitter)
		time.Sleep(jitter)
	}
	defer page.Close()
	if err := page.Open(url); err != nil {
		return err
	}

	// Setup the viewport and render the results view.
	if err := page.SetViewportSize(s.ImgX, s.ImgY); err != nil {
		return err
	}
	screenshotFilename := fmt.Sprintf("%v/%v_%v_%v_%v.png", s.OutputDirectory, host.Protocol, host.HostAddr, host.Port, path)
	if err := page.Render(screenshotFilename, "png", 100); err != nil {
		return err
	}
	return
}
