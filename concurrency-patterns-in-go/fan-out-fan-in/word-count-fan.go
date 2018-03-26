package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"runtime"
	"sort"
	"sync"
	"time"
)

type Word struct {
	w string
	c int
}

type WordCount []Word

func (w WordCount) Len() int           { return len(w) }
func (w WordCount) Swap(i, j int)      { w[i], w[j] = w[j], w[i] }
func (w WordCount) Less(i, j int) bool { return w[i].c > w[j].c }

func main() {
	start := time.Now()

	urls := []string{
		"http://www.gutenberg.org/files/1921/1921-0.txt",
		"http://www.gutenberg.org/files/53362/53362-0.txt",
		"http://www.gutenberg.org/files/53357/53357-0.txt",
		"http://www.gutenberg.org/files/53358/53358-0.txt",
		"http://www.gutenberg.org/files/53359/53359-0.txt",
		"http://www.gutenberg.org/files/53360/53360-0.txt",
		"http://www.gutenberg.org/files/53361/53361-0.txt",
		"http://www.gutenberg.org/files/53362/53352-0.txt",
	}

	done := make(chan bool)
	defer close(done)

	urlGen := GenUrls(done, urls...)

	numFinders := runtime.NumCPU()
	fmt.Printf("Spinning up %d WordCounts.\n", numFinders)
	finders := make([]<-chan string, numFinders)
	for i := 0; i < numFinders; i++ {
		finders[i] = GetWords(done, urlGen)
	}

	wordCountMap := FanIn(done, finders...)

	var wordCountStruct WordCount
	wordCountMap.Range(func(k, v interface{}) bool {
		wordCountStruct = append(wordCountStruct, Word{w: k.(string), c: v.(int)})
		return true
	})
	sort.Sort(wordCountStruct)
	log.Println(wordCountStruct[:10])

	log.Printf("Took: %fs", time.Since(start).Seconds())
}

func GenUrls(done <-chan bool, urls ...string) <-chan string {
	urlStream := make(chan string)
	go func() {
		defer close(urlStream)
		for {
			for _, u := range urls {
				select {
				case <-done:
					return
				case urlStream <- u:
				}
			}
		}
	}()
	return urlStream
}

func GetWords(done <-chan bool, urls <-chan string) <-chan string {
	wordStream := make(chan string)
	go func() {
		defer close(wordStream)
		for {
			select {
			case <-done:
				return
			case u := <-urls:
				r, err := http.Get(u)
				if err != nil {
					return
				}
				b, err := ioutil.ReadAll(r.Body)
				if err != nil {
					return
				}
				words := bytes.Split(bytes.ToLower(b), []byte(" "))
				for _, word := range words {
					wordStream <- string(word)
				}
			}
		}
	}()
	return wordStream
}

func FanIn(done <-chan bool, channels ...<-chan string) sync.Map { // <1>
	var wg sync.WaitGroup
	var syncMap sync.Map

	multiplex := func(c <-chan string) { // <3>
		defer wg.Done()
		for {
			select {
			case <-done:
				return
			case s := <-c:
				i, ok := syncMap.Load(s)
				if ok {
					syncMap.Store(s, i.(int)+1)
				} else {
					syncMap.Store(s, 0)
				}
			}
		}
	}

	// Select from all the channels
	wg.Add(len(channels)) // <4>
	for _, c := range channels {
		go multiplex(c)
	}

	// Wait for all the reads to complete
	go func() { // <5>
		wg.Wait()
	}()

	return syncMap
}

// 2018/03/22 16:53:24 [{ 67794} {the 64008} {of 35455} {and 25325} {to 25025} {a 18207} {in 16339} {was 9302} {that 9276} {he 6842}]
