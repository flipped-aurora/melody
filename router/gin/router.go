package gin

import (
	"context"
	"github.com/gin-gonic/gin"
	"melody/config"
	"melody/logging"
	"melody/proxy"
	"melody/router"
	"net/http"
	"strings"
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

func DefaultFactory(proxyFactory proxy.Factory, logger logging.Logger) router.Factory {
	return NewFactory(
		Config{
			Engine:         gin.Default(),
			MiddleWares:    []gin.HandlerFunc{},
			HandlerFactory: EndpointHandler,
			ProxyFactory:   proxyFactory,
			Logger:         logger,
			RunServer:      router.DefaultRunServer,
		},
	)
}

func (r ginRouter) Run(config config.ServiceConfig) {
	if !config.Debug {
		gin.SetMode(gin.ReleaseMode)
	} else {
		r.cfg.Logger.Debug("Melody Debug enable")
	}

	router.InitHTTPDefaultTransport(config)

	r.cfg.Engine.RedirectTrailingSlash = true
	r.cfg.Engine.RedirectFixedPath = true
	r.cfg.Engine.HandleMethodNotAllowed = true

	// 启用 Middleware
	r.cfg.Engine.Use(r.cfg.MiddleWares...)

	// 注册 Debug 路由
	if config.Debug {
		r.registerDebugEndpoints()
	}

	// 注册所有Endpoints
	r.registerMelodyEndpoints(config.Endpoints)

	// 处理404请求
	r.cfg.Engine.NoRoute(func(c *gin.Context) {
		c.Header(router.HeaderCompleteKey, router.HeaderInCompleteResponseValue)
	})

	// Run Melody server
	if err := r.RunServer(r.ctx, config, r.cfg.Engine); err != nil {
		r.cfg.Logger.Error(err.Error())
	}

	r.cfg.Logger.Info("Melody server execution ended")
}

func (r ginRouter) registerMelodyEndpoints(endpoints []*config.EndpointConfig) {
	for _, e := range endpoints {
		proxyStack, err := r.cfg.ProxyFactory.New(e)
		if err != nil {
			r.cfg.Logger.Error("calling the ProxyFactory", err.Error())
			continue
		}

		r.registerMelodyEndpoint(e.Method, e.Endpoint, r.cfg.HandlerFactory(e, proxyStack), len(e.Backends))
	}
}

func (r ginRouter) registerMelodyEndpoint(method, path string, handler gin.HandlerFunc, totBackends int) {
	//if requestMethod != http.MethodGet && totBackends > 1 {
	//
	//}
	method = strings.ToTitle(method)
	switch method {
	case http.MethodGet:
		r.cfg.Engine.GET(path, handler)
	case http.MethodPost:
		r.cfg.Engine.POST(path, handler)
	case http.MethodPut:
		r.cfg.Engine.PUT(path, handler)
	case http.MethodPatch:
		r.cfg.Engine.PATCH(path, handler)
	case http.MethodDelete:
		r.cfg.Engine.DELETE(path, handler)
	default:
		r.cfg.Logger.Error("Unsupported method", method)
	}
}

func (r ginRouter) registerDebugEndpoints() {
	debugHandler := DebugHandler(r.cfg.Logger)
	r.cfg.Engine.GET("/__debug/*param", debugHandler)
	r.cfg.Engine.POST("/__debug/*param", debugHandler)
	r.cfg.Engine.PUT("/__debug/*param", debugHandler)
	r.cfg.Engine.DELETE("/__debug/*param", debugHandler)
}
