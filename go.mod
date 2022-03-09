module github.com/odpf/siren

go 1.16

require (
	github.com/DATA-DOG/go-sqlmock v1.5.0
	github.com/MakeNowJust/heredoc v1.0.0
	github.com/envoyproxy/protoc-gen-validate v0.6.2
	github.com/go-openapi/loads v0.20.1 // indirect
	github.com/go-openapi/runtime v0.19.26
	github.com/go-openapi/spec v0.20.2 // indirect
	github.com/go-playground/universal-translator v0.18.0 // indirect
	github.com/gofrs/uuid v4.0.0+incompatible // indirect
	github.com/golang/snappy v0.0.3 // indirect
	github.com/grafana/cortex-tools v0.7.2
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.6.0
	github.com/gtank/cryptopasta v0.0.0-20170601214702-1f550f6f2f69
	github.com/kr/pretty v0.3.0 // indirect
	github.com/leodido/go-urn v1.2.1 // indirect
	github.com/lib/pq v1.3.0
	github.com/mitchellh/mapstructure v1.4.1
	github.com/newrelic/go-agent/v3 v3.12.0
	github.com/newrelic/go-agent/v3/integrations/nrgrpc v1.3.1
	github.com/odpf/salt v0.0.0-20220106155451-62e8c849ae81
	github.com/pkg/errors v0.9.1
	github.com/prometheus/alertmanager v0.21.1-0.20200911160112-1fdff6b3f939
	github.com/prometheus/prometheus v1.8.2-0.20201014093524-73e2ce1bd643
	github.com/rogpeppe/go-internal v1.8.0 // indirect
	github.com/slack-go/slack v0.9.3
	github.com/spf13/cobra v1.2.1
	github.com/stretchr/testify v1.7.0
	go.uber.org/multierr v1.7.0 // indirect
	go.uber.org/zap v1.19.0
	golang.org/x/crypto v0.0.0-20210711020723-a769d52b0f97 // indirect
	golang.org/x/net v0.0.0-20210825183410-e898025ed96a
	golang.org/x/sys v0.0.0-20210823070655-63515b42dcdf // indirect
	google.golang.org/genproto v0.0.0-20210903162649-d08c68adba83
	google.golang.org/grpc v1.40.0
	google.golang.org/protobuf v1.27.1
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/go-playground/assert.v1 v1.2.1 // indirect
	gopkg.in/go-playground/validator.v9 v9.31.0
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
	gorm.io/driver/postgres v1.0.8
	gorm.io/gorm v1.20.12
)

replace k8s.io/client-go => k8s.io/client-go v0.19.2

replace github.com/grafana/cortex-tools v0.7.2 => github.com/kevinbheda/cortex-tools v0.8.0
