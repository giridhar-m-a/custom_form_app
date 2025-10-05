package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func GetUserDetails(accessToken string)(map[string]any, error) {
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