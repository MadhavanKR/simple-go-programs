package main

/*
	reference: https://github.com/gophercises/sitemap
*/

import (
	"flag"
	"fmt"

	"./siteMapBuilder"
)

func main() {
	domainUrl := flag.String("rootUrl", "https://gophercises.com", "the root url for which you wish to obtain a siteMap")
	flag.Parse()
	fmt.Println("Preparing sitemap for: ", *domainUrl)
	buildSiteErr := siteMapBuilder.BuildSiteMap(*domainUrl)
	if buildSiteErr != nil {
		panic(buildSiteErr)
	}
}
