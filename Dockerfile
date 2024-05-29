FROM node:22-alpine3.19 AS node-builder
RUN mkdir -p /build/dist
WORKDIR /build

COPY internal/admin/stoke-admin-ui ./
RUN npm install && npm run build --emptyOutDir

FROM golang:1.22.2-alpine3.19 AS go-builder
RUN mkdir -p /build
WORKDir /build

COPY . ./
COPY --from=node-builder /build/dist ./internal/admin/
RUN apk add build-base

ARG EXTRA_BUILD_ARGS
RUN go mod tidy && go build -o ./stoke-server $EXTRA_BUILD_ARGS ./cmd/

FROM alpine:3.19
RUN mkdir /etc/stoke /var/stoke/log -p
RUN addgroup -S stoke && adduser -HDS stoke stoke
RUN chown stoke:stoke -R /var/stoke /etc/stoke

COPY --from=go-builder /build/stoke-server /usr/bin/stoke-server

USER stoke
ENTRYPOINT ["/usr/bin/stoke-server", "-config", "/etc/stoke/config.yaml"]
