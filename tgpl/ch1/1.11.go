package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

const (
	URL_LIMIT = 100
)

func main() {
	start := time.Now()
	sitesFile := os.Args[1]
	f, err := os.Open(sitesFile)
	if err != nil {
		fmt.Printf("file error: %v", err)
		return
	}
	r := csv.NewReader(bufio.NewReader(f))
	records, err := r.ReadAll()

	ch := make(chan string)
	for index, record := range records {
		url := fmt.Sprintf("http://%s", record[1])
		go fetch(url, ch)
		if index >= URL_LIMIT {
			break
		}
	}
	for i:= 0; i<=URL_LIMIT; i++ {
		fmt.Println(<-ch)
	}
	fmt.Printf(".2fs elapsed\n", time.Since(start).Seconds())
	f.Close()
}

func fetch(url string, ch chan<- string) {
	start := time.Now()
	resp, err := http.Get(url)
	if err != nil {
		ch <- fmt.Sprint(err)
		return
	}
	nbytes, err := io.Copy(ioutil.Discard, resp.Body)
	resp.Body.Close()
	if err != nil {
		ch <- fmt.Sprintf("while reading %s: %v", url, err)
		return
	}
	secs := time.Since(start).Seconds()
	ch <- fmt.Sprintf("%.2fs %7d %s", secs, nbytes, url)
}
