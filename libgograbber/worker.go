package libgograbber

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"strings"
	"sync"
	"time"
)

func RoutineManager(s *State, ScanChan chan Host, DirbustChan chan Host, ScreenshotChan chan Host, wg *sync.WaitGroup) {
	defer wg.Done()
	threadChan := make(chan struct{}, s.Threads)
	var scanWg = sync.WaitGroup{}
	var dirbWg = sync.WaitGroup{}
	var screenshotWg = sync.WaitGroup{}
	currTime := GetTimeString()

	if s.Debug {
		ticker := time.NewTicker(10 * time.Second)
		startTime := time.Now()
		go func() {
			var currTime time.Duration
			for t := range ticker.C {
				currTime = t.Sub(startTime)
				fmt.Printf(LineSep())
				Debug.Printf("Elapsed %v\n", currTime)
				fmt.Printf(LineSep())

			}
		}()
	}
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
		sWriteChan := make(chan []byte)
		var portScanOutFile string

		if s.ProjectName != "" {
			portScanOutFile = fmt.Sprintf("%v/hosts_%v_%v_%v.txt", s.ScanOutputDirectory, strings.ToLower(strings.Replace(s.ProjectName, " ", "_", -1)), currTime, rand.Int63())
		} else {
			portScanOutFile = fmt.Sprintf("%v/hosts_%v_%v_%v.txt", s.ScanOutputDirectory, currTime, rand.Int63())
		}
		go writerWorker(sWriteChan, portScanOutFile)
		for host := range s.Targets {
			scanWg.Add(1)
			threadChan <- struct{}{}
			go ConnectHost(&scanWg, s.Timeout, s.Jitter, s.Debug, host, ScanChan, threadChan, sWriteChan)
		}
		scanWg.Wait()
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
		// Do dirbusting
		var dirbustOutFile string

		dWriteChan := make(chan []byte)

		if s.ProjectName != "" {
			dirbustOutFile = fmt.Sprintf("%v/urls_%v_%v_%v.txt", s.DirbustOutputDirectory, strings.ToLower(SanitiseFilename(s.ProjectName)), currTime, rand.Int63())
		} else {
			dirbustOutFile = fmt.Sprintf("%v/urls_%v_%v_%v.txt", s.DirbustOutputDirectory, currTime, rand.Int63())
		}
		go writerWorker(dWriteChan, dirbustOutFile)
		// var xwg = sync.WaitGroup{}
		// dirbWg.Add(1)
		for host := range ScanChan {
			dirbWg.Add(1)
			host.Cookies = s.Cookies
			for hostHeader, _ := range s.HostHeaders.Set {
				dirbWg.Add(1)
				host.HostHeader = hostHeader
				if s.URLProvided {
					var h Host
					h = host
					// I think the modification inplace of the host object was creating a problem when accessed later in the dir.go file?
					dirbWg.Add(1)
					go func() {
						defer dirbWg.Done()
						// defer xwg.Done()

						if s.Soft404Detection {
							knary := RandString()
							if s.Canary != "" {
								knary = s.Canary
							}
							randURL := fmt.Sprintf("%v://%v:%v/%v", h.Protocol, h.HostAddr, h.Port, knary)
							if s.Debug {
								Debug.Printf("Soft404 checking [%v]\n", randURL)
							}
							_, randResp, err := h.makeHTTPRequest(randURL)
							if err != nil {
								if s.Debug {
									Error.Printf("Soft404 check failed... [%v] Err:[%v] \n", randURL, err)
								}
							} else {
								defer randResp.Body.Close()
								data, err := ioutil.ReadAll(randResp.Body)
								if err != nil {
									Error.Printf("uhhh... [%v]\n", err)
									return
								}
								h.Soft404RandomURL = randURL
								h.Soft404RandomPageContents = strings.Split(string(data), " ")
							}
						}
						for path, _ := range s.Paths.Set {
							var p string
							p = fmt.Sprintf("%v/%v", strings.TrimSuffix(h.Path, "/"), strings.TrimPrefix(path, "/"))
							dirbWg.Add(1)
							threadChan <- struct{}{}
							go HTTPGetter(&dirbWg, h, s.Debug, s.Jitter, s.Soft404Detection, s.StatusCodesIgn, s.Ratio, p, DirbustChan, threadChan, s.ProjectName, s.HTTPResponseDirectory, dWriteChan, s.FollowRedirects)
						}
					}()
				} else {
					for scheme := range s.Protocols.Set {
						var h Host
						h = host
						h.Protocol = scheme // Weird hack to fix a random race condition...
						// I think the modification inplace of the host object was creating a problem when accessed later in the dir.go file?
						// xwg.Add(1)
						dirbWg.Add(1)
						go func() {
							defer dirbWg.Done()
							// defer xwg.Done()

							if s.Soft404Detection {
								knary := RandString()
								if s.Canary != "" {
									knary = s.Canary
								}
								randURL := fmt.Sprintf("%v://%v:%v/%v", h.Protocol, h.HostAddr, h.Port, knary)
								if s.Debug {
									Debug.Printf("Soft404 checking [%v]\n", randURL)
								}
								_, randResp, err := h.makeHTTPRequest(randURL)
								if err != nil {
									if s.Debug {
										Error.Printf("Soft404 check failed... [%v] Err:[%v] \n", randURL, err)
									}
								} else {
									defer randResp.Body.Close()
									data, err := ioutil.ReadAll(randResp.Body)
									if err != nil {
										Error.Printf("uhhh... [%v]\n", err)
										return
									}
									h.Soft404RandomURL = randURL
									h.Soft404RandomPageContents = strings.Split(string(data), " ")
								}
							}

							for path, _ := range s.Paths.Set {
								dirbWg.Add(1)
								threadChan <- struct{}{}
								go HTTPGetter(&dirbWg, h, s.Debug, s.Jitter, s.Soft404Detection, s.StatusCodesIgn, s.Ratio, path, DirbustChan, threadChan, s.ProjectName, s.HTTPResponseDirectory, dWriteChan, s.FollowRedirects)
							}

						}()
					}
				}
				dirbWg.Done()
			}
			dirbWg.Done()
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
	}()
}
