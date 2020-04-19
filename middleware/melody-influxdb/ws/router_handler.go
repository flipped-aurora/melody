package ws

import (
	"errors"
	"melody/middleware/melody-influxdb/ws/handler"
	"net/http"
)

func (wsc WebSocketClient) GetRouterDirection() http.HandlerFunc {
	return wsc.WebSocketHandler(func(request *http.Request, data map[string]interface{}) (i interface{}, err error) {
		cmd := wsc.generateCommand(`
SELECT 
sum("total") AS "sum_total", sum("current") AS "sum_current", sum("gauge") AS "sum_gauge"
FROM 
"%s"."autogen"."router" 
WHERE 
time > %s - %s AND time < %s AND ("direction"='out' OR "direction"='in')
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
		var total []int64
		var current []int64
		var gauge []int64

		handler.ResultDataHandler(&times, values, GetTimeFormat(), &total, &current, &gauge)

		return map[string]interface{}{
			"title": "Router Direction",
			"times": times,
			"series": []map[string]interface{}{
				{
					"data": total,
					"name": "total",
					"type": "line",
				},
				{
					"data": current,
					"name": "current",
					"type": "line",
				},
				{
					"data": gauge,
					"name": "gauge",
					"type": "line",
				},
			},
		}, nil
	})
}
