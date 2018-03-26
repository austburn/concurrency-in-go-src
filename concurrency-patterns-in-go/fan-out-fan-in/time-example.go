package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	integerGen := func(done <-chan bool, n int) <-chan int {
		r := rand.New(rand.NewSource(99))
		intGen := make(chan int)
		go func() {
			defer close(intGen)
			for i := 0; i < n; i++ {
				intGen <- r.Intn(10)
			}
		}()
		return intGen
	}

	timeGen := func(done <-chan bool, intChan <-chan int) <-chan string {
		msgChan := make(chan string)
		go func() {
			defer close(msgChan)
			for x := range intChan {
				select {
				case <-done:
					return
				default:
					fmt.Println("Starting to sleep for", x, "seconds")
					time.Sleep(time.Duration(x) * time.Second)
					msgChan <- fmt.Sprintf("I just finished sleeping for %d seconds", x)
				}
			}
		}()
		return msgChan
	}

	done := make(chan bool)
	defer close(done)
	start := time.Now()
	messages := timeGen(done, integerGen(done, 10))

	for m := range messages {
		fmt.Println(m)
	}
	fmt.Printf("Took %f seconds", time.Since(start).Seconds())
}
