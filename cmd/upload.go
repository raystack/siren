package cmd

import (
	"github.com/odpf/siren/app"
	"github.com/odpf/siren/uploader"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "upload",
		Short: "Upload Rules or Templates YAML file",
		RunE:  uploadRules,
	})
}

func uploadRules(cmd *cobra.Command, args []string) error {
	c := app.LoadConfig()
	s := uploader.NewService(&c.SirenService)
	return s.Upload(args[0])
}
