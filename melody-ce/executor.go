package melody

import (
	"context"
	"melody/cmd"
	"melody/config"
)

func NewExecutor(ctx context.Context) cmd.Executor {
	return func(config config.ServiceConfig) {
		//TODO start melody server with config file.
	}
}
