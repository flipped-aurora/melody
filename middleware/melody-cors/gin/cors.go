package gin

import (
	"melody/config"

	gincors "melody/middleware/melody-cors"

	"github.com/gin-gonic/gin"
	"github.com/rs/cors"
	wrapper "github.com/rs/cors/wrapper/gin"
)

// New 创建cors middleware
func New(extra config.ExtraConfig) gin.HandlerFunc {
	c := gincors.GetConfig(extra)
	if c == nil {
		return nil
	}

	t, ok := c.(gincors.Config)
	if !ok {
		return nil
	}

	return wrapper.New(
		cors.Options{
			AllowedOrigins:   t.AllowOrigins,
			AllowedMethods:   t.AllowMethods,
			AllowedHeaders:   t.AllowHeaders,
			ExposedHeaders:   t.ExposeHeaders,
			AllowCredentials: t.AllowCredentials,
			MaxAge:           int(t.MaxAge.Seconds()),
		})
}
