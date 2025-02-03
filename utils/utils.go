package utils

import (
	"net/url"
	"strings"
)

func CheckURLisValid(urlString string) bool {
	_, err := url.ParseRequestURI(urlString)

	if err != nil {
		return false
	}

	return true
}

func GetHostFromURL(urlString string) (string) {
	var hostname string

	// check it's all lower
	urlString = strings.ToLower(urlString)

	// strip common prefixes
	if strings.HasPrefix(urlString, "https://www.") {
		// strip prefix
		hostname = urlString[12:]
	} else if strings.HasPrefix(urlString, "http://www.") {
		// strip prefix
		hostname = urlString[11:]
	} else if strings.HasPrefix(urlString, "https://") {
		// strip prefix
		hostname = urlString[8:]
	} else if strings.HasPrefix(urlString, "http://") {
		// strip prefix
		hostname = urlString[7:]
	} else if strings.HasPrefix(urlString, "www.") {
		// strip prefix
		hostname = urlString[4:]
	} else {
		hostname = urlString
	}

	// strip english prefix
	hostname = strings.TrimPrefix(hostname, "en.")

	// strip mobile prefixes
	hostname = strings.TrimPrefix(hostname, "m.")
	hostname = strings.TrimPrefix(hostname, "mobile.")

	// strip uk prefixes
	hostname = strings.TrimPrefix(hostname, "uk.")

	if strings.Contains(hostname, "/") {
		// strip path
		hostname = hostname[:strings.Index(hostname, "/")]
	}

	return hostname
}
