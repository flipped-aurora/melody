package gin

import (
	"errors"
	"net/http"

	botmonitor "melody/middleware/melody-botmonitor"
	melody "melody/middleware/melody-botmonitor/melody"
	"melody/config"
	"melody/logging"
	"melody/proxy"
	melodygin "melody/router/gin"
	"github.com/gin-gonic/gin"
)

// Register 检查配置，按照需要在gin引擎上注册一个bot检测中间件
func Register(cfg config.ServiceConfig, l logging.Logger, engine *gin.Engine) {
	detectorCfg, err := melody.ParseConfig(cfg.ExtraConfig)
	if err == melody.ErrNoConfig {
		l.Debug("botmonitor middleware: ", err.Error())
		return
	}
	if err != nil {
		l.Warning("botmonitor middleware: ", err.Error())
		return
	}
	d, err := botmonitor.New(detectorCfg)
	if err != nil {
		l.Warning("botmonitor middleware: unable to createt the LRU detector:", err.Error())
		return
	}
	engine.Use(middleware(d))
}

// New 检查配置
func New(hf melodygin.HandlerFactory, l logging.Logger) melodygin.HandlerFactory {
	return func(cfg *config.EndpointConfig, p proxy.Proxy) gin.HandlerFunc {
		next := hf(cfg, p)

		detectorCfg, err := melody.ParseConfig(cfg.ExtraConfig)
		if err == melody.ErrNoConfig {
			l.Debug("botmonitor: ", err.Error())
			return next
		}
		if err != nil {
			l.Warning("botmonitor: ", err.Error())
			return next
		}

		d, err := botmonitor.New(detectorCfg)
		if err != nil {
			l.Warning("botmonitor: unable to create the LRU detector:", err.Error())
			return next
		}
		return handler(d, next)
	}
}

func middleware(f botmonitor.DetectorFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		if f(c.Request) {
			c.AbortWithError(http.StatusForbidden, errBotRejected)
			return
		}

		c.Next()
	}
}

func handler(f botmonitor.DetectorFunc, next gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		if f(c.Request) {
			c.AbortWithError(http.StatusForbidden, errBotRejected)
			return
		}

		next(c)
	}
}

var errBotRejected = errors.New("bot rejected")
