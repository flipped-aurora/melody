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
sum("total") AS "sum_total", sum("current") AS "sum_current"
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

		handler.ResultDataHandler(&times, values, GetTimeFormat(), &total, &current)

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
			},
		}, nil
	})
}

func (wsc WebSocketClient) GetRouterTime() http.HandlerFunc {
	return wsc.WebSocketHandler(func(request *http.Request, data map[string]interface{}) (i interface{}, err error) {
		message, ok := data["message"]
		if !ok {
			return map[string]interface{}{
				"title": "Router Time",
			}, nil
		}
		cmd := wsc.generateCommand(`
SELECT 
sum("max") AS "sum_max", sum("mean") AS "sum_mean", sum("min") AS "sum_min"
FROM 
"%s"."autogen"."router.response-time" 
WHERE 
time > %s - %s AND time < %s AND "name"='` + message.(string) + `'
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
		if len(result.Series) == 0 {
			return map[string]interface{}{
				"title": "No Data",
			}, nil
		}
		var times []string
		var max []int64
		var mean []int64
		var min []int64

		handler.ResultDataHandler(&times, values, GetTimeFormat(), &max, &mean, &min)

		return map[string]interface{}{
			"title": "Router Time " + message.(string),
			"times": times,
			"series": []map[string]interface{}{
				{
					"data":      max,
					"name":      "max",
					"type":      "line",
					"areaStyle": map[string]interface{}{},
				},
				{
					"data":      mean,
					"name":      "mean",
					"type":      "line",
					"areaStyle": map[string]interface{}{},
				},
				{
					"data":      min,
					"name":      "min",
					"type":      "line",
					"areaStyle": map[string]interface{}{},
				},
			},
		}, nil
	})
}

func (wsc WebSocketClient) GetRouterSize() http.HandlerFunc {
	return wsc.WebSocketHandler(func(request *http.Request, data map[string]interface{}) (i interface{}, err error) {
		message, ok := data["message"]
		if !ok {
			return map[string]interface{}{
				"title": "Router Size",
			}, nil
		}
		cmd := wsc.generateCommand(`
SELECT 
sum("max") AS "sum_max", sum("mean") AS "sum_mean", sum("min") AS "sum_min"
FROM 
"%s"."autogen"."router.response-size" 
WHERE 
time > %s - %s AND time < %s AND "name"='` + message.(string) + `'
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
		if len(result.Series) == 0 {
			return map[string]interface{}{
				"title": "No Data",
			}, nil
		}
		values := result.Series[0].Values
		var times []string
		var max []int64
		var mean []int64
		var min []int64

		handler.ResultDataHandler(&times, values, GetTimeFormat(), &max, &mean, &min)

		return map[string]interface{}{
			"title": "Router Size " + message.(string),
			"times": times,
			"series": []map[string]interface{}{
				{
					"data":      max,
					"name":      "max",
					"type":      "line",
					"areaStyle": map[string]interface{}{},
				},
				{
					"data":      mean,
					"name":      "mean",
					"type":      "line",
					"areaStyle": map[string]interface{}{},
				},
				{
					"data":      min,
					"name":      "min",
					"type":      "line",
					"areaStyle": map[string]interface{}{},
				},
			},
		}, nil
	})
}
