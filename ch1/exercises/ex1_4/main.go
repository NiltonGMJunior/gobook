package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	counts := make(map[string]map[string]int)
	files := os.Args[1:]
	if len(files) == 0 {
		countLines(os.Stdin, counts)
	} else {
		for _, arg := range files {
			f, err := os.Open(arg)
			if err != nil {
				fmt.Fprintf(os.Stderr, "dup2: %v\n", err)
				continue
			}
			counts[f.Name()] = make(map[string]int)
			countLines(f, counts)
			f.Close()
		}
	}
	for filename, lineCounts := range counts {
		for _, n := range lineCounts {
			if n > 1 {
				fmt.Printf("%s\n", filename)
				continue
			}
		}
	}
}

func countLines(f *os.File, counts map[string]map[string]int) {
	input := bufio.NewScanner(f)
	for input.Scan() {
		counts[f.Name()][input.Text()]++
	}
	// NOTE: ignoring potential errors from input.Err()
}
