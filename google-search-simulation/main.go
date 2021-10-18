package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Result string
type Search func(query string) Result

// for each type of search we have two replica servers
var (
	Web1   = fakeSearch("web1")
	Web2   = fakeSearch("web2")
	Image1 = fakeSearch("image1")
	Image2 = fakeSearch("image2")
	Video1 = fakeSearch("video1")
	Video2 = fakeSearch("video2")
)

// randomly sleep to simulate server response time
func fakeSearch(kind string) Search {
	return func(query string) Result {
		time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
		return Result(fmt.Sprintf("%s result for %q\n", kind, query))
	}
}

// only returns the results from the fastest server
func First(query string, replicas ...Search) Result {
	c := make(chan Result)

	for i := range replicas {
		go func(i int) {
			// write result to the channel
			c <- replicas[i](query)
		}(i)
	}

	// Only waits for the first response
	return <-c
}

func Google(query string) []Result {
	c := make(chan Result)

	// performs each type of search in a separate goroutine
	go func() { c <- First(query, Web1, Web2) }()
	go func() { c <- First(query, Image1, Image2) }()
	go func() { c <- First(query, Video1, Video2) }()

	var results []Result

	// defines a global timeout of 80ms
	// all results that are taking longer than 80ms are ignored
	timeout := time.After(80 * time.Millisecond)

	for i := 0; i < 3; i++ {
		select {
		case r := <-c:
			results = append(results, r)
		// ignores the slowest server
		case <-timeout:
			fmt.Println("timed out")
			return results
		}
	}

	return results
}

func main() {
	rand.Seed(time.Now().UnixNano())
	start := time.Now()
	results := Google("golang")
	elapsed := time.Since(start)
	fmt.Println(results)
	fmt.Println(elapsed)
}
