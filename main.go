package main

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"
)

type counter struct {
	sync.Mutex
	count int
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func parse(reader io.ReadCloser) int {
	body, err := ioutil.ReadAll(reader)
	checkError(err)
	return strings.Count(string(body), "Go")
}

func main() {
	k := 5
	urls := []string{
		"https://google.com",
		"https://golang.org",
		"https://golang.org",
		"https://golang.org",
		"https://golang.org",
		"https://golang.org",
		"https://google.com",
	}

	var wg sync.WaitGroup
	goroutines := make(chan struct{}, k)
	count := &counter{}

	for _, url := range urls {

		goroutines <- struct{}{}
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			request, err := http.NewRequest("GET", url, nil)
			checkError(err)
			response, err := http.DefaultClient.Do(request)
			checkError(err)

			localCount := parse(response.Body)

			count.Lock()
			count.count += localCount
			count.Unlock()

			err = response.Body.Close()
			checkError(err)

			log.Println(url, " count:", localCount)
			<-goroutines
		}(url)

	}
	wg.Wait()
	log.Println(count.count)
}
