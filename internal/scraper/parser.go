package scraper

import (
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

type Parser struct {
	baseURL string
}

func NewParser(baseURL string) *Parser {
	return &Parser{baseURL: baseURL}
}

func (p *Parser) ParseTitle(n *html.Node) string {
	if n.Type == html.ElementNode && n.Data == "title" {
		if n.FirstChild != nil {
			text := strings.TrimSpace(n.FirstChild.Data)
			if text != "" {
				return text
			}
		}
		return ""
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if title := p.ParseTitle(c); title != "" {
			return title
		}
	}

	return ""
}

func (p *Parser) ParseDescription(n *html.Node) string {
	if n.Type == html.ElementNode && n.Data == "meta" {
		name := getAttr(n, "name")
		if name == "description" {
			return getAttr(n, "content")
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if desc := p.ParseDescription(c); desc != "" {
			return desc
		}
	}

	return ""
}

func (p *Parser) ParseLinks(n *html.Node, links *[]string) {
	if n.Type == html.ElementNode && n.Data == "a" {
		href := getAttr(n, "href")
		if href != "" {
			*links = append(*links, p.resolveURL(href))
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		p.ParseLinks(c, links)
	}
}

func (p *Parser) ParseImages(n *html.Node, images *[]string) {
	if n.Type == html.ElementNode && n.Data == "img" {
		src := getAttr(n, "src")
		if src != "" {
			*images = append(*images, p.resolveURL(src))
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		p.ParseImages(c, images)
	}
}

func (p *Parser) ParseOGTags(n *html.Node, meta *Metadata) {
	if n.Type == html.ElementNode && n.Data == "meta" {
		property := getAttr(n, "property")
		content := getAttr(n, "content")

		switch property {
		case "og:title":
			meta.OGTitle = content
		case "og:description":
			meta.Description = content
		case "og:image":
			meta.OGImage = content
		case "og:url":
			meta.OGURL = content
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		p.ParseOGTags(c, meta)
	}
}

func (p *Parser) ParseStructuredData(n *html.Node, data *[]string) {
	if n.Type == html.ElementNode && n.Data == "script" {
		scriptType := getAttr(n, "type")
		if scriptType == "application/ld+json" {
			if n.FirstChild != nil {
				*data = append(*data, n.FirstChild.Data)
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		p.ParseStructuredData(c, data)
	}
}

func (p *Parser) resolveURL(href string) string {
	if strings.HasPrefix(href, "http://") || strings.HasPrefix(href, "https://") {
		return href
	}

	if strings.HasPrefix(href, "//") {
		return "https:" + href
	}

	if strings.HasPrefix(href, "/") {
		return strings.TrimSuffix(p.baseURL, "/") + href
	}

	return p.baseURL + "/" + href
}

func getAttr(n *html.Node, key string) string {
	for _, attr := range n.Attr {
		if attr.Key == key {
			return attr.Val
		}
	}
	return ""
}

var linkRegex = regexp.MustCompile(`href=["']([^"']+)["']`)

func ExtractLinksFromHTML(htmlContent string) []string {
	matches := linkRegex.FindAllStringSubmatch(htmlContent, -1)
	links := make([]string, 0, len(matches))
	for _, match := range matches {
		if len(match) > 1 {
			links = append(links, match[1])
		}
	}
	return links
}
