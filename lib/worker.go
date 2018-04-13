package lib

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"sync"
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
	targetHost := make(TargetHost, s.Threads)
	var err error
	var scanWg = sync.WaitGroup{}
	var dirbWg = sync.WaitGroup{}
	var screenshotWg = sync.WaitGroup{}
	var firstRunS bool = true
	var firstRunD bool = true
	var firstRunSS bool = true
	var doneScan bool = false
	var doneDirbust bool = false
	var doneScreenshot bool = false

	var cnt int
	for {
		select {
		case host, ok := <-s.Targets:
			if ok {
				if !s.Scan {
					// We're not supposed to scan, so let's pump it into the output chan!
					ScanChan <- host
					break
				}
				scanWg.Add(1)
				go targetHost.ConnectHost(&scanWg, s.Jitter, s.Debug, host, ScanChan)
				if firstRunS {
					go func() {
						wg.Wait()
						close(ScanChan)
					}()
				}
			} else {
				fmt.Printf("lol scan done\n")
				doneScan = true
			}

		case host, ok := <-ScanChan:
			if ok {
				var fuggoff bool
				// Do dirbusting
				if !s.Dirbust {
					// We're not supposed to dirbust, so let's pump it into the output chan!
					DirbustChan <- host
					break
				}
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
						break
						// panic(err)
					}
					data, err := ioutil.ReadAll(randResp.Body)
					if err != nil {
						// panic(err)
						fuggoff = true
						break
					}
					randResp.Body.Close()
					host.Soft404RandomURL = randURL
					host.Soft404RandomPageContents = strings.Split(string(data), " ")
					s.Soft404edHosts[host.Soft404Hash()] = true
				}
				if !fuggoff {
					if !s.URLProvided {
						for path, _ := range s.Paths.Set {
							// fmt.Printf("HTTP GET to [%v://%v:%v/%v]\n", host.Protocol, host.HostAddr, host.Port, host.Path)
							dirbWg.Add(1)

							go targetHost.HTTPGetter(&dirbWg, host, s.Debug, s.Jitter, s.Soft404Detection, s.StatusCodesIgn, s.Ratio, path, DirbustChan)
						}
					} else {
						dirbWg.Add(1)

						go targetHost.HTTPGetter(&dirbWg, host, s.Debug, s.Jitter, s.Soft404Detection, s.StatusCodesIgn, s.Ratio, host.Path, DirbustChan)
					}
					if firstRunD {
						go func() {
							wg.Wait()
							close(DirbustChan)
						}()
					}
				}
			} else {
				fmt.Printf("lol dirb done\n")

				doneDirbust = true
			}
		case host, ok := <-DirbustChan:
			if ok {
				// Do Screenshotting
				if !s.Screenshot {
					// We're not supposed to screenshot, so let's pump it into the output chan!
					ScreenshotChan <- host
					break
				}
				screenshotWg.Add(1)
				go targetHost.ScreenshotAURL(&screenshotWg, s, cnt, host, ScreenshotChan)

				cnt = cnt + 1
				if firstRunSS {
					go func() {
						wg.Wait()
						close(ScreenshotChan)
					}()
				}
			} else {
				fmt.Printf("lol screenshot done\n")

				doneScreenshot = true
			}
		case <-targetHost:
			break
		default:
			if doneScan && doneDirbust && doneScreenshot {
				fmt.Printf("lol %v hosts donereiszed\n", cnt)
				return
			}
		}

	}

}
