package ws

import (
	"errors"
	"melody/middleware/melody-influxdb/ws/convert"
	"net/http"
)

func (wsc WebSocketClient) PushTestArray() http.HandlerFunc {
	return wsc.WebSocketHandler(func(request *http.Request, data map[string]interface{}) (i interface{}, err error) {
		if data != nil {
			if v, ok := data["message"].(string); ok {
				return v, nil
			}
		}
		return
	})
}

func (wsc WebSocketClient) GetDebugNumGC() http.HandlerFunc {
	return wsc.WebSocketHandler(func(request *http.Request, data map[string]interface{}) (interface{}, error) {
		cmd := wsc.generateCommand(`SELECT sum("GCStats.NumGC")
					AS "GCStats.NumGC" FROM "%s"."autogen"."debug" WHERE time > %s - %s AND time <
				%s GROUP BY time(%s) FILL(null)`)
		resu, err := wsc.executeQuery(cmd)
		if err != nil {
			return nil, err
		}
		result := resu[0]
		if result.Err != "" {
			return nil, errors.New(result.Err)
		}
		values := result.Series[0].Values
		var times []string
		var debugNumGC []float64
		for _, v := range values {
			// time
			if t, ok := convert.ObjectToStringTime(v[0], GetTimeFormat()); ok {
				times = append(times, t)
			} else {
				continue
			}
			// value
			if f, ok := convert.ObjectToFloat(v[1]); ok {
				debugNumGC = append(debugNumGC, f)
			} else {
				debugNumGC = append(debugNumGC, 0)
			}
		}

		return map[string]interface{}{
			"title": "debug.NumGC",
			"times": times,
			"series": []map[string]interface{}{
				{
					"data":   debugNumGC,
					"name":   "GCNum",
					"type":   "line",
					"smooth": true,
				},
			},
		}, nil
	})
}
