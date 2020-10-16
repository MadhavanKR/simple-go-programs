package main

/*
	reference: https://github.com/gophercises/link/tree/master
*/
import (
	"fmt"
	"strings"

	"./htmlParser"
)

func main() {
	var sampleHtml = `
	<html>
	<body>
	  <h1>Hello!</h1>
	  <a href="/other-page">
		A link to another page
		<span> some span  </span>
	  </a>
	  <a href="/page-two">A link to a second page</a>
	</body>
	</html>
	`
	_ = sampleHtml
	var sampleHtml2 = `<html>
	<head>
	  <link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/font-awesome/4.7.0/css/font-awesome.min.css">
	</head>
	<body>
	  <h1>Social stuffs</h1>
	  <div>
		<a href="https://www.twitter.com/joncalhoun">
		  Check me out on twitter
		  <i class="fa fa-twitter" aria-hidden="true"></i>
		</a>
		<a href="https://github.com/gophercises">
		  Gophercises is on <strong>Github</strong>!
		</a>
	  </div>
	</body>
	</html>`

	links, parseErr := htmlParser.ParseHtml(strings.NewReader(sampleHtml2))
	if parseErr != nil {
		panic(parseErr)
	}
	for _, link := range links {
		fmt.Printf("Href: %s, Text: %s\n", link.Href, link.Text)
	}
}
