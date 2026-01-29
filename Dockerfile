#syntax=docker/dockerfile:1

FROM --platform=$BUILDPLATFORM golang:1.25.6-alpine AS builder
WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG TARGETOS
ARG TARGETARCH
RUN --mount=type=cache,target=/root/.cache \
  CGO_ENABLED=0 GOOS="$TARGETOS" GOARCH="$TARGETARCH" \
  go build -ldflags='-w -s'


FROM alpine:3.23
LABEL org.opencontainers.image.source="https://github.com/clevyr/yampl"
WORKDIR /data

RUN apk add --no-cache git jq yq-go

COPY --from=builder /app/yampl /usr/local/bin/

ENTRYPOINT ["yampl"]
