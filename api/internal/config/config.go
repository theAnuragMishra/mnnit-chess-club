package config

import (
	"net/url"
	"os"
	"strings"
)

var (
	BaseURL        string
	FrontendURL    string
	AllowedOrigins []string
)

func Config() {
	BaseURL = os.Getenv("BASE_URL")
	FrontendURL = os.Getenv("FRONTEND_URL")

	originsMap := map[string]bool{
		"http://localhost:*":    true,
		"https://localhost:*":   true,
		"http://localhost:5173":  true,
		"http://localhost:4173":  true,
		"http://localhost:3000":  true,
	}

	// Extract origin from FrontendURL (e.g. "https://user.github.io/repo" -> "https://user.github.io")
	if FrontendURL != "" {
		if parsed, err := url.Parse(FrontendURL); err == nil && parsed.Scheme != "" && parsed.Host != "" {
			origin := parsed.Scheme + "://" + parsed.Host
			originsMap[origin] = true
		}
		trimmed := strings.TrimRight(FrontendURL, "/")
		if trimmed != "" {
			originsMap[trimmed] = true
		}
	}

	// Additional origins from CORS_ALLOWED_ORIGINS (comma-separated)
	if corsEnv := os.Getenv("CORS_ALLOWED_ORIGINS"); corsEnv != "" {
		for _, o := range strings.Split(corsEnv, ",") {
			o = strings.TrimSpace(o)
			if o != "" {
				originsMap[o] = true
				if parsed, err := url.Parse(o); err == nil && parsed.Scheme != "" && parsed.Host != "" {
					originsMap[parsed.Scheme+"://"+parsed.Host] = true
				}
			}
		}
	}

	AllowedOrigins = make([]string, 0, len(originsMap))
	for o := range originsMap {
		AllowedOrigins = append(AllowedOrigins, o)
	}
}

// IsAllowedOrigin checks if a given origin is allowed (used for WebSocket origin check)
func IsAllowedOrigin(origin string) bool {
	if origin == "" {
		return true
	}

	trimmedOrigin := strings.TrimRight(origin, "/")

	// Also extract scheme://host in case the origin header contains anything unusual
	parsedOrigin := ""
	if parsed, err := url.Parse(trimmedOrigin); err == nil && parsed.Scheme != "" && parsed.Host != "" {
		parsedOrigin = parsed.Scheme + "://" + parsed.Host
	}

	for _, allowed := range AllowedOrigins {
		if allowed == "*" || allowed == trimmedOrigin || (parsedOrigin != "" && allowed == parsedOrigin) {
			return true
		}
		// Match localhost port wildcards like "http://localhost:*"
		if strings.HasSuffix(allowed, ":*") {
			prefix := strings.TrimSuffix(allowed, ":*")
			if strings.HasPrefix(trimmedOrigin, prefix+":") || trimmedOrigin == prefix {
				return true
			}
		}
	}
	return false
}
