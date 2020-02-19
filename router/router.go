package router

import (
	"context"
	"melody/config"
	http "melody/transport/http/server"
)

// DefaultRunServer default run server func
var DefaultRunServer = http.RunServer

// Router 暴露出去的接口
type Router interface {
	Run(config.ServiceConfig)
}

// Factory 暴露出去的接口
type Factory interface {
	New() Router
	NewWithContext(context.Context) Router
}
