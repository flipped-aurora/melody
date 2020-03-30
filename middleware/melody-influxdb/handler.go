package influxdb

import (
	"github.com/gin-gonic/gin"
	"melody/middleware/melody-influxdb/response"
	"melody/middleware/melody-influxdb/ws"
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

func (cw clientWrapper) Query() gin.HandlerFunc {
	return func(c *gin.Context) {
		var q query
		if err := c.ShouldBindJSON(&q); err != nil {
			cw.logger.Error("parse request body to query object error:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}
		cw.logger.Debug("-> query influxdb with query:", q)
		res, err := ws.ExecuteQuery(cw.client, q.Command, cw.config.db)
		if err != nil {
			cw.logger.Error("query influxdb error:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}
		cw.logger.Debug("<- query success")
		response.Ok(c, http.StatusOK, "", res[0].Series[0])
		return
	}
}

func (cw clientWrapper) Ping() gin.HandlerFunc {
	return func(c *gin.Context) {
		var con AuthConfig
		err := c.ShouldBindJSON(&con)
		if err != nil {
			cw.logger.Debug("parse request body to query object error:", err)
			response.Ok(c, requestFailCode, "parse request body error", nil)
			return
		}

		if con.Username != cw.config.password || con.Password != cw.config.password {
			cw.logger.Debug("influx db username or password incorrect")
			response.Ok(c, requestFailCode, "username or password incorrect", nil)
			return
		}

		response.Ok(c, http.StatusOK, "ping success", cw.config.db)
	}
}
