package utils

import (
	"net/url"
	"regexp"
	"strings"
)

var urlRegex = regexp.MustCompile(`^[a-zA-Z0-9+.-]+:.*$`)

func IsValidURL(rawURL string) bool {
	return urlRegex.MatchString(rawURL)
}

func ParseURL(rawURL string) (*url.URL, error) {
	if !strings.HasPrefix(rawURL, "http://") && !strings.HasPrefix(rawURL, "https://") {
		rawURL = "https://" + rawURL
	}

	return url.Parse(rawURL)
}

func ExtractFileName(rawURL string) string {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "download"
	}
	path := u.Path
	if idx := strings.LastIndex(path, "/"); idx != -1 {
		path = path[idx+1:]
	}
	if path == "" {
		return "download"
	}
	return path
}

func GetBaseURL(rawURL string) string {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return rawURL
	}

	return parsed.Scheme + "://" + parsed.Host
}

func ResolveURL(base, relative string) string {
	baseParsed, err := url.Parse(base)
	if err != nil {
		return relative
	}

	relativeParsed, err := url.Parse(relative)
	if err != nil {
		return relative
	}

	return baseParsed.ResolveReference(relativeParsed).String()
}

func IsSameDomain(url1, url2 string) bool {
	parsed1, err1 := url.Parse(url1)
	parsed2, err2 := url.Parse(url2)

	if err1 != nil || err2 != nil {
		return false
	}

	return parsed1.Host == parsed2.Host
}

func GetHostname(rawURL string) string {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return ""
	}
	return parsed.Host
}
