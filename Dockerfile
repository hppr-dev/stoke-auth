FROM node:24.13.1-alpine3.22 AS node-builder
RUN mkdir -p /build/
WORKDIR /build

COPY internal/admin/ .
RUN cd stoke-admin-ui && npm install && npm run build

FROM golang:1.25.7-alpine3.23 AS go-builder
RUN apk add build-base

RUN mkdir -p /build
WORKDIR /build

COPY . ./
COPY --from=node-builder /build/stoke-admin-ui/dist ./internal/admin/stoke-admin-ui/dist

ARG EXTRA_BUILD_ARGS
RUN go mod tidy
RUN go build -o ./stoke-server $EXTRA_BUILD_ARGS ./cmd/

FROM alpine:3.23
RUN mkdir /etc/stoke /var/stoke/log -p
RUN addgroup -S stoke -g 1000 && adduser -HDS -u 1000 stoke stoke
RUN chown stoke:stoke -R /var/stoke /etc/stoke

COPY --from=go-builder /build/stoke-server /usr/bin/stoke-server

USER stoke
ENTRYPOINT ["/usr/bin/stoke-server", "-config", "/etc/stoke/config.yaml"]
