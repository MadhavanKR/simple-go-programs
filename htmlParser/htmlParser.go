package htmlParser

import (
	"fmt"
	"io"
	"strings"

	"golang.org/x/net/html"
)

type Link struct {
	Href string
	Text string
}

func ParseHtml(htmlContent io.Reader) ([]Link, error) {
	htmlRootNode, parseErr := html.Parse(htmlContent)
	if parseErr != nil {
		fmt.Println("error while parsing html: ", parseErr)
		return nil, parseErr
	}
	links := traverse(htmlRootNode)
	return links, nil
}

func buildLink(node *html.Node) Link {
	var link Link
	for _, attr := range node.Attr {
		if attr.Key == "href" {
			link.Href = attr.Val
			break
		}
	}

	linkText := extractText(node)
	link.Text = linkText
	return link
}

func extractText(node *html.Node) string {
	if node.Type == html.TextNode {
		return node.Data
	} else if node.Type != html.ElementNode {
		return ""
	}
	linkText := ""
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		linkText += strings.TrimSpace(extractText(child))
	}
	return linkText
}

func traverse(node *html.Node) []Link {
	if node.Type == html.ElementNode && node.Data == "a" {
		return []Link{buildLink(node)}
	}
	var links []Link
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		links = append(links, traverse(child)...)
	}
	return links
}
