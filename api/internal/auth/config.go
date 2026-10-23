package auth

import (
	"net/http"
	"os"
	"strings"

	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/config"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var oauthCfg *oauth2.Config

type cookieConfig struct {
	SameSite    http.SameSite
	Secure      bool
	Partitioned bool
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

	sameSiteEnv := strings.ToLower(os.Getenv("COOKIE_SAMESITE"))
	switch sameSiteEnv {
	case "none":
		// Required for cross-site deployments (e.g. frontend on GitHub Pages, API on Render)
		CookieCfg.SameSite = http.SameSiteNoneMode
		CookieCfg.Secure = true
		CookieCfg.Partitioned = true
	case "lax":
		CookieCfg.SameSite = http.SameSiteLaxMode
		CookieCfg.Secure = os.Getenv("APP_ENV") == "prod" || os.Getenv("COOKIE_SECURE") == "true"
		CookieCfg.Partitioned = false
	case "strict":
		CookieCfg.SameSite = http.SameSiteStrictMode
		CookieCfg.Secure = os.Getenv("APP_ENV") == "prod" || os.Getenv("COOKIE_SECURE") == "true"
		CookieCfg.Partitioned = false
	default:
		// Default behavior: Strict in prod (for same-domain/proxy setups like Caddy), Default in dev
		if os.Getenv("APP_ENV") == "prod" {
			CookieCfg.SameSite = http.SameSiteStrictMode
			CookieCfg.Secure = true
			CookieCfg.Partitioned = false
		} else {
			CookieCfg.SameSite = http.SameSiteDefaultMode
			CookieCfg.Secure = false
			CookieCfg.Partitioned = false
		}
	}
}
