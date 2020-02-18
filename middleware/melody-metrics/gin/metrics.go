package gin

import (
	"context"
	"melody/config"
	"melody/logging"
	metrics "melody/middleware/melody-metrics"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rcrowley/go-metrics/exp"
)

type Metrics struct {
	*metrics.Metrics
}

func New(c context.Context, e config.ExtraConfig, logger logging.Logger) *Metrics {
	metricsController := Metrics{metrics.New(c, e, logger)}
	if metricsController.Config != nil && !metricsController.Config.EndpointDisabled {
		metricsController.RunEndpoint(c, metricsController.NewEngine(), logger)
	}
	return &metricsController
}

func (m *Metrics) RunEndpoint(c context.Context, engine *gin.Engine, logger logging.Logger) {
	server := &http.Server{
		Addr:    m.Config.ListenAddr,
		Handler: engine,
	}

	go func() {
		logger.Debug("Metrics server listening in", m.Config.ListenAddr)
		logger.Info(server.ListenAndServe())
	}()

	go func() {
		<-c.Done()
		logger.Info("shutting down the stats handler")
		ctx, cancel := context.WithTimeout(c, time.Second)
		server.Shutdown(ctx)
		cancel()
		os.Exit(1)
	}()
}

func (m *Metrics) NewEngine() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	// 紧急恢复middleware
	engine.Use(gin.Recovery())
	// 启用自动重定向
	// 例: /fo/ -> /fo
	engine.RedirectTrailingSlash = true
	// 启用过滤重定向
	// 例: /../fo -> /fo
	engine.RedirectFixedPath = true
	engine.HandleMethodNotAllowed = true

	engine.GET("/__stats", m.NewExpHandler())

	return engine
}

func (m *Metrics) NewExpHandler() gin.HandlerFunc {
	return gin.WrapH(exp.ExpHandler(*m.Registry))
}
