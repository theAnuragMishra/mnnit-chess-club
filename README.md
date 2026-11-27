# MNNIT Chess Club's official website

[mnnitchessclub.pages.dev](https://mnnitchessclub.pages.dev)

> [!NOTE]
> The deployment at [mnnitchessclub.pages.dev](https://mnnitchessclub.pages.dev) is for purposes of project demonstration and casual play. It's still as secure as the production deployment I ran at mnnitchessclub.com, but now that I have moved the project to free tiers requests may be slow and there may be cold starts. Also since the database itself is in free tier too, I make no guarantee of data availability and your games or even your account may disappear at any time. For casual play or testing my project though, it's all working!

## Features:

1. PvP gameplay
2. Random Matchmaking
3. Challenge a person
4. Arena style tournaments

### Planned:

1. Friends & Followers

## Development

A high level design of the architecture (well, not very well-drawn) can be found [here](https://mnnitchessclub.pages.dev/dev).

### Stack

1. Frontend: Svelte
2. Server : Golang
3. Database: Postgres
4. Deployment: Docker & compose on a vps

### API rate limiting

The Go API uses in-memory token buckets. Limits reset when the Render service restarts or spins down.

| Traffic | Sustained limit | Burst |
| --- | --- | --- |
| All HTTP requests, including WebSocket handshakes | 40 per second across the API | 120 |
| HTTP requests per client IPv4 address or IPv6 /64 | 8 per second | 60 |
| Google OAuth login and callback combined, per client | 1 per second | 20 |
| WebSocket messages per signed-in account, across tabs | 10 per second | 25 |

HTTP requests over budget return `429` with `Retry-After`. WebSockets close with code `1008` if the account exceeds its message budget. Incoming WebSocket messages can be up to 16 KiB.

Set the same `RATE_LIMIT_PROXY_SECRET` on the Render API and as a Cloudflare Pages secret so the API can identify individual visitors behind the Pages proxy. Generate a value with `openssl rand -hex 32`. Without it, traffic through Pages can share the proxy's IP bucket. Direct requests to Render use the IP supplied by Render's Cloudflare ingress.

## Contributing

Hey! I will be really happy to have people contributing to the development of the site. You can contribute in any of the following ways:

- Create issues for bugs! Yes that's the simplest way to contribute. We don't want to let a single bug go free XD
- Add features. Yes, do it. I love reviewing pull requests.
- Just spread the word. We will be so happy to have more users play on our site.
