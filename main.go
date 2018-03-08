package main

import (
	"bufio"
	"flag"
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

// readLines reads a whole file into memory
// and returns a slice of its lines.
func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func main() {
	foundUrls := make(map[string]URLShit)

	argUrls := flag.String("urls", "", "Provide space separated list of urls to grab eh")
	argUrlFile := flag.String("url_file", "	", "Provide space separated list of urls to grab eh")
	var urls []string
	flag.Parse()
	if *argUrlFile != "" {
		argFileUrls, err := readLines(*argUrlFile)
		if err != nil {
			fmt.Print("File: %s does not exist, or you do not have permz (%s)", *argUrlFile, err)
		}
		urls = append(urls, argFileUrls...)
	}
	urls = append(urls, strings.Split(*argUrls, " ")...)
	// Channels
	chUrls := make(chan URLShit)
	chFinished := make(chan struct{})

	wg := sync.WaitGroup{}
	// Kick off the crawl process (concurrently)
	for _, url := range urls {
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
