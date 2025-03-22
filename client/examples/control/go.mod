module control

go 1.22.2

replace hppr.dev/stoke v0.0.0 => ../../stoke

replace engine v0.0.0 => ../engine

require (
	engine v0.0.0
	github.com/golang-jwt/jwt/v5 v5.2.2
	golang.org/x/net v0.26.0
	google.golang.org/grpc v1.64.1
	hppr.dev/stoke v0.0.0
)

require (
	github.com/go-faster/errors v0.7.1 // indirect
	github.com/go-faster/jx v1.1.0 // indirect
	github.com/go-logr/logr v1.4.1 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/segmentio/asm v1.2.0 // indirect
	go.opentelemetry.io/otel v1.25.0 // indirect
	go.opentelemetry.io/otel/metric v1.25.0 // indirect
	go.opentelemetry.io/otel/trace v1.25.0 // indirect
	golang.org/x/sys v0.21.0 // indirect
	golang.org/x/text v0.16.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240318140521-94a12d6c2237 // indirect
	google.golang.org/protobuf v1.34.1 // indirect
)
