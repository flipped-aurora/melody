package gin

import (
	"errors"
	"melody/config"
	httpsecure "melody/middleware/melody-httpsecure"

	"github.com/gin-gonic/gin"
	"github.com/unrolled/secure"
)

var (
	errorNoSecureConfig = errors.New("no extra config for http secure middleware")
)

// Register http security middleware in engine
func Register(c config.ExtraConfig, engine *gin.Engine) error {
	v, ok := httpsecure.GetConfig(c).(secure.Options)
	if !ok {
		return errorNoSecureConfig
	}

	engine.Use(secureMw(v))

	return nil
}

func secureMw(o secure.Options) gin.HandlerFunc {
	secureMiddle := secure.New(o)
	return func(c *gin.Context) {
		err := secureMiddle.Process(c.Writer, c.Request)
		if err != nil {
			c.Abort()
			return
		}

		if status := c.Writer.Status(); status > 300 && status < 399 {
			c.Abort()
		}
	}
}
