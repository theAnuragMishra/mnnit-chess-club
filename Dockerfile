# syntax=docker/dockerfile:1
FROM node AS build-stage
WORKDIR /app

ARG PUBLIC_BASE_URL
ARG PUBLIC_FRONTEND_URL
ARG PUBLIC_WS_URL

RUN corepack enable && corepack install -g pnpm
COPY frontend/package.json frontend/pnpm-lock.yaml ./
RUN pnpm install --frozen-lockfile
COPY frontend .
RUN pnpm build

FROM caddy AS caddy

COPY --from=build-stage /app/build /usr/share/caddy
COPY Caddyfile /etc/caddy/Caddyfile

CMD ["caddy", "run", "--config", "/etc/caddy/Caddyfile", "--adapter", "caddyfile"]