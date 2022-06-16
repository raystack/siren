module github.com/odpf/siren

go 1.16

require (
	github.com/DATA-DOG/go-sqlmock v1.5.0
	github.com/MakeNowJust/heredoc v1.0.0
	github.com/benbjohnson/clock v1.3.0 // indirect
	github.com/envoyproxy/protoc-gen-validate v0.6.7
	github.com/go-openapi/loads v0.20.1 // indirect
	github.com/go-openapi/runtime v0.19.26
	github.com/go-openapi/spec v0.20.2 // indirect
	github.com/go-playground/universal-translator v0.18.0 // indirect
	github.com/gofrs/uuid v4.2.0+incompatible // indirect
	github.com/google/go-cmp v0.5.8
	github.com/grafana/cortex-tools v0.7.2
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.10.3
	github.com/gtank/cryptopasta v0.0.0-20170601214702-1f550f6f2f69
	github.com/leodido/go-urn v1.2.1 // indirect
	github.com/lib/pq v1.3.0
	github.com/mitchellh/mapstructure v1.5.0
	github.com/newrelic/go-agent/v3 v3.12.0
	github.com/newrelic/go-agent/v3/integrations/nrgrpc v1.3.1
	github.com/odpf/salt v0.0.0-20220106155451-62e8c849ae81
	github.com/pelletier/go-toml/v2 v2.0.2 // indirect
	github.com/pkg/errors v0.9.1
	github.com/prometheus/alertmanager v0.21.1-0.20200911160112-1fdff6b3f939
	github.com/prometheus/prometheus v1.8.2-0.20201014093524-73e2ce1bd643
	github.com/rogpeppe/go-internal v1.8.0 // indirect
	github.com/slack-go/slack v0.11.0
	github.com/spf13/cobra v1.4.0
	github.com/spf13/viper v1.12.0 // indirect
	github.com/stretchr/testify v1.7.2
	github.com/subosito/gotenv v1.4.0 // indirect
	go.buf.build/odpf/gw/odpf/proton v1.1.97
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/multierr v1.8.0 // indirect
	go.uber.org/zap v1.21.0
	golang.org/x/crypto v0.0.0-20220525230936-793ad666bf5e // indirect
	golang.org/x/net v0.0.0-20220526153639-5463443f8c37
	golang.org/x/sys v0.0.0-20220610221304-9f5ed59c137d // indirect
	golang.org/x/tools v0.1.11 // indirect
	google.golang.org/genproto v0.0.0-20220525015930-6ca3db687a9d
	google.golang.org/grpc v1.47.0
	google.golang.org/protobuf v1.28.0
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/go-playground/assert.v1 v1.2.1 // indirect
	gopkg.in/go-playground/validator.v9 v9.31.0
	gopkg.in/ini.v1 v1.66.6 // indirect
	gopkg.in/yaml.v3 v3.0.1
	gorm.io/driver/postgres v1.0.8
	gorm.io/gorm v1.20.12
)

replace k8s.io/client-go => k8s.io/client-go v0.19.2

replace github.com/grafana/cortex-tools v0.7.2 => github.com/kevinbheda/cortex-tools v0.8.0
