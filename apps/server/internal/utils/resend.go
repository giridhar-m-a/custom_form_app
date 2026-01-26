package utils

import (
	"fmt"

	"github.com/resend/resend-go/v3"
)

var ResendClient *resend.Client

func InitResend() {
	apiKey := GetEnv("API_KEY", "")
	if apiKey == "" {
		fmt.Println("Resend Credentials not found")
	}

	ResendClient = resend.NewClient(apiKey)
	fmt.Println("Mail client initialised successfully")
}
