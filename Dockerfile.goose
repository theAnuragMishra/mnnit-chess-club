FROM alpine:3.21.2

RUN apk update && \
    apk upgrade && \
    apk add bash && \
    rm -rf /var/cache/apk/*

ADD https://github.com/pressly/goose/releases/download/v3.24.3/goose_linux_x86_64 /bin/goose
RUN chmod +x /bin/goose

WORKDIR /migrations