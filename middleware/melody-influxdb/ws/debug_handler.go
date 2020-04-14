package ws

import (
	"errors"
	"melody/middleware/melody-influxdb/ws/convert"
	"melody/middleware/melody-influxdb/ws/handler"
	"net/http"
)

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

// SELECT sum("GCStats.PauseTotal") AS "sum_GCStats.PauseTotal" FROM "melody_data_p1"."autogen"."debug" WHERE time > :dashboardTime: AND time < :upperDashboardTime: GROUP BY time(:interval:) FILL(null)

func (wsc WebSocketClient) GetDebugFreeTotal() http.HandlerFunc {
	return wsc.WebSocketHandler(func(request *http.Request, data map[string]interface{}) (interface{}, error) {
		cmd := wsc.generateCommand(`
SELECT
sum("GCStats.PauseTotal") AS "sum_GCStats.PauseTotal"
FROM
"%s"."autogen"."debug"
WHERE
time > %s - %s AND time < %s 
GROUP BY
time(%s) FILL(null)
`)
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
		var freeTotal []int64
		handler.ResultDataHandler(&times, values, GetTimeFormat(), &freeTotal)
		return map[string]interface{}{
			"title": "debug.FreeTotal",
			"times": times,
			"series": []map[string]interface{}{
				{
					"data":      freeTotal,
					"name":      "FreeTotal",
					"type":      "line",
					"areaStyle": map[string]interface{}{},
				},
			},
		}, nil
	})
}
