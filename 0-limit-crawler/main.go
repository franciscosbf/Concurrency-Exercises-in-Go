//////////////////////////////////////////////////////////////////////
//
// Your task is to change the code to limit the crawler to at most one
// page per second, while maintaining concurrency (in other words,
// Crawl() must be called concurrently)
//
// @hint: you can achieve this by adding 3 lines
//

package main

import (
	"fmt"
	"sync"
	"time"
)

const rateLimit = time.Second / 1 // One call per second

// Crawl uses `fetcher` from the `mockfetcher.go` file to imitate a
// real crawler. It crawls until the maximum depth has reached.
func Crawl(url string, depth int, wg *sync.WaitGroup, throttle <-chan time.Time) {
	defer wg.Done()

	// Testing purposes
	defer func() { fetchSignal <-true }()

	<-throttle

	if depth <= 0 {
		return
	}

	body, urls, err := fetcher.Fetch(url)
	if err != nil {
		fmt.Println(err)

		// It means that has finiched too soon.
		// Then we need to add this waiting thing
		// so we can throttle even more time
		<-throttle 

		return
	}

	fmt.Printf("found: %s %q\n", url, body)

	wg.Add(len(urls))
	for _, u := range urls {
		<-throttle

		// Do not remove the `go` keyword, as Crawl() must be
		// called concurrently
		go Crawl(u, depth-1, wg, throttle)
	}

	return
}

func main() {
	var wg sync.WaitGroup

	throttle := time.Tick(rateLimit)

	wg.Add(1)
	Crawl("http://golang.org/", 4, &wg, throttle)
	wg.Wait()
}
