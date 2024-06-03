FROM golang:1.22.2-alpine3.19 AS go-builder
RUN mkdir -p /build

WORKDIR /build

COPY . ./
RUN apk add build-base

WORKDIR /build/examples/engine
RUN go mod tidy && go build -o ./engine-room .

FROM alpine:3.19

COPY --from=go-builder /build/examples/control/engine-room /usr/bin/engine-room

EXPOSE 6060

ENTRYPOINT ["/usr/bin/engine-room"]
