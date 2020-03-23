package gin

import (
	"context"
	"melody/config"
	"melody/logging"
	metrics "melody/middleware/melody-metrics"
	"melody/proxy"
	melodygin "melody/router/gin"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rcrowley/go-metrics/exp"
)

// Metrics 定义了包装计数器
type Metrics struct {
	*metrics.Metrics
}

type ginResponseWriter struct {
	gin.ResponseWriter
	name string
	begin time.Time
	rm *metrics.RouterMetrics
}

// New 返回一个基础的计数控制器
func New(c context.Context, e config.ExtraConfig, logger logging.Logger) *Metrics {
	metricsController := Metrics{metrics.New(c, e, logger)}
	if metricsController.Config != nil && !metricsController.Config.EndpointDisabled {
		metricsController.RunEndpoint(c, metricsController.NewEngine(), logger)
	}
	return &metricsController
}

// RunEndpoint 驱动计数器server，开始计数
func (m *Metrics) RunEndpoint(c context.Context, engine *gin.Engine, logger logging.Logger) {
	server := &http.Server{
		Addr:    m.Config.ListenAddr,
		Handler: engine,
	}

	go func() {
		logger.Info("Metrics server listening in", m.Config.ListenAddr, "🎁")
		logger.Error(server.ListenAndServe())
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

// NewEngine 返回一个gin.Engine去驱动metrics的运行
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

// NewExpHandler 返回一个json的数据统计结果
func (m *Metrics) NewExpHandler() gin.HandlerFunc {
	return gin.WrapH(exp.ExpHandler(*m.Registry))
}

func (m *Metrics) NewHTTPHandleFactory(handleFactory melodygin.HandlerFactory) melodygin.HandlerFactory {
	if m.Config == nil || m.Config.RouterDisabled {
		return handleFactory
	}
	return NewHTTPHandleFactory(m.Router, handleFactory)
}

func NewHTTPHandleFactory(routerMetrics *metrics.RouterMetrics, handleFactory melodygin.HandlerFactory) melodygin.HandlerFactory {
	return func(endpointConfig *config.EndpointConfig, proxy proxy.Proxy) gin.HandlerFunc {
		next := handleFactory(endpointConfig, proxy)
		routerMetrics.RegisterResponseWriterMetrics(endpointConfig.Endpoint)
		return func(c *gin.Context) {
			rw := &ginResponseWriter{
				ResponseWriter: c.Writer,
				name:           endpointConfig.Endpoint,
				begin:          time.Now(),
				rm:             routerMetrics,
			}
			c.Writer = rw

			next(c)

			rw.end()
			routerMetrics.Disconnection()
		}
	}
}

func (gw *ginResponseWriter) end() {
	duration := time.Since(gw.begin)
	gw.rm.Counter("response", gw.name, "status", strconv.Itoa(gw.Status()), "count").Inc(1)
	gw.rm.Histogram("response", gw.name, "size").Update(int64(gw.Size()))
	gw.rm.Histogram("response", gw.name, "time").Update(int64(duration))
}


