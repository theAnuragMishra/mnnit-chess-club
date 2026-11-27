package control

import (
	"crypto/subtle"
	"net"
	"net/http"
	"net/netip"
	"strconv"
	"time"

	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/ratelimit"
	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/utils"
)

const maxRateLimitIPs = 4096

type rateLimitSource struct {
	onRender    bool
	proxySecret string
}

func (source rateLimitSource) clientIP(r *http.Request) netip.Addr {
	// The Pages function supplies its visitor IP on requests it authenticates.
	if source.proxySecret != "" && subtle.ConstantTimeCompare(
		[]byte(source.proxySecret), []byte(r.Header.Get("X-Rate-Limit-Proxy-Secret")),
	) == 1 {
		if addr, err := netip.ParseAddr(r.Header.Get("X-Rate-Limit-Client-IP")); err == nil {
			return addr.Unmap()
		}
	}
	// Public ingress on Render's free tier passes through Cloudflare.
	if source.onRender {
		if addr, err := netip.ParseAddr(r.Header.Get("CF-Connecting-IP")); err == nil {
			return addr.Unmap()
		}
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		host = r.RemoteAddr
	}
	addr, _ := netip.ParseAddr(host)
	return addr.Unmap()
}

func (source rateLimitSource) middleware(perSecond float64, burst int, overall *ratelimit.Limiter[struct{}]) func(http.Handler) http.Handler {
	clients := ratelimit.New[netip.Addr](perSecond, burst, maxRateLimitIPs)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := source.clientIP(r)
			if ip.Is6() {
				// Group IPv6 privacy addresses from the same network.
				ip = netip.PrefixFrom(ip, 64).Masked().Addr()
			}
			retry := clients.Take(ip)
			if retry == 0 && overall != nil {
				retry = overall.Take(struct{}{})
			}
			if retry > 0 {
				seconds := (retry + time.Second - 1) / time.Second
				w.Header().Set("Retry-After", strconv.FormatInt(int64(seconds), 10))
				w.Header().Set("Cache-Control", "no-store")
				utils.RespondWithError(w, http.StatusTooManyRequests, "Too many requests. Try again shortly.")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
