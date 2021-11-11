package cmd

import (
	"github.com/odpf/siren/app"
	"github.com/odpf/siren/config"
	"github.com/spf13/cobra"
)

func serveCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "serve",
		Short: "Run server",
		RunE: func(cmd *cobra.Command, args []string) error {
			c := config.LoadConfig()
			return app.RunServer(c)
		},
	}
}
