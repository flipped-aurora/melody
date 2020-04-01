package ws

import (
	"errors"
	"fmt"
	"math/rand"
	"melody/middleware/melody-influxdb/ws/convert"
	"net/http"
)

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
		cmd := fmt.Sprintf(`SELECT sum("GCStats.NumGC")
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
			t, ok := convert.ObjectToStringTime(v[0], GetTimeFormat())
			if !ok {
				continue
			}
			xAxis = append(xAxis, t)
			// value
			if f, ok := convert.ObjectToFloat(v[1]); ok {
				yAxis = append(yAxis, f)
			} else {
				yAxis = append(yAxis, 0)
			}
		}

		return map[string]interface{}{
			"xAxis":   xAxis,
			"yAxis":   yAxis,
			"columns": columns,
			"title":   "GCStats.NumGC",
		}, nil
	})
}
