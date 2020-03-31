package ws

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/influxdata/influxdb/client/v2"
	"math/rand"
	"melody/logging"
	"net/http"
	"time"
)

type WebSocketClient struct {
	Client   client.Client
	Upgrader websocket.Upgrader
	Logger   logging.Logger
	DB       string
	Refresh  chan int
}

func (wsc WebSocketClient) PushTestArray() http.HandlerFunc {

	return wsc.WebSocketHandler(func(request *http.Request) (i interface{}, err error) {
		var array []int
		for i := 0; i < 7; i++ {
			randNum := rand.Int() % 1500
			array = append(array, randNum)
		}
		return map[string]interface{}{
			"xAxis": []string{"1", "2", "3", "4", "5", "6", "7"},
			"yAxis": array,
		}, nil
	})
}

func (wsc WebSocketClient) GetDebugNumGC() http.HandlerFunc {
	return wsc.WebSocketHandler(func(request *http.Request) (interface{}, error) {
		cmd := fmt.Sprintf(`SELECT mean("GCStats.NumGC")
					AS "mean_GCStats.NumGC" FROM "%s"."autogen"."debug" WHERE time > %s - %s AND time <
				%s GROUP BY time(%s) FILL(null)`, wsc.DB, WsTimeControl.MinTime, WsTimeControl.TimeInterval, WsTimeControl.MaxTime, WsTimeControl.GroupTime)
		resu, err := ExecuteQuery(wsc.Client, cmd, wsc.DB)
		if err != nil {
			return nil, err
		}
		result := resu[0]
		if result.Err != "" {
			return nil, errors.New(result.Err)
		}
		values := result.Series[0].Values
		columns := result.Series[0].Columns
		var xAxis []string
		var yAxis []float64
		for _, v := range values {
			// time
			if ts, ok := (v[0]).(json.Number); ok {
				if ti, err := ts.Int64(); err == nil {
					t := time.Unix(ti, 0).Format("15:04:05")
					xAxis = append(xAxis, t)
				} else {
					continue
				}
			}
			// value
			if vs, ok := (v[1]).(json.Number); ok {
				if val, err := vs.Float64(); err == nil {
					yAxis = append(yAxis, val)
				}
			} else {
				yAxis = append(yAxis, 0)
			}
		}

		return map[string]interface{}{
			"xAxis":   xAxis,
			"yAxis":   yAxis,
			"columns": columns,
			"title":   "debug.GCStats.NumGC",
		}, nil
	})
}
