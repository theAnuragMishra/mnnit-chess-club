package config

import "os"

var BaseURL, FrontendURL string

func Config() {
	BaseURL = os.Getenv("BASE_URL")
	FrontendURL = os.Getenv("FRONTEND_URL")
}
