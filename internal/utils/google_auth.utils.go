package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// GetUserDetails retrieves user information from Google's OAuth2 userinfo endpoint using the provided access token.
// It returns a map containing the decoded JSON user claims, or an error if the HTTP request or JSON decoding fails.
func GetUserDetails(accessToken string) (map[string]any, error) {
	userInfoEndpoint := "https://www.googleapis.com/oauth2/v2/userinfo"
	resp, err := http.Get(fmt.Sprintf("%s?access_token=%s", userInfoEndpoint, accessToken))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var userInfo map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, err
	}

	return userInfo, nil
}