package auth

import (
	"net/http"
	"os"

	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/config"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var oauthCfg *oauth2.Config

type cookieConfig struct {
	SameSite http.SameSite
	Secure   bool
}

var CookieCfg cookieConfig

func Config() {

	oauthCfg = &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  config.BaseURL + "/auth/callback/google",
		Scopes:       []string{"email", "profile", "openid"},
		Endpoint:     google.Endpoint,
	}
	if os.Getenv("APP_ENV") == "prod" {
		CookieCfg.SameSite = http.SameSiteLaxMode
		CookieCfg.Secure = true
	} else {
		CookieCfg.SameSite = http.SameSiteDefaultMode
		CookieCfg.Secure = false
	}

}
