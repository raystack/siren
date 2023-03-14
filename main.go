package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/goto/siren/cli"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	if err := cli.New().ExecuteContext(ctx); err != nil {
		fmt.Printf("%+v\n", err)
	}
}
