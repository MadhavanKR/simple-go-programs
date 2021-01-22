package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	ch := make(chan int)
	for i := 0; i < 5; i++ {
		go sampleGoRoutine(i, ch)
	}
	for i := 0; i < 5; i++ {
		fmt.Println("recieved data from: ", <-ch)
	}
	for i := 5; i < 10; i++ {
		go sampleGoRoutine(i, ch)
	}

	for i := 0; i < 5; i++ {
		fmt.Println("recieved data from: ", <-ch)
	}
}

func sampleGoRoutine(id int, ch chan int) {
	waitSeconds := rand.Int() % 5
	fmt.Printf("in method call %d, waiting for %d seconds \n", id, waitSeconds)
	//time.Sleep(time.Duration(waitSeconds) * time.Second)
	time.Sleep(time.Duration(5) * time.Second)
	ch <- id
}
