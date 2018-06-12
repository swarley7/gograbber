package libgograbber

import (
	"fmt"
	"sync"
	"time"
)

func RoutineManager(s *State, ScanChan chan Host, DirbustChan chan Host, ScreenshotChan chan Host, wg *sync.WaitGroup) {
	defer wg.Done()
	threadChan := make(chan struct{}, s.Threads)
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

	// Start our operations
	wg.Add(1)
	go Scan(s, s.Targets, ScanChan, currTime, threadChan, wg)
	wg.Add(1)
	go Dirbust(s, ScanChan, DirbustChan, currTime, threadChan, wg)
	wg.Add(1)
	go Screenshot(s, DirbustChan, ScreenshotChan, currTime, threadChan, wg)
}
