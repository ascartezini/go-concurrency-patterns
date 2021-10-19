package main

import "fmt"

func main() {
	fmt.Println(getNumber())
}

func getNumber() int {
	var i int
	done := make(chan bool)
	go func() {
		i = 5
		done <- true
	}()

	<-done
	return i
}
