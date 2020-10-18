package siteMapBuilder

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"../htmlParser"
)

func fetchAndParseSiteContent(url string) ([]htmlParser.Link, error) {
	response, httpErr := http.Get(url)
	if httpErr != nil {
		fmt.Println("error while fetching site content: ", httpErr)
		return nil, httpErr
	}
	buf := new(strings.Builder)
	_, readErr := io.Copy(buf, response.Body)
	if readErr != nil {
		fmt.Println("error while reading response body: ", readErr)
		return nil, readErr
	}
	//fmt.Println("siteContent: ", buf.String())
	links, parseErr := htmlParser.ParseHtml(strings.NewReader(buf.String()))
	if parseErr != nil {
		return nil, parseErr
	}
	return links, nil
}

func BuildSiteMap(url string) error {
	rootSiteLinks, fetchAndParseErr := fetchAndParseSiteContent(url)
	//fmt.Println("rootSiteLinks: ", rootSiteLinks)
	if fetchAndParseErr != nil {
		return fetchAndParseErr
	}
	siteVisitedMap := make(map[string]bool)
	for _, rootSiteLink := range rootSiteLinks {
		if isPathLink(rootSiteLink.Href) {
			siteVisitedMap[getFullDomainLink(url, rootSiteLink.Href)] = false
		} else {
			if isDomainLink(rootSiteLink.Href, url) {
				siteVisitedMap[rootSiteLink.Href] = false
			}
		}
	}
	fmt.Println("len(siteVisitedMap) before: ", len(siteVisitedMap))
	//fmt.Println("siteVistedmap: ", siteVisitedMap)
	exploreAllLinks(siteVisitedMap, url)
	//fmt.Println("all links: ", siteVisitedMap)
	fmt.Println("len(siteVisitedMap) after: ", len(siteVisitedMap))
	for k := range siteVisitedMap {
		fmt.Println(k)
	}
	return nil
}

func exploreAllLinks(siteVisitedMap map[string]bool, domainUrl string) {
	allLinksVisited := true
	for _, v := range siteVisitedMap {
		if v == false {
			allLinksVisited = false
			break
		}
	}
	if allLinksVisited == true {
		return
	}
	for k, v := range siteVisitedMap {
		if v == true {
			continue
		}
		siteLinks, fetchAndParseErr := fetchAndParseSiteContent(k)
		siteVisitedMap[k] = true
		for _, siteLink := range siteLinks {
			href := strings.TrimSpace(siteLink.Href)
			if href == "" {
				//fmt.Println("empty string.. continuing")
				continue
			}
			var domainSiteLink string
			if isPathLink(href) {
				domainSiteLink = getFullDomainLink(domainUrl, href)
			} else {
				if isDomainLink(href, domainUrl) {
					domainSiteLink = strings.Trim(href, "/")
				}
			}
			if _, ok := siteVisitedMap[domainSiteLink]; !ok {
				if domainSiteLink != "" {
					siteVisitedMap[domainSiteLink] = false
				}
			}
		}
		if fetchAndParseErr != nil {
			fmt.Printf("error while parsing links from %s: %v ", k, fetchAndParseErr)
		} else {
			exploreAllLinks(siteVisitedMap, domainUrl)
		}
	}
}

func isDomainLink(linkUrl string, url string) bool {
	return strings.Contains(linkUrl, url)
}

func getFullDomainLink(url string, path string) string {
	domainLink := strings.Trim(url, "/") + "/" + strings.Trim(path, "/")
	return strings.Trim(domainLink, "/")
}

func isPathLink(url string) bool {
	if (strings.HasPrefix(url, "http") == false) && (strings.HasPrefix(url, "/")) {
		//fmt.Printf("%s is pathLink \n", url)
		return true
	} else {
		return false
	}
}
