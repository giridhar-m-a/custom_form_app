package services

import (
	"context"
	"log"

	"github.com/giridhar-m-a/custom_form_app/configs"
	"golang.org/x/oauth2"
)

func GoogleAuthService(code string) (*oauth2.Token, error) {
    token, err := configs.GoogleOauthConfig.Exchange(context.Background(), code)
    if err != nil {
        log.Printf("Error exchanging code for token: %v", err)
        return nil, err
    }
    return token, nil
}

