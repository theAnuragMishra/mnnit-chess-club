FROM alpine:3.21.2

RUN apk add --no-cache curl

RUN curl -fsSL -o /bin/goose https://github.com/pressly/goose/releases/download/v3.24.3/goose_linux_x86_64 \
    && chmod +x /bin/goose

WORKDIR /migrations
