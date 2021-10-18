package main

import (
	"fmt"
	"time"
)

// generator: function that returns a channel
func talk(msg string, sleep int) <-chan string { // returns receive-only channel of strings
	c := make(chan string)

	go func() { // launchs goroutine inside the function
		for i := 0; ; i++ {
			c <- fmt.Sprintf("%s %d", msg, i)
			time.Sleep(time.Duration(sleep) * time.Millisecond)
		}
	}()

	return c // returns channel to the caller
}

// takes two channels as parameters and create a third one
func fanIn(input1, input2 <-chan string) <-chan string {
	c := make(chan string)

	// launchs two independent goroutines to write to the single channel
	go func() {
		for {
			c <- <-input1
		}
	}()
	go func() {
		for {
			c <- <-input2
		}
	}()

	return c
}

func main() {
	fastTalker := talk("I talk really fast!", 500)
	slowTalker := talk("I talk really slow!", 2000)

	// using the fan-in function in order to the slowest talker does not block the fastest one
	c := fanIn(fastTalker, slowTalker)

	for i := 0; i < 10; i++ {
		fmt.Printf("You say: %q\n", <-c)
	}

	fmt.Println("I'm leaving, bye!")
}
