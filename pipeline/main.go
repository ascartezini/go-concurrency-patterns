package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strings"
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

func writableStream(data <-chan string, file *os.File, quit chan bool) {

	go func() {
		for v := range data {
			file.WriteString(v)
		}

		err := file.Close()

		if err != nil {
			fmt.Println("file closed")
		}

		quit <- true
	}()
}

func main() {
	path, err := os.Getwd()

	if err != nil {
		fmt.Println("error when getting working directory")
	}

	file, err := os.Create(path + "/data.csv")
	rand.Seed(time.Now().UnixNano())
	start := time.Now()

	quit := make(chan bool)
	c := readableStream()
	t := transform(c)
	writableStream(t, file, quit)

	<-quit

	elapsed := time.Since(start)
	fmt.Println(elapsed)
}
