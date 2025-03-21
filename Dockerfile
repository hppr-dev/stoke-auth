FROM node:22-alpine3.19 AS node-builder
RUN mkdir -p /build/dist
WORKDIR /build

COPY internal/admin/stoke-admin-ui ./
RUN npm install && npm run build --emptyOutDir

FROM golang:1.23.7-alpine3.20 AS go-builder
RUN apk add build-base

RUN mkdir -p /build
WORKDir /build

COPY . ./
COPY --from=node-builder /build/dist ./internal/admin/

ARG EXTRA_BUILD_ARGS
RUN go mod tidy
RUN go build -o ./stoke-server $EXTRA_BUILD_ARGS ./cmd/

FROM alpine:3.20
RUN mkdir /etc/stoke /var/stoke/log -p
RUN addgroup -S stoke -g 1000 && adduser -HDS -u 1000 stoke stoke
RUN chown stoke:stoke -R /var/stoke /etc/stoke

COPY --from=go-builder /build/stoke-server /usr/bin/stoke-server

USER stoke
ENTRYPOINT ["/usr/bin/stoke-server", "-config", "/etc/stoke/config.yaml"]
