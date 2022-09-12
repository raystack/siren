package cli

import (
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/odpf/salt/cmdx"
	"github.com/odpf/salt/printer"
	"github.com/spf13/cobra"
)

func configCmd(cmdxConfig *cmdx.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config <command>",
		Short: "Manage siren CLI configuration",
		Example: heredoc.Doc(`
			$ siren config init
			$ siren config list
		`),
	}
	cmd.AddCommand(configInitCommand(cmdxConfig))
	cmd.AddCommand(configListCommand(cmdxConfig))
	return cmd
}

func configInitCommand(cmdxConfig *cmdx.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Initialize CLI configuration",
		Example: heredoc.Doc(`
			$ siren config init
		`),
		Annotations: map[string]string{
			"group": "core",
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdxConfig.Init(&ClientConfig{}); err != nil {
				return err
			}

			fmt.Printf("config created: %v\n", cmdxConfig.File())
			printer.SuccessIcon()
			return nil
		},
	}
}

func configListCommand(cmdxConfig *cmdx.Config) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "list",
		Short: "List client configuration settings",
		Example: heredoc.Doc(`
			$ siren config list
		`),
		Annotations: map[string]string{
			"group": "core",
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			data, err := cmdxConfig.Read()
			if err != nil {
				return ErrClientConfigNotFound
			}

			fmt.Println(data)
			return nil
		},
	}
	return cmd
}
