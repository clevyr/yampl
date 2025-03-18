#syntax=docker/dockerfile:1

FROM --platform=$BUILDPLATFORM tonistiigi/xx:1.6.1 AS xx

FROM --platform=$BUILDPLATFORM golang:1.24.1-alpine AS builder
WORKDIR /app

COPY --from=xx / /

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Set Golang build envs based on Docker platform string
ARG TARGETPLATFORM
RUN --mount=type=cache,target=/root/.cache \
  CGO_ENABLED=0 xx-go build -ldflags='-w -s'


FROM alpine:3.21
LABEL org.opencontainers.image.source="https://github.com/clevyr/yampl"
WORKDIR /data

RUN apk add --no-cache git jq yq-go

COPY --from=builder /app/yampl /usr/local/bin/

ENTRYPOINT ["yampl"]
