package main

import (
	"fmt"

	"flag"

	"./urlShortner"
)

func main() {
	port := flag.String("port", "2527", "Port at which the http server must be started")
	yamlFilename := flag.String("yaml", "urls.yml", "yaml file containing the map of short url to actual urls")
	flag.Parse()
	fmt.Println("starting the server at ", *port)

	urlShortner.StartHttpServer(":"+*port, *yamlFilename)
}
