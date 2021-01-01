package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
)

func main() {
	var numRequests, blockSize int
	flag.IntVar(&numRequests, "numCalls", 1, "defines number of calls to be made")
	flag.IntVar(&blockSize, "blockSize", 99999, "defines number of calls to be made in one go")
	flag.Parse()
	waitGroup := sync.WaitGroup{}
	count := 0
	for i := 0; i < numRequests; i++ {
		waitGroup.Add(1)
		count = count + 1
		go fetchToken(&waitGroup)
		if count == blockSize {
			waitGroup.Wait()
			count = 0
		}
	}
	waitGroup.Wait()
}

func fetchToken(waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()

	data := url.Values{}
	data.Set("grant_type", "password")
	data.Set("username", "<username>")
	data.Set("password", "<password>")
	tokenRequest, rqstCreateErr := http.NewRequest("POST", "<url>", strings.NewReader(data.Encode()))
	if rqstCreateErr != nil {
		fmt.Println("error while creating post request for token: ", rqstCreateErr)
		os.Exit(1)
	}

	tokenRequest.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	tokenRequest.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	tokenRequest.Header.Add("Authorization", "Basic <base64encode>")

	var httpClient = &http.Client{}
	tokenResponse, tokenFetchErr := httpClient.Do(tokenRequest)
	if tokenFetchErr != nil {
		fmt.Println("error while fetching token: ", tokenFetchErr)
		os.Exit(1)
	}
	defer httpClient.CloseIdleConnections()
	defer tokenResponse.Body.Close()
	if tokenResponse.StatusCode != 200 {
		fmt.Println("failed to fetch token: Http Status ", tokenResponse.Status)
	} else {
		fmt.Println("successfully fetched token: ")
	}
}
