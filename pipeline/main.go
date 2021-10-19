package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"
)

type Person struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func readableStream() <-chan string {
	c := make(chan string)

	go func() {
		for i := 0; i < 1e6; i++ {
			person := &Person{Id: +i, Name: fmt.Sprintf("Andre-%d", i)}
			data, err := json.Marshal(person)

			if err != nil {
				fmt.Println(err)
				break
			}

			c <- string(data)
		}
		close(c)
	}()

	return c
}

func transform(data <-chan string) <-chan string {

	c := make(chan string)

	go func() {
		c <- "id,name\n"

		for v := range data {
			person := Person{}
			json.Unmarshal([]byte(v), &person)
			transformedData := fmt.Sprintf("%d,%s\n", person.Id, strings.ToUpper(person.Name))

			c <- transformedData
		}
		close(c)
	}()

	return c
}

func writableStream(data <-chan string, file *os.File, done chan bool) <-chan string {

	c := make(chan string)

	go func() {
		for v := range data {
			file.WriteString(v)
		}

		err := file.Close()

		if err != nil {
			fmt.Println("file closed")
		}

		close(c)
		done <- true
	}()

	return c
}

func merge(cs ...<-chan string) <-chan string {
	var wg sync.WaitGroup
	out := make(chan string)

	// Start an output goroutine for each input channel in cs.  output
	// copies values from c to out until c is closed, then calls wg.Done.
	output := func(c <-chan string) {
		for n := range c {
			out <- n
		}
		wg.Done()
	}
	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
	}

	// Start a goroutine to close out once all the output goroutines are
	// done.  This must start after the wg.Add call.
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func main() {
	path, err := os.Getwd()

	if err != nil {
		fmt.Println("error when getting working directory")
	}

	file, err := os.Create(path + "/data.csv")
	rand.Seed(time.Now().UnixNano())
	start := time.Now()

	done := make(chan bool)
	c := readableStream()
	t1 := transform(c)
	t2 := transform(c)
	mergedTransform := merge(t1, t2)
	writableStream(mergedTransform, file, done)

	<-done

	elapsed := time.Since(start)
	fmt.Println(elapsed)
}
