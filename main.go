package main

import (
	"context"
	"fmt"

	"github.com/odpf/siren/cli"
)

func main() {
	if err := cli.New(context.Background()).Execute(); err != nil {
		fmt.Printf("%+v\n", err)
	}
}
