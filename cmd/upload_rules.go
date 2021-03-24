package cmd

import (
	"github.com/odpf/siren/app"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "upload:rules",
		Short: "UploadRule",
		RunE:  uploadRules,
	})
}

func uploadRules(cmd *cobra.Command, args []string) error {
	c := app.LoadConfig()
	return app.UploadRules(c, args[0], args[1])
}
