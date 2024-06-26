FROM alpine:3.20
LABEL org.opencontainers.image.source="https://github.com/clevyr/yampl"
WORKDIR /data

RUN apk add --no-cache git jq yq

COPY yampl /usr/local/bin

ENTRYPOINT ["yampl"]
