package auth

import (
	"fmt"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	api_header := headers.Get("Authorization")
	apiSlices := strings.Split(api_header, " ")
	if len(apiSlices) < 2 {
		return "", fmt.Errorf("No API Header")
	}

	return apiSlices[1], nil
}
