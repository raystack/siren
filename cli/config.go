package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var (
	DefaultHost   = "localhost"
	FileName      = ".siren"
	FileExtension = "yaml"
)

type configuration struct {
	Host string `yaml:"host"`
}

func configCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "manage siren CLI configuration",
	}
	cmd.AddCommand(configInitCommand())
	return cmd
}

func configInitCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "initialize CLI configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			config := configuration{
				Host: DefaultHost,
			}

			b, err := yaml.Marshal(config)
			if err != nil {
				return err
			}

			filepath := fmt.Sprintf("%v.%v", FileName, FileExtension)
			if err := os.WriteFile(filepath, b, 0655); err != nil {
				return err
			}
			fmt.Printf("config created: %v", filepath)

			return nil
		},
	}
}

func readConfig() (*configuration, error) {
	var c configuration
	filepath := fmt.Sprintf("%v.%v", FileName, FileExtension)
	b, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(b, &c); err != nil {
		return nil, err
	}

	return &c, nil
}
