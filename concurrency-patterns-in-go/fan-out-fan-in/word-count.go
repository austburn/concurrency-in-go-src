package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
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
	wordCount := make(map[string]int)

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

	for _, u := range urls {
		r, err := http.Get(u)
		if err != nil {
			log.Fatal(err)
		}
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
		}
		words := bytes.Split(bytes.ToLower(b), []byte(" "))
		for _, word := range words {
			wordStr := string(word)
			i, ok := wordCount[wordStr]
			if ok {
				wordCount[wordStr] = i + 1
			} else {
				wordCount[wordStr] = 0
			}
		}
	}

	var wordCountStruct WordCount
	for k, v := range wordCount {
		wordCountStruct = append(wordCountStruct, Word{w: k, c: v})
	}
	sort.Sort(wordCountStruct)
	log.Println(wordCountStruct[:10])

	log.Printf("Took: %fs", time.Since(start).Seconds())
}

// 2018/03/22 16:53:24 [{ 67794} {the 64008} {of 35455} {and 25325} {to 25025} {a 18207} {in 16339} {was 9302} {that 9276} {he 6842}]
