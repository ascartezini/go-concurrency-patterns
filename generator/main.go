package main

import (
	"fmt"
	"math/rand"
	"time"
)

// generator: function that returns a channel
func sendMessage(msg string) <-chan string { // returns receive-only channel of strings
	c := make(chan string)

	go func() { // launchs goroutine inside the function
		for i := 0; ; i++ {
			c <- fmt.Sprintf("%s %d", msg, i)
			time.Sleep(time.Duration(rand.Intn(1e3)) * time.Millisecond)
		}
	}()

	return c // returns channel to the caller
}

func main() {
	c := sendMessage("hello!")

	for i := 0; i < 5; i++ {
		fmt.Printf("You say: %q\n", <-c)
	}

	fmt.Println("I'm leaving, bye!")
}
