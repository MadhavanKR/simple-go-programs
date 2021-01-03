package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"sync"
)

func main() {
	var numRequests, blockSize int
	flag.IntVar(&numRequests, "numCalls", 1, "defines number of calls to be made")
	flag.IntVar(&blockSize, "blockSize", 99999, "defines number of calls to be made in one go")
	flag.Parse()
	waitGroup := sync.WaitGroup{}
	loadTestConfig := parseConfigJSON("loadtestConfig.json")
	for _, testConfig := range loadTestConfig {
		waitGroup.Add(1)
		go startTest(testConfig, &waitGroup)
	}
	waitGroup.Wait()
}

func startTest(requestConfig config, startTestWg *sync.WaitGroup) {
	defer startTestWg.Done()
	fmt.Println("starting test on: ", requestConfig.URL)
	waitGroup := sync.WaitGroup{}
	count := 0
	for i := 0; i < requestConfig.NumCalls; i++ {
		waitGroup.Add(1)
		count = count + 1
		go fetchToken(&waitGroup, requestConfig)
		if count == requestConfig.BlockSize {
			waitGroup.Wait()
			count = 0
		}
	}
	waitGroup.Wait()
}

/*
{
        "url": "http://localhost:8080/",
        "httpMethod": "GET",
        "data": {},
        "headers": {
            "Content-Type": "application/json",
            "Authorization": "Basic dXNlcm5hbWU6cGFzc3dvcmQ="
        }
    }
*/

type config struct {
	URL        string                 `json: "url"`
	HTTPMethod string                 `json: "httpMethod`
	Data       map[string]interface{} `json: "data"`
	Headers    map[string]string      `json: "headers"`
	NumCalls   int                    `json: "numCalls"`
	BlockSize  int                    `json: "blockSize"`
}

func parseConfigJSON(configFileName string) []config {
	configBytes, configReadErr := ioutil.ReadFile(configFileName)
	if configReadErr != nil {
		fmt.Println("error while reading config file: ", configReadErr)
		os.Exit(1)
	}
	var loadTestConfig []config
	unmarshalErr := json.Unmarshal(configBytes, &loadTestConfig)
	//fmt.Println(loadTestConfig)
	if unmarshalErr != nil {
		fmt.Println("error while unmarshalling config file: ", unmarshalErr)
		os.Exit(1)
	}
	return loadTestConfig
}

func fetchToken(waitGroup *sync.WaitGroup, requestConfig config) {
	defer waitGroup.Done()
	var dataBytes []byte
	if requestConfig.Headers["Content-Type"] == "application/x-www-form-urlencoded" {
		data := url.Values{}
		for key, value := range requestConfig.Data {
			if valueStr, ok := value.(string); ok {
				data.Set(key, valueStr)
			} else {
				fmt.Printf("all key values must be string for content type: %s\n", requestConfig.Headers["Content-Type"])
				return
			}
		}
		dataBytes = []byte(data.Encode())
	} else if requestConfig.Headers["Content-Type"] == "application/json" {
		data, marshalErr := json.Marshal(requestConfig.Data)
		fmt.Println("data is: ", string(data))
		if marshalErr != nil {
			fmt.Println("error while marshalling request data: ", marshalErr)
			return
		}
		dataBytes = data
	} else {
		fmt.Printf("Content-Type: %s not supported\n", requestConfig.Headers["Content-Type"])
		return
	}
	tokenRequest, rqstCreateErr := http.NewRequest(requestConfig.HTTPMethod, requestConfig.URL, bytes.NewReader(dataBytes))
	if rqstCreateErr != nil {
		fmt.Println("error while creating post request for token: ", rqstCreateErr)
		os.Exit(1)
	}

	//adding headers
	for headerKey, headerValue := range requestConfig.Headers {
		tokenRequest.Header.Add(headerKey, headerValue)
	}

	var httpClient = &http.Client{}
	tokenResponse, tokenFetchErr := httpClient.Do(tokenRequest)
	if tokenFetchErr != nil {
		fmt.Println("error while fetching token: ", tokenFetchErr)
		os.Exit(1)
	}
	defer httpClient.CloseIdleConnections()
	defer tokenResponse.Body.Close()
	if tokenResponse.StatusCode != 200 {
		errorDescBytes, _ := ioutil.ReadAll(tokenResponse.Body)
		fmt.Printf("request failed: Http Status %s, reason: %s", tokenResponse.Status, string(errorDescBytes))
	} else {
		fmt.Println("request completed for: ", requestConfig.URL)
	}
}
