package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"

	"github.com/cepmap/otus-system-monitoring/internal/config"
	"github.com/cepmap/otus-system-monitoring/internal/logger"
	"github.com/cepmap/otus-system-monitoring/internal/network/server"
)

func main() {
	if err := config.InitConfig(); err != nil {
		logger.Fatal(err.Error())
	}

	logger.Info(fmt.Sprintf("Current config: %+v", config.DaemonConfig))

	ctx, stop := signal.NotifyContext(context.Background(),
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	defer stop()

	srv := server.NewStatsDaemonServer(ctx)

	go func() {
		if err := srv.Start(); err != nil {
			logger.Error(fmt.Sprintf("Server error: %v", err))
		}
	}()

	<-ctx.Done()
	logger.Info("Received shutdown signal")
}
