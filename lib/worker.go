package lib

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"time"
)

// func statusUpdater() {
// 	//update output every 3 seconds or so
// 	tick := time.Tick(time.Second * 3)
// }

func writerWorker(writeChan chan []byte, filename string) {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if os.IsNotExist(err) {
		file, err = os.Create(filename)
	}
	if err != nil {
		panic(err)
	}
	writer := bufio.NewWriter(file)
	for {
		b := <-writeChan
		if len(b) > 0 {
			writer.Write(b)
			writer.Flush()
		}
	}
}

func RoutineManager(s *State, ScanChan chan Host, DirbustChan chan Host, ScreenshotChan chan Host, wg *sync.WaitGroup) {
	defer wg.Done()
	threadChan := make(chan struct{}, s.Threads)
	var err error
	var scanWg = sync.WaitGroup{}
	var dirbWg = sync.WaitGroup{}
	var screenshotWg = sync.WaitGroup{}

	ticker := time.NewTicker(10 * time.Second)
	go func() {
		for t := range ticker.C {
			fmt.Printf("Tick at %v\n", t)
		}
	}()
	wg.Add(1)
	go func() {
		defer func() {
			close(ScanChan)
			wg.Done()
		}()
		if !s.Scan {
			for host := range s.Targets {
				ScanChan <- host
			}
			return
		}
		for host := range s.Targets {
			scanWg.Add(1)
			threadChan <- struct{}{}
			go ConnectHost(&scanWg, s.Timeout*time.Second, s.Jitter, s.Debug, host, ScanChan, threadChan)
		}
		scanWg.Wait()
		return
	}()

	wg.Add(1)
	go func() {
		defer func() {
			close(DirbustChan)
			wg.Done()
		}()

		if !s.Dirbust {
			for host := range ScanChan {
				if !s.URLProvided {
					for scheme := range s.Protocols.Set {
						host.Protocol = scheme
						DirbustChan <- host
					}

				} else {
					DirbustChan <- host
				}
			}
			return
		}
		var fuggoff bool
		// Do dirbusting
		for host := range ScanChan {
			fuggoff = false
			if !s.URLProvided && !host.PrefetchDoneCheck(s.PrefetchedHosts) {
				host, err = Prefetch(host, s)
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
				// fmt.Printf("Soft404 checking [%v]\n", randURL)
				randResp, err := cl.Get(randURL)
				if err != nil {
					fuggoff = true
					continue
					// panic(err)
				}
				data, err := ioutil.ReadAll(randResp.Body)
				if err != nil {
					// panic(err)
					fuggoff = true
					continue
				}
				randResp.Body.Close()
				host.Soft404RandomURL = randURL
				host.Soft404RandomPageContents = strings.Split(string(data), " ")
				s.Soft404edHosts[host.Soft404Hash()] = true
			}
			if fuggoff {
				continue
			}
			if !s.URLProvided {
				for path, _ := range s.Paths.Set {
					// fmt.Printf("HTTP GET to [%v://%v:%v/%v]\n", host.Protocol, host.HostAddr, host.Port, host.Path)
					threadChan <- struct{}{}
					dirbWg.Add(1)
					go HTTPGetter(&dirbWg, host, s.Debug, s.Jitter, s.Soft404Detection, s.StatusCodesIgn, s.Ratio, path, DirbustChan, threadChan, s.ProjectName, s.HTTPResponseDirectory)
				}
			} else {
				threadChan <- struct{}{}
				dirbWg.Add(1)
				go HTTPGetter(&dirbWg, host, s.Debug, s.Jitter, s.Soft404Detection, s.StatusCodesIgn, s.Ratio, host.Path, DirbustChan, threadChan, s.ProjectName, s.HTTPResponseDirectory)
			}
		}
		dirbWg.Wait()
	}()

	wg.Add(1)
	go func() {
		defer func() {
			close(ScreenshotChan)
			wg.Done()
		}()

		if !s.Screenshot {
			for host := range DirbustChan {
				ScreenshotChan <- host
			}
			return
		}
		var cnt int
		for host := range DirbustChan {
			threadChan <- struct{}{}
			screenshotWg.Add(1)
			go ScreenshotAURL(&screenshotWg, s, cnt, host, ScreenshotChan, threadChan)
			cnt++
		}
		screenshotWg.Wait()
		return
	}()

	// go func() {
	// 	for {
	// 		select {
	// 		case <-targetHost:
	// 			{
	// 				continue
	// 			}

	// 		}
	// 	}
	// }()
}
