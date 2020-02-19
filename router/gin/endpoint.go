package gin

import (
	"melody/config"

	"github.com/gin-gonic/gin"
)

type HandleFactory func(*config.EndpointConfig) gin.HandlerFunc
