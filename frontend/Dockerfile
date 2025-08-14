# syntax=docker/dockerfile:1
FROM node AS builder
WORKDIR /app

ARG PUBLIC_BASE_URL
ARG PUBLIC_FRONTEND_URL
ARG PUBLIC_WS_URL

RUN corepack enable && corepack install -g pnpm
COPY package.json pnpm-lock.yaml ./
RUN pnpm install --frozen-lockfile
COPY . .
RUN pnpm build

FROM alpine as output
WORKDIR /app
COPY --from=builder /app/build ./build