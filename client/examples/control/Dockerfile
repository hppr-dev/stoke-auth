FROM golang:1.22.2-alpine3.19 AS go-builder
RUN mkdir -p /build

WORKDIR /build

COPY . ./
RUN apk add build-base

WORKDIR /build/examples/control
RUN go mod tidy && go build -o ./ship-control .

FROM alpine:3.19

COPY --from=go-builder /build/examples/control/ship-control /usr/bin/ship-control

EXPOSE 8080

ENTRYPOINT ["/usr/bin/ship-control"]
