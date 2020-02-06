package main

import (
	"context"
	"log"
	"melody/cmd"
	"melody/config"
	"melody/core/melody"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		select {
		case sig := <-sigs:
			log.Println("Signal intercepted:", sig)
			cancel()
		case <-ctx.Done():
		}
	}()

	cfg := config.New()
	cmd.Execute(cfg, melody.NewExecutor(ctx))
}
