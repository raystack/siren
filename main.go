package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/goto/siren/cli"
)

const (
	exitError = 1
	exitOK    = 0
)

func main() {
	os.Exit(execute())
}

func execute() int {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	if err := cli.New().ExecuteContext(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "\n%v\n", err)
		return exitError
	}

	return exitOK
}
