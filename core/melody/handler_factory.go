package melody

import (
	"melody/logging"
	botmonitor "melody/middleware/melody-botmonitor/gin"
	jose "melody/middleware/melody-jose"
	ginjose "melody/middleware/melody-jose/gin"
	metrics "melody/middleware/melody-metrics/gin"
	juju "melody/middleware/melody-ratelimit/juju/router/gin"
	router "melody/router/gin"
)

// NewHandlerFactory 返回一个Handler工厂
// 根据不同的EndpointConfig定制Handler
// 这里的Handler旨在处理Endpoint层的逻辑
func NewHandlerFactory(logger logging.Logger, rejecter jose.RejecterFactory, metrics *metrics.Metrics) router.HandlerFactory {
	handlerFactory := router.EndpointHandler
	handlerFactory = juju.NewRateLimiterMw(handlerFactory)
	handlerFactory = ginjose.HandlerFactory(handlerFactory, logger, rejecter)
	handlerFactory = botmonitor.New(handlerFactory, logger)
	handlerFactory = metrics.NewHTTPHandleFactory(handlerFactory)
	return handlerFactory
}
