package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/cbeimers113/zorra/internal/cli"
)

func main() {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer stop()

	cli.Execute(ctx)
}
