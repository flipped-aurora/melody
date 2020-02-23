package plugin

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

		// load plugin(s)
		r, ok := serverRegister.Get(Namespace)
		if !ok {
			logger.Debug("melody_http_server_handler: no plugins registered for the module")
			return next(ctx, cfg, handler)
		}

		name, nameOk := extra["name"].(string)
		fifoRaw, fifoOk := extra["name"].([]interface{})
		if !nameOk && !fifoOk {
			logger.Debug("melody_http_server_handler: no plugins required in the extra config")
			return next(ctx, cfg, handler)
		}

		fifo := []string{}

		if !fifoOk {
			fifo = []string{name}
		} else {
			for _, x := range fifoRaw {
				if v, ok := x.(string); ok {
					fifo = append(fifo, v)
				}
			}
		}

		for _, name := range fifo {
			rawHf, ok := r.Get(name)
			if !ok {
				logger.Debug("melody_http_server_handler: no plugin resgistered as", name)
				return next(ctx, cfg, handler)
			}

			hf, ok := rawHf.(func(context.Context, map[string]interface{}, http.Handler) (http.Handler, error))
			if !ok {
				logger.Warning("melody_http_server_handler: wrong plugin handler type:", name)
				return next(ctx, cfg, handler)
			}

			handlerWrapper, err := hf(context.Background(), extra, handler)
			if err != nil {
				logger.Warning("melody_http_server_handler: error getting the plugin handler:", err.Error())
				return next(ctx, cfg, handler)
			}

			logger.Debug("melody_http_server_handler: injecting plugin", name)
			handler = handlerWrapper
		}

		return next(ctx, cfg, handler)
	}
}
