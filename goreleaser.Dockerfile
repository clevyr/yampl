FROM alpine:3.24
WORKDIR /data

RUN apk add --no-cache git jq yq

ARG TARGETPLATFORM
COPY $TARGETPLATFORM/yampl /usr/local/bin

ENTRYPOINT ["yampl"]
