package ws

import (
	"errors"
	"fmt"
	"melody/middleware/melody-influxdb/ws/handler"
	"net/http"
	"strings"
)

func (wsc WebSocketClient) GetRequestsComplete() http.HandlerFunc {
	return wsc.WebSocketHandler(func(request *http.Request, data map[string]interface{}) (i interface{}, err error) {
		cmd := wsc.generateCommand(`
SELECT 
sum("total") AS "sum_total", sum("count") AS "sum_count" 
FROM 
"%s"."autogen"."requests" 
WHERE 
time > %s - %s AND time < %s AND "complete"='true'
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
		var count []int64
		handler.ResultDataHandler(&times, values, GetTimeFormat(), &total, &count)

		return map[string]interface{}{
			"title": "Requests Complete",
			"times": times,
			"series": []map[string]interface{}{
				{
					"data": total,
					"name": "total",
					"type": "line",
				},
				{
					"data": count,
					"name": "count",
					"type": "line",
				},
			},
		}, nil
	})
}

func (wsc WebSocketClient) GetRequestsError() http.HandlerFunc {
	return wsc.WebSocketHandler(func(request *http.Request, data map[string]interface{}) (i interface{}, err error) {
		cmd := wsc.generateCommand(`
SELECT 
sum("total") AS "sum_total", sum("count") AS "sum_count" 
FROM 
"%s"."autogen"."requests" 
WHERE 
time > %s - %s AND time < %s AND "error"='true'
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
		var count []int64
		handler.ResultDataHandler(&times, values, GetTimeFormat(), &total, &count)

		return map[string]interface{}{
			"title": "Requests Error",
			"times": times,
			"series": []map[string]interface{}{
				{
					"data": total,
					"name": "total",
					"type": "line",
				},
				{
					"data": count,
					"name": "count",
					"type": "line",
				},
			},
		}, nil
	})
}

func (wsc WebSocketClient) GetRequestsEndpoints() http.HandlerFunc {
	return wsc.WebSocketHandler(func(request *http.Request, data map[string]interface{}) (i interface{}, err error) {
		if _, ok := data["message"]; !ok {
			return map[string]interface{}{
				"title": "Requests Endpoints",
			}, nil
		}
		status := data["message"].(string)
		title := status
		if status == "ALL" {
			status = ""
		} else {
			status = ` AND "` + strings.ToLower(status) + `"='true'`
		}
		cmd := wsc.generateCommand(`
SELECT 
sum("total") AS "sum_total", sum("count") AS "sum_count" 
FROM 
"%s"."autogen"."requests" 
WHERE 
time > %s - %s AND time < %s AND "layer"='endpoint'` + status + `
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
		var count []int64
		handler.ResultDataHandler(&times, values, GetTimeFormat(), &total, &count)

		return map[string]interface{}{
			"title": fmt.Sprintf("Requests Endpoints %s", title),
			"times": times,
			"series": []map[string]interface{}{
				{
					"data": total,
					"name": "total",
					"type": "line",
				},
				{
					"data": count,
					"name": "count",
					"type": "line",
				},
			},
		}, nil
	})
}

func (wsc WebSocketClient) GetRequestsBackends() http.HandlerFunc {
	return wsc.WebSocketHandler(func(request *http.Request, data map[string]interface{}) (i interface{}, err error) {
		if _, ok := data["message"]; !ok {
			return map[string]interface{}{
				"title": "Requests Backends",
			}, nil
		}
		status := data["message"].(string)
		title := status
		if status == "ALL" {
			status = ""
		} else {
			status = ` AND "` + strings.ToLower(status) + `"='true'`
		}
		fmt.Println(status, title)
		cmd := wsc.generateCommand(`
SELECT 
sum("total") AS "sum_total", sum("count") AS "sum_count" 
FROM 
"%s"."autogen"."requests" 
WHERE 
time > %s - %s AND time < %s AND "layer"='backend'` + status + `
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
		var count []int64
		handler.ResultDataHandler(&times, values, GetTimeFormat(), &total, &count)

		return map[string]interface{}{
			"title": fmt.Sprintf("Requests Backends %s", title),
			"times": times,
			"series": []map[string]interface{}{
				{
					"data": total,
					"name": "total",
					"type": "line",
				},
				{
					"data": count,
					"name": "count",
					"type": "line",
				},
			},
		}, nil
	})
}

func (wsc WebSocketClient) GetRequestsAPI() http.HandlerFunc {
	return wsc.WebSocketHandler(func(request *http.Request, data map[string]interface{}) (i interface{}, err error) {
		message, ok := data["message"]
		if !ok {
			return map[string]interface{}{
				"title": "Requests API",
			}, nil
		}

		api := strings.Fields(message.(string))
		cmd := wsc.generateCommand(`
SELECT 
sum("total") AS "sum_total", sum("count") AS "sum_count" 
FROM 
"%s"."autogen"."requests" 
WHERE 
time > %s - %s AND time < %s AND "name"='` + api[0] + `' AND "` + strings.ToLower(api[1]) + `"='true'
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
		var count []int64
		handler.ResultDataHandler(&times, values, GetTimeFormat(), &total, &count)

		return map[string]interface{}{
			"title": fmt.Sprintf("%s %s", api[0], api[1]),
			"times": times,
			"series": []map[string]interface{}{
				{
					"data": total,
					"name": "total",
					"type": "line",
				},
				{
					"data": count,
					"name": "count",
					"type": "line",
				},
			},
		}, nil
	})
}
