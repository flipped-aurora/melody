package main

import (
	"context"
	"log"
	"melody/cmd"
	"melody/core/melody"
	viper "melody/middleware/melody-viper"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	melody.RegisterEncoders()

	go func() {
		select {
		case sig := <-sigs:
			log.Println("Signal intercepted:", sig)
			cancel()
		case <-ctx.Done():
		}
	}()

	parser := viper.New()
	cmd.Execute(parser, melody.NewExecutor(ctx))
}
