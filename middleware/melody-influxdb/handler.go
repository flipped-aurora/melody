package influxdb

import (
	"github.com/gin-gonic/gin"
	"melody/middleware/melody-influxdb/response"
	"melody/middleware/melody-influxdb/ws"
	"net/http"
	"time"
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

func (cw *clientWrapper) Query() gin.HandlerFunc {
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

func (cw *clientWrapper) Ping() gin.HandlerFunc {
	return func(c *gin.Context) {
		var con AuthConfig
		err := c.ShouldBindJSON(&con)
		if err != nil {
			cw.logger.Error("parse request body to query object error:", err)
			response.Ok(c, requestFailCode, "parse request body error", nil)
			return
		}

		if con.Username != cw.config.password || con.Password != cw.config.password {
			cw.logger.Error("influx db username or password incorrect")
			response.Ok(c, requestFailCode, "username or password incorrect", nil)
			return
		}

		response.Ok(c, http.StatusOK, "ping success", cw.config.db)
	}
}

func (cw *clientWrapper) ModifyTimeControl() gin.HandlerFunc {
	return func(context *gin.Context) {
		var t ws.TimeControl
		err := context.ShouldBindJSON(&t)
		if err != nil {
			cw.logger.Error("parse request body to time control error:", err)
			response.Ok(context, requestFailCode, "parse request body error", nil)
			return
		}

		d, err := time.ParseDuration(t.RefreshParam)
		if err != nil {
			cw.logger.Error("refresh time can not convert to time.Duration :", err)
			response.Ok(context, requestFailCode, "refresh time can not convert to time.Duration :", nil)
			return
		}
		t.RefreshTime = d
		ws.SetTimeControl(t)
		cw.Refresh <- 1
		response.Ok(context, http.StatusOK, "modify success", nil)
	}
}
