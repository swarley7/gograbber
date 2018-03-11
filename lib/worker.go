package lib

import (
	"bufio"
	"bytes"
	"fmt"
	"net/http"
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

func headder(t string) (status int, length int64, err error) {
	req, err := http.NewRequest("HEAD", t+"/.git/index", nil)
	if err != nil {
		fmt.Println(strings.Repeat("~", 40))
		fmt.Println("Bad Req Construction")
		fmt.Println(err)
		fmt.Println(strings.Repeat("~", 40))
		return 0, 0, err
	}
	res, err := cl.Do(req)
	if res != nil && res.Body != nil {
		defer res.Body.Close()
	}
	if err != nil {
		return 0, 0, err
	}
	return res.StatusCode, res.ContentLength, nil
}

func getter(t string) (code int, body []byte, err error) {
	req, err := http.NewRequest("GET", t+"/.git/config", nil)
	if err != nil {
		fmt.Println(strings.Repeat("~", 40))
		fmt.Println("Bad Req Construction")
		fmt.Println(err)
		fmt.Println(strings.Repeat("~", 40))
	}

	//send off that request
	resp, err := cl.Do(req)
	if resp != nil {
		defer resp.Body.Close() //remember to close the body bro
	}
	if err != nil {
		return 404, nil, err
	}
	buf := &bytes.Buffer{}
	buf.ReadFrom(resp.Body)
	body = buf.Bytes()

	return resp.StatusCode, body, nil
}

var filled = false

func routineManager(finishedInput chan struct{}, threads int, indexChan chan string, configChan chan string, writeChan chan []byte, wg *sync.WaitGroup) {
	q := make(chan struct{}, threads)
	lolgroup := sync.WaitGroup{}
	for {
		//fmt.Println("starting..", len(q), filled, len(indexChan), len(configChan), len(writeChan))
		select {
		case q <- struct{}{}:
			lolgroup.Add(1)
			go taskWorker(&lolgroup, indexChan, configChan, writeChan, q)
		case _, ok := <-finishedInput:
			if !ok {
				filled = true
				finishedInput = nil
			}
		}
		if filled && len(indexChan) == 0 && len(configChan) == 0 && len(writeChan) == 0 {
			fmt.Println("WAITING")
			lolgroup.Wait()
			fmt.Println("BREAKING")
			wg.Done()
			return
		}
		//fmt.Println("waiting..", len(q), filled, len(indexChan), len(configChan), len(writeChan))
		//time.Sleep(time.Second * 2)
	}
}

func taskWorker(lolgroup *sync.WaitGroup, indexChan chan string, configChan chan string, writeChan chan []byte, finishedIndicator chan struct{}) {
	defer func() {
		_ = <-finishedIndicator
		lolgroup.Done()
	}()
	select {
	case t := <-indexChan:
		//send a HEAD request for the index file. Add to config queue if successful
		code, le, err := headder(t)
		if err != nil {
			return
		}

		if le < 10 || code != 200 {
			return
		}
		//content length is over 10, status code is 200, likely a good result.

		configChan <- t

	case t := <-configChan:
		//try to get the config file, write to disk if successful
		code, body, err := getter(t)

		if err != nil {
			return
		}

		if code != 200 {
			return
		}

		if strings.Contains(string(body), "[core]") {
			if strings.Contains(strings.ToLower(string(body)), "<!doctype") {
				return //if we got some weird html through for some reason
			}
			//we got a vuln
			writeChan <- append([]byte("\n~$$ Git found: "+t+" $$~\n"), body...)
		}
	}
}
