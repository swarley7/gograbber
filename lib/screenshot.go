package lib

import (
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/benbjohnson/phantomjs"
)

// func Screenshot(s *State, targets chan Host, results chan Host, wgExt *sync.WaitGroup) {
// 	// defer close(results)
// 	defer wgExt.Done()
// 	wg := sync.WaitGroup{}
// 	targetHost := make(TargetHost, s.Threads)
// 	var cnt int
// 	for {
// 		select {
// 		case host := <-targets:
// 			{
// 				for path := range host.Paths.Set {
// 					wg.Add(1) //MAKE SURE SCREENSHOTURL HAS A DONE CALL IN IT JFC
// 					routineId := Counter{cnt}
// 					targetHost <- routineId
// 					fmt.Printf("Screenshotting: [%v]\n", cnt)
// 					go targetHost.ScreenshotAURL(s, cnt, host, path, results, &wg)
// 					cnt++
// 				}
// 			}
// 		} // wg.Add(1)
// 		// go distributeScreenshotWorkers(s, URLComponent, hostChan, respChan, &wg)
// 	}
// 	wg.Wait()
// }

func (target TargetHost) ScreenshotHost(s *State, host Host, results chan Host) {
	// defer close(results)
	wg := sync.WaitGroup{}
	// targetHost := make(TargetHost, s.Threads)
	// var cnt int
	wg.Wait()
}

func (target TargetHost) ScreenshotAURL(s *State, cnt int, host Host, path string, hostChan chan Host, wg *sync.WaitGroup) (err error) {
	defer func() {
		<-target
		wg.Done()
	}()
	page, err := s.PhantomProcesses[cnt%len(s.PhantomProcesses)].CreateWebPage()
	url := fmt.Sprintf("%v://%v:%v/%v", host.Protocol, host.HostAddr, host.Port, path)

	if err != nil {
		fmt.Printf("Unable to Create webpage: %v (%v)\n", url, err)
		return err
	}
	defer page.Close()
	page.SetSettings(phantomjs.WebPageSettings{ResourceTimeout: s.Timeout}) // Time out the page if it takes too long to load. Sometimes JS is fucky and takes wicked long to do nothing forever :(

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
		<-target
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
	return err
}
