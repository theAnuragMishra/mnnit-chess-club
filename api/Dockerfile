# syntax=docker/dockerfile:1
FROM golang:1.24 as build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /api ./cmd

FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

COPY --from=build-stage /api /api

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/api"]