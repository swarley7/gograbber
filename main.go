package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"
)

type URLShit struct {
	URL      string
	BodyData string
	Host     string
	Protocol string
}

// Discover web pages
func crawl(url string, ch chan URLShit, chFinished chan struct{}, wg *sync.WaitGroup) {
	resp, err := http.Get(url)
	defer func() {
		// Notify that we're done after this function
		wg.Done()
		chFinished <- struct{}{}
	}()
	defer resp.Body.Close()
	if err != nil {
		// Something? IDK! fuck
		return
	}
	if resp.StatusCode == http.StatusOK {
		bodyBytes, err2 := ioutil.ReadAll(resp.Body)
		if err2 != nil {
			// Something? IDK! fuck
			return
		}
		bodyString := string(bodyBytes)

		ch <- URLShit{URL: url, BodyData: bodyString, Host: resp.Request.Host, Protocol: resp.Request.Proto}
	}
}

func main() {
	foundUrls := make(map[string]URLShit)
	seedUrls := os.Args[1:]

	// Channels
	chUrls := make(chan URLShit)
	chFinished := make(chan struct{})

	wg := sync.WaitGroup{}
	// Kick off the crawl process (concurrently)
	for _, url := range seedUrls {
		wg.Add(1)
		go crawl(url, chUrls, chFinished, &wg)
	}

	// Subscribe to both channels
	go func() {
		for {
			data := <-chUrls
			foundUrls[data.URL] = data
		}
	}()

	wg.Wait()
	// We're done! Print the results...

	fmt.Println("\nFound", len(foundUrls), "unique urls:")

	for url := range foundUrls {
		fmt.Println(" - " + url)
		filename := fmt.Sprintf("./%s.txt", strings.Replace(foundUrls[url].Host, ".", "_", -1))
		f, err := os.Create(filename)
		fmt.Printf("test? %s", filename)
		if err != nil {
			fmt.Printf("Uh %s %s", err, filename)
			continue
		}
		f.WriteString(foundUrls[url].BodyData)
	}

	close(chUrls)
}
