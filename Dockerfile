#syntax=docker/dockerfile:1.4

FROM --platform=$BUILDPLATFORM golang:1.20-alpine as builder
WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Set Golang build envs based on Docker platform string
ARG TARGETPLATFORM
RUN <<EOT
  set -eux
  case "$TARGETPLATFORM" in
    'linux/amd64') export GOARCH=amd64 ;;
    'linux/arm/v6') export GOARCH=arm GOARM=6 ;;
    'linux/arm/v7') export GOARCH=arm GOARM=7 ;;
    'linux/arm64') export GOARCH=arm64 ;;
    *) echo "Unsupported target: $TARGETPLATFORM" && exit 1 ;;
  esac
  go build -ldflags='-w -s' -o yampl
EOT


FROM alpine:3.17
LABEL org.opencontainers.image.source="https://github.com/clevyr/yampl"
WORKDIR /data

RUN apk add --no-cache git jq yq

COPY --from=builder /app/yampl /usr/local/bin/

ENTRYPOINT ["yampl"]
