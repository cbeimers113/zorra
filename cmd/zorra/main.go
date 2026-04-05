package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/cbeimers113/zorra/internal/channel"
	"github.com/cbeimers113/zorra/internal/cli"
	"github.com/cbeimers113/zorra/internal/log"
)

func main() {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer stop()

	if err := channel.LoadChannels(); err != nil {
		log.Warnf("Unable to load channel map: %s", err.Error())
	}

	cli.Execute(ctx)
}
