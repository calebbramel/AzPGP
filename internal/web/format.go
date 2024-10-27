package web

import (
	"fmt"
	"net/url"
	"regexp"
)

func sanitizeString(input string) string {
	// regex for all non-alphanumeric characters except dashes
	re := regexp.MustCompile(`[^a-zA-Z0-9-]+`)
	// Replace matches with a dash
	return re.ReplaceAllString(input, "-")
}

func urlDecode(encodedString string) (string, error) {
	decodedString, err := url.QueryUnescape(encodedString)
	if err != nil {
		return "", fmt.Errorf("error decoding email: %v", err)
	}
	return decodedString, nil
}
