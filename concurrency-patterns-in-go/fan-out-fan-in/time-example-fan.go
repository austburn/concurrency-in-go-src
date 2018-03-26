package main

import (
	"fmt"
	"math/rand"
	"sync"
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

	merge := func(done <-chan bool, cs ...<-chan string) <-chan string {
		var wg sync.WaitGroup
		out := make(chan string)

		// Start an output goroutine for each input channel in cs.  output
		// copies values from c to out until c is closed, then calls wg.Done.
		output := func(c <-chan string) {
			for n := range c {
				select {
				case <-done:
					return
				case out <- n:
				}
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

	done := make(chan bool)
	defer close(done)
	start := time.Now()

	workers := make([]<-chan string, 4)
	integers := integerGen(done, 10)

	for i := 0; i < 4; i++ {
		workers[i] = timeGen(done, integers)
	}

	for m := range merge(done, workers...) {
		fmt.Println(m)
	}
	fmt.Printf("Took %f seconds", time.Since(start).Seconds())
}
