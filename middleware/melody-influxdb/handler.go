package influxdb

import (
	"github.com/gin-gonic/gin"
	"github.com/influxdata/influxdb/client/v2"
	"melody/logging"
	"melody/middleware/melody-influxdb/response"
	"net/http"
)

const (
	requestFailCode = 201
)

type query struct {
	Command   string `json:"command"`
	Precision string `json:"precision"`
}

type AuthConfig struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Query(cli client.Client, logger logging.Logger, config influxdbConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		var q query
		if err := c.ShouldBindJSON(&q); err != nil {
			logger.Error("parse request body to query object error:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}
		logger.Debug("-> query influxdb with query:", q)
		resp, err := cli.Query(client.NewQuery(q.Command, config.db, q.Precision))
		if err != nil || resp.Err != "" {
			errMsg := err.Error()
			if resp != nil && resp.Err != "" {
				errMsg = resp.Err
			}
			logger.Error("query influxdb error:", errMsg)
			c.JSON(http.StatusBadRequest, gin.H{"error": errMsg})
			return
		}
		logger.Debug("<- query success")
		response.Ok(c, http.StatusOK, "", resp.Results[0].Series[0])
		return
	}
}

func Ping(logger logging.Logger, config influxdbConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		var con AuthConfig
		err := c.ShouldBindJSON(&con)
		if  err != nil {
			logger.Debug("parse request body to query object error:", err)
			response.Ok(c, requestFailCode, "parse request body error", nil)
			return
		}

		if con.Username != config.password || con.Password != config.password{
			logger.Debug("influx db username or password incorrect")
			response.Ok(c, requestFailCode, "username or password incorrect", nil)
			return
		}

		response.Ok(c, http.StatusOK, "ping success", config.db)
	}
}
