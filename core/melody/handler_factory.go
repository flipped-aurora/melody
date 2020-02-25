package melody

import (
	"melody/logging"
	juju "melody/middleware/melody-ratelimit/juju/router/gin"
	router "melody/router/gin"
)

// NewHandlerFactory 返回一个Handler工厂
// 根据不同的EndpointConfig定制Handler
// 这里的Handler旨在处理Endpoint层的逻辑
func NewHandlerFactory(logger logging.Logger) router.HandlerFactory {
	handlerFactory := router.EndpointHandler
	handlerFactory = juju.NewRateLimiterMw(handlerFactory)
	return handlerFactory
}
