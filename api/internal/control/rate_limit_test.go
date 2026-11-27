package control

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRateLimitSource(t *testing.T) {
	for _, check := range []struct {
		name           string
		onRender       bool
		providedSecret string
		wantIP         string
	}{
		{"local ignores forwarding headers", false, "wrong", "192.0.2.30"},
		{"direct Render uses ingress address", true, "wrong", "192.0.2.20"},
		{"Pages uses authenticated visitor address", true, "trusted", "192.0.2.10"},
	} {
		t.Run(check.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "/me", nil)
			request.RemoteAddr = "192.0.2.30:1234"
			request.Header.Set("CF-Connecting-IP", "192.0.2.20")
			request.Header.Set("X-Forwarded-For", "192.0.2.99")
			request.Header.Set("X-Rate-Limit-Client-IP", "::ffff:192.0.2.10")
			request.Header.Set("X-Rate-Limit-Proxy-Secret", check.providedSecret)
			source := rateLimitSource{onRender: check.onRender, proxySecret: "trusted"}
			if ip := source.clientIP(request).String(); ip != check.wantIP {
				t.Fatalf("client IP = %s, want %s", ip, check.wantIP)
			}
		})
	}
}

func TestRateLimitHTTPResponseAndIPv6Buckets(t *testing.T) {
	limiter := rateLimitSource{}.middleware(0.01, 1, nil)
	endpoint := limiter(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	for _, check := range []struct {
		remoteAddr string
		status     int
	}{
		{"[2001:db8:1::1]:1234", http.StatusNoContent},
		{"[2001:db8:1::2]:5678", http.StatusTooManyRequests},
		{"[2001:db8:2::1]:1234", http.StatusNoContent},
	} {
		request := httptest.NewRequest(http.MethodGet, "/me", nil)
		request.RemoteAddr = check.remoteAddr
		response := httptest.NewRecorder()
		endpoint.ServeHTTP(response, request)
		if response.Code != check.status {
			t.Fatalf("request from %s returned %d, want %d", check.remoteAddr, response.Code, check.status)
		}
		if check.status == http.StatusTooManyRequests {
			if response.Header().Get("Retry-After") == "" || response.Header().Get("Cache-Control") != "no-store" {
				t.Fatalf("missing rate limit headers: %v", response.Header())
			}
		}
	}
}
