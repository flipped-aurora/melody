package ws

import (
	"errors"
	"fmt"
	"melody/middleware/melody-influxdb/ws/convert"
	"net/http"
)

func (wsc WebSocketClient) GetGoroutineAndThreadNum() http.HandlerFunc {
	return wsc.WebSocketHandler(func(request *http.Request) (i interface{}, err error) {
		cmd := fmt.Sprintf(`SELECT sum("NumGoroutine") AS "mean_NumGoroutine", 
		sum("NumThread") AS "mean_NumThread" FROM "%s"."autogen"."runtime" WHERE time
		> %s - %s AND time < %s GROUP BY time(%s) FILL(null)`,wsc.DB, WsTimeControl.MinTime, WsTimeControl.TimeInterval, WsTimeControl.MaxTime, WsTimeControl.GroupTime)
		re, err := ExecuteQuery(wsc.Client, cmd, wsc.DB)
		if err != nil {
			return nil, err
		}
		result := re[0]
		if result.Err != "" {
			return nil, errors.New(result.Err)
		}
		values := result.Series[0].Values

		var times []string
		var numGoroutine, numThread []int64
		for _, v := range values {
			// time
			if t, ok := convert.ObjectToStringTime(v[0], GetTimeFormat()); ok {
				times = append(times, t)
			} else {
				continue
			}
			// goroutine
			if g, ok := convert.ObjectToInt(v[1]); ok {
				numGoroutine = append(numGoroutine, g)
			} else {
				numGoroutine = append(numGoroutine, 0)
			}
			// thread
			if t, ok := convert.ObjectToInt(v[2]); ok {
				numThread = append(numThread, t)
			} else {
				numThread = append(numThread, 0)
			}
		}

		return map[string]interface{}{
			"title": "NumGoroutine&NumThread",
			"times": times,
			"series": []map[string]interface{}{
				{
					"data": numGoroutine,
					"name": "NumGoroutine",
					"type": "line",
					"smooth": true,
				},
				{
					"data": numThread,
					"name": "NumThread",
					"type": "line",
					"smooth": true,
				},
			},
		}, nil
	})
}
