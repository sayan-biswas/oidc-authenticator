FROM alpine:latest

LABEL maintainer="Sayan Biswas"
LABEL version="1.0.0"

COPY oidc-service /usr/bin
COPY 