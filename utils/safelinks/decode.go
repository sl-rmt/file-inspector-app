package safelinks

import (
	"net/url"
	"strings"
)

const (
	https = "https://"
	slURL = "safelinks.protection.outlook.com/?url="
)

// they have a region-specific prefix such as:
// https://gbr01.safelinks.protection.outlook.com/?url=
func IsSafelink(urlString string) bool {
	if strings.HasPrefix(urlString, https) && strings.Contains(urlString, slURL) {
		return true
	}

	return false
}

func ExtractOriginalURL(urlString string) (string, error) {
	// get the encoded URL from the string
	startIndex := strings.Index(urlString, slURL) + len(slURL)
	encodedURL := urlString[startIndex:]
	
	// deode it
	decodedURL, err := url.QueryUnescape(encodedURL)

	if err != nil {
		return "", err
	}

	return decodedURL, nil
}
