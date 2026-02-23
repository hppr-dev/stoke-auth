module stoke

go 1.24.0

require (
	entgo.io/contrib v0.3.5
	entgo.io/ent v0.13.1
	github.com/ghodss/yaml v1.0.0
	github.com/go-faster/errors v0.7.1
	github.com/go-faster/jx v1.1.0
	github.com/go-ldap/ldap/v3 v3.4.7
	github.com/go-sql-driver/mysql v1.8.1
	github.com/golang-jwt/jwt/v5 v5.2.2
	github.com/lib/pq v1.10.9
	github.com/mattn/go-sqlite3 v1.14.19
	github.com/ogen-go/ogen v1.0.0
	github.com/prometheus/client_golang v1.19.0
	github.com/rs/zerolog v1.32.0
	github.com/vincentfree/opentelemetry/otelzerolog v0.0.10
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.49.0
	go.opentelemetry.io/otel v1.25.0
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc v1.25.0
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp v1.25.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.24.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp v1.25.0
	go.opentelemetry.io/otel/exporters/prometheus v0.46.0
	go.opentelemetry.io/otel/metric v1.25.0
	go.opentelemetry.io/otel/sdk v1.25.0
	go.opentelemetry.io/otel/sdk/metric v1.25.0
	go.opentelemetry.io/otel/trace v1.25.0
	go.uber.org/multierr v1.11.0
	golang.org/x/crypto v0.45.0
	hppr.dev/stoke v0.0.0
)

replace entgo.io/contrib v0.3.5 => github.com/swalkerhppr/ent-contrib v0.0.0-20240415204700-185b78ea4766

replace ariga.io/ogent v0.0.0-20230621041143-ed3e5d4da458 => github.com/swalkerhppr/ogent v0.0.0-20240416025651-87740ceb514a

replace hppr.dev/stoke v0.0.0 => ./client/stoke

require (
	ariga.io/atlas v0.19.1-0.20240203083654-5948b60a8e43 // indirect
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/Azure/go-ntlmssp v0.0.0-20221128193559-754e69321358 // indirect
	github.com/agext/levenshtein v1.2.3 // indirect
	github.com/apparentlymart/go-textseg/v13 v13.0.0 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cenkalti/backoff/v4 v4.3.0 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/dlclark/regexp2 v1.11.0 // indirect
	github.com/fatih/color v1.16.0 // indirect
	github.com/felixge/httpsnoop v1.0.4 // indirect
	github.com/go-asn1-ber/asn1-ber v1.5.5 // indirect
	github.com/go-faster/yaml v0.4.6 // indirect
	github.com/go-logr/logr v1.4.1 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-openapi/inflect v0.19.0 // indirect
	github.com/google/go-cmp v0.6.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.19.0 // indirect
	github.com/hashicorp/hcl/v2 v2.15.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mitchellh/go-wordwrap v1.0.1 // indirect
	github.com/prometheus/client_model v0.6.1 // indirect
	github.com/prometheus/common v0.48.0 // indirect
	github.com/prometheus/procfs v0.12.0 // indirect
	github.com/rs/xid v1.5.0 // indirect
	github.com/segmentio/asm v1.2.0 // indirect
	github.com/sergi/go-diff v1.1.0 // indirect
	github.com/stoewer/go-strcase v1.2.0 // indirect
	github.com/zclconf/go-cty v1.12.1 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.25.0 // indirect
	go.opentelemetry.io/proto/otlp v1.1.0 // indirect
	go.uber.org/zap v1.27.0 // indirect
	golang.org/x/exp v0.0.0-20230725093048-515e97ebf090 // indirect
	golang.org/x/mod v0.29.0 // indirect
	golang.org/x/net v0.47.0 // indirect
	golang.org/x/sync v0.18.0 // indirect
	golang.org/x/sys v0.38.0 // indirect
	golang.org/x/text v0.31.0 // indirect
	golang.org/x/tools v0.38.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20240318140521-94a12d6c2237 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240401170217-c3f982113cda // indirect
	google.golang.org/grpc v1.64.1 // indirect
	google.golang.org/protobuf v1.33.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)
