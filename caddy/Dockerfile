# syntax=docker/dockerfile:1
FROM caddy AS caddy

COPY --from=mnnit-chess-club-frontend /app/build /usr/share/caddy
COPY Caddyfile /etc/caddy/Caddyfile

CMD ["caddy", "run", "--config", "/etc/caddy/Caddyfile", "--adapter", "caddyfile"]