// Fetchall fetches URLs in parallel and reports their times and sizes.
package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

func main() {
	start := time.Now()
	ch := make(chan string)
	outputFile, err := os.Create(os.Args[1])
	defer outputFile.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "while creating output file: %v\n", err)
		os.Exit(1)
	}
	writer := bufio.NewWriter(outputFile)
	for _, url := range os.Args[2:] {
		go fetch(url, ch) // start a goroutine
	}
	for range os.Args[2:] {
		writer.WriteString(fmt.Sprintf("%s\n", <-ch))
	}
	writer.Flush()
	fmt.Printf("%.2fs elapsed\n", time.Since(start).Seconds())
}

func fetch(url string, ch chan<- string) {
	start := time.Now()
	resp, err := http.Get(url)
	if err != nil {
		ch <- fmt.Sprint(err) // send to channel ch
		return
	}

	nbytes, err := io.Copy(ioutil.Discard, resp.Body)
	resp.Body.Close() // don´t leak resources
	if err != nil {
		ch <- fmt.Sprintf("while reading %s: %v", url, err)
		return
	}
	secs := time.Since(start).Seconds()
	ch <- fmt.Sprintf("%.2fs %7d %s", secs, nbytes, url)
}
