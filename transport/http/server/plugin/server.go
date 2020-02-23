package server

import (
	"context"
	"melody/config"
	"melody/logging"
	"net/http"
)

// Namespace 是否运行使用server handler
const Namespace = "melody_http_server_handler"

// RunServer 定义了运行http.Server的函数结构
type RunServer func(context.Context, config.ServiceConfig, http.Handler) error

// New 返回下一个RunServer
func New(logger logging.Logger, next RunServer) RunServer {
	return func(ctx context.Context, cfg config.ServiceConfig, handler http.Handler) error {
		// 根据配置文件检察是否开启handler
		v, ok := cfg.ExtraConfig[Namespace]
		if !ok {
			logger.Debug("melody_http_server_handler: no extra config")
			return next(ctx, cfg, handler)
		}
		extra, ok := v.(map[string]interface{})
		if !ok {
			logger.Debug("melody_http_server_handler: wrong extra config type")
			return next(ctx, cfg, handler)
		}

		// TODO load plugin(s)
		return next(ctx, cfg, handler)
	}
}
