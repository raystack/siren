package cmd

import (
	"github.com/odpf/siren/app"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "upload",
		Short: "Upload",
		RunE:  uploadTemplate,
	})
}

func uploadTemplate(cmd *cobra.Command, args []string) error {
	c := app.LoadConfig()
	return app.UploadTemplates(c, args[0])
}
