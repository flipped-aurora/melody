package gin

import (
	"context"
	"melody/config"
	"melody/logging"
	"melody/proxy"
	"melody/router"
	"net/http"

	"github.com/gin-gonic/gin"
)

// RunServerFunc 定义了Melody Server的运行方法
type RunServerFunc func(context.Context, config.ServiceConfig, http.Handler) error

// Config 整个melody server的结构
type Config struct {
	Engine         *gin.Engine
	MiddleWares    []gin.HandlerFunc
	HandlerFactory HandlerFactory
	ProxyFactory   proxy.Factory
	Logger         logging.Logger
	RunServer      RunServerFunc
}

type ginRouter struct {
	cfg       Config
	ctx       context.Context
	RunServer RunServerFunc
}

type factory struct {
	cfg Config
}

// NewFactory 返回默认router factory
func NewFactory(cfg Config) router.Factory {
	return factory{cfg: cfg}
}

func (f factory) New() router.Router {
	return f.NewWithContext(context.Background())
}

func (f factory) NewWithContext(ctx context.Context) router.Router {
	return ginRouter{
		cfg:       f.cfg,
		ctx:       ctx,
		RunServer: f.cfg.RunServer,
	}
}

func (r ginRouter) Run(config config.ServiceConfig) {

}

func (r ginRouter) registerMelodyEndpoints(endpoints []*config.EndpointConfig) {

}

func (r ginRouter) registerMelodyEndpoint(method, path string, handler gin.HandlerFunc, totBackends int) {

}
