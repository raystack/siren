package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/odpf/siren/cli"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	if err := cli.New().ExecuteContext(ctx); err != nil {
		fmt.Printf("%+v\n", err)
	}
}
