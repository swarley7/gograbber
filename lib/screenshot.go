package lib

import (
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"
)

func Screenshot(s *State) (h []Host) {
	for true {
		page, err := s.PhantomProcess.CreateWebPage()
		if err != nil {
			fmt.Println(err)
			time.Sleep(time.Second)
			page.Close()
			continue
		}
		if err := page.Open("http://localhost:20202/"); err != nil {
			fmt.Println(err)
			time.Sleep(time.Second)
			page.Close()
			continue
		}
		page.Close()
		break
	}
	hostChan := make(chan Host, s.Threads)
	respChan := make(chan *http.Response, s.Threads)

	wg := sync.WaitGroup{}
	for _, URLComponent := range s.URLComponents {
		wg.Add(1)
		go distributeScreenshotWorkers(s, URLComponent, hostChan, respChan, &wg)
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

func distributeScreenshotWorkers(s *State, host Host, hostChan chan Host, respChan chan *http.Response, wg *sync.WaitGroup) {
	//wg.Add called before this, so we FUCKING DEFER DONE IT
	defer wg.Done()
	for path := range host.Paths.Set {
		wg.Add(1) //MAKE SURE SCREENSHOTURL HAS A DONE CALL IN IT JFC
		go ScreenshotAURL(s, host, path, hostChan, respChan, wg)
	}
}

func ScreenshotAURL(s *State, host Host, path string, hostChan chan Host, respChan chan *http.Response, wg *sync.WaitGroup) (err error) {
	defer wg.Done()
	page, err := s.PhantomProcess.CreateWebPage()
	url := fmt.Sprintf("%v://%v:%v/%v", host.Protocol, host.HostAddr, host.Port, path)

	if err != nil {
		fmt.Printf("Unable to Create webpage: %v (%v)\n", url, err)
		return err
	}
	defer page.Close()

	if strings.HasPrefix(path, "/") {
		path = path[1:] // strip preceding '/' char
	}
	if s.Debug {
		fmt.Printf("Trying to screenshot URL: %v\n", url)
	}
	if s.Jitter > 0 {
		jitter := time.Duration(rand.Intn(s.Jitter)) * time.Millisecond
		if s.Debug {
			fmt.Printf("Jitter: %v\n", jitter)
		}
		time.Sleep(jitter)
	}
	if err := page.Open(url); err != nil {
		fmt.Printf("Unable to open page: %v (%v)\n", url, err)
		return err
	}

	// Setup the viewport and render the results view.
	if err := page.SetViewportSize(s.ImgX, s.ImgY); err != nil {
		fmt.Printf("Unable to set Viewport size: %v (%v)\n", url, err)
		return err
	}
	currTime := strings.Replace(time.Now().Format(time.RFC3339), ":", "_", -1)
	var screenshotFilename string
	if s.ProjectName != "" {
		screenshotFilename = fmt.Sprintf("%v/%v_%v_%v_%v-%v_%v.png", s.ScreenshotDirectory, strings.ToLower(strings.Replace(s.ProjectName, " ", "_", -1)), host.Protocol, host.HostAddr, host.Port, currTime, rand.Int63())
	} else {
		screenshotFilename = fmt.Sprintf("%v/%v_%v_%v-%v_%v.png", s.ScreenshotDirectory, host.Protocol, host.HostAddr, host.Port, currTime, rand.Int63())
	}
	fmt.Println(screenshotFilename)
	if err := page.Render(screenshotFilename, "png", s.ScreenshotQuality); err != nil {
		fmt.Printf("Unable to save Screenshot: %v (%v)\n", url, err)
		return err
	}
	host.ScreenshotFilename = screenshotFilename
	hostChan <- host
	return
}
