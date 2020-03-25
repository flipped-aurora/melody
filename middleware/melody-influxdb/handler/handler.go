package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/influxdata/influxdb/client/v2"
	"melody/logging"
	"net/http"
)

type query struct {
	Command   string `json:"command"`
	Database  string `json:"database"`
	Precision string `json:"precision"`
}

func Query(cli client.Client, logger logging.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		var q query
		if err := c.ShouldBindJSON(&q); err != nil {
			if err != nil {
				logger.Error("parse request body to query object error:", err)
				c.JSON(http.StatusBadRequest, gin.H{"error": err})
				return
			}
		}
		logger.Debug("-> query influxdb with query:", q)
		resp, err := cli.Query(client.NewQuery(q.Command, q.Database, q.Precision))
		if err != nil || resp.Err != "" {
			logger.Error("query influxdb error:", err, ", Err:", resp.Err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err, "Err": resp.Err})
			return
		}
		logger.Debug("<- query success")
		c.JSON(http.StatusOK, resp.Results)
		return
	}
}
