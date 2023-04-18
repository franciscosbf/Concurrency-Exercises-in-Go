//////////////////////////////////////////////////////////////////////
//
// Given is a producer-consumer scenario, where a producer reads in
// tweets from a mockstream and a consumer is processing the
// data. Your task is to change the code so that the producer as well
// as the consumer can run concurrently
//

package main

import (
	"fmt"
	"sync"
	"time"
)

func producer(tweets chan *Tweet, stream Stream, wg *sync.WaitGroup) {
	defer func(){ wg.Done() }()

	for {
		tweet, err := stream.Next()
		if err == ErrEOF {
			break	
		}

		tweets <-tweet
	}

	// Closes the channel so the consumer nows
	// when there's no more items to fetch
	close(tweets)
}

func consumer(tweets chan *Tweet, wg *sync.WaitGroup) {
	defer func(){ wg.Done() }()

	for t := range tweets {
		if t.IsTalkingAboutGo() {
			fmt.Println(t.Username, "\ttweets about golang")
		} else {
			fmt.Println(t.Username, "\tdoes not tweet about golang")
		}
	}
}

func main() {
	var wg sync.WaitGroup

	// This chan is used to send tweets 
	// from the producer to the consumer
	buffer := make(chan *Tweet, 20)

	start := time.Now()
	stream := GetMockStream()

	// Producer
	wg.Add(1)
	go producer(buffer, stream, &wg)

	// Consumer
	wg.Add(1)
	go consumer(buffer, &wg)

	wg.Wait()

	fmt.Printf("Process took %s\n", time.Since(start))
}
