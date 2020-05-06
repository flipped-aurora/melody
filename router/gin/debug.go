package gin

import (
	"github.com/gin-gonic/gin"

	"melody/logging"
)

// DebugHandler creates a dummy handler function, useful for quick integration tests
func DebugHandler(logger logging.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	}
}
