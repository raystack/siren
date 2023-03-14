package cli

import (
	"context"
	"time"

	"github.com/goto/salt/cmdx"
	"github.com/goto/salt/config"
	"github.com/goto/siren/pkg/errors"
	sirenv1beta1 "github.com/goto/siren/proto/gotocompany/siren/v1beta1"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ClientConfig struct {
	Host string `yaml:"host" cmdx:"host" default:"localhost:8080"`
}

func loadClientConfig(cmd *cobra.Command, cmdxConfig *cmdx.Config) (*ClientConfig, error) {
	var clientConfig ClientConfig

	if err := cmdxConfig.Load(
		&clientConfig,
		cmdx.WithFlags(cmd.Flags()),
	); err != nil {
		if !errors.As(err, new(config.ConfigFileNotFoundError)) {
			return nil, err
		}
	}

	if err := validateClientConfig(&clientConfig); err != nil {
		return nil, err
	}

	return &clientConfig, nil
}

func validateClientConfig(cfg *ClientConfig) error {
	if cfg.Host == "" {
		return errors.New("`host` is missing")
	}
	return nil
}

func createConnection(ctx context.Context, host string) (*grpc.ClientConn, error) {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	}

	return grpc.DialContext(ctx, host, opts...)
}

func createClient(ctx context.Context, host string) (sirenv1beta1.SirenServiceClient, func(), error) {
	dialTimeoutCtx, dialCancel := context.WithTimeout(ctx, time.Second*2)
	conn, err := createConnection(dialTimeoutCtx, host)
	if err != nil {
		dialCancel()
		return nil, nil, err
	}

	cancel := func() {
		dialCancel()
		conn.Close()
	}

	client := sirenv1beta1.NewSirenServiceClient(conn)
	return client, cancel, nil
}
