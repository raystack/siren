package cmd

import (
	"errors"
	"fmt"

	"github.com/odpf/siren/client"
	"github.com/odpf/siren/config"
	"github.com/odpf/siren/pkg/uploader"
	"github.com/spf13/cobra"
)

func uploadCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "upload",
		Short: "Upload Rules or Templates YAML file",
		RunE:  upload,
	}
}

func upload(cmd *cobra.Command, args []string) error {
	c := config.LoadConfig()
	s := uploader.NewService(&c.SirenService)
	result, err := s.Upload(args[0])
	//print all resources(succeed or failed in upsert)
	if err != nil {
		fmt.Println(err)
		return err
	}
	switch obj := result.(type) {
	case *client.Template:
		printTemplate(obj)
	case []*client.Rule:
		printRules(obj)
	default:
		return errors.New("unknown response")
	}
	return nil
}

func printRules(rules []*client.Rule) {
	for i := 0; i < len(rules); i++ {
		fmt.Println("Upserted Rule")
		fmt.Println("ID:", rules[i].Id)
		fmt.Println("Name:", rules[i].Name)
		fmt.Println("Name:", rules[i].Status)
		fmt.Println("CreatedAt At:", rules[i].CreatedAt)
		fmt.Println("UpdatedAt At:", rules[i].UpdatedAt)
		fmt.Println()
	}
}

func printTemplate(template *client.Template) {
	if template == nil {
		return
	}
	fmt.Println("Upserted Template")
	fmt.Println("ID:", template.Id)
	fmt.Println("Name:", template.Name)
	fmt.Println("CreatedAt At:", template.CreatedAt)
	fmt.Println("UpdatedAt At:", template.UpdatedAt)
	fmt.Println("Tags:", template.Tags)
	fmt.Println("Variables:", template.Variables)
}
