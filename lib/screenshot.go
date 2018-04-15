package lib

import (
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/benbjohnson/phantomjs"
)

func ScreenshotAURL(wg *sync.WaitGroup, s *State, cnt int, host Host, results chan Host, threads chan struct{}) (err error) {
	defer func() {
		<-threads

		wg.Done()
	}()
	page, err := s.PhantomProcesses[cnt%len(s.PhantomProcesses)].CreateWebPage()
	url := fmt.Sprintf("%v://%v:%v/%v", host.Protocol, host.HostAddr, host.Port, host.Path)

	if err != nil {
		fmt.Printf("Unable to Create webpage: %v (%v)\n", url, err)
		return err
	}
	defer page.Close()
	page.SetSettings(phantomjs.WebPageSettings{ResourceTimeout: s.Timeout}) // Time out the page if it takes too long to load. Sometimes JS is fucky and takes wicked long to do nothing forever :(

	if strings.HasPrefix(host.Path, "/") {
		host.Path = host.Path[1:] // strip preceding '/' char
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
		// <-target
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
	results <- host
	return err
}
