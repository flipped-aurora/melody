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

func (wsc WebSocketClient) GetRequestsEndpointsPie() http.HandlerFunc {
	var endpoints []string
	for _, endpointCfg := range wsc.Cfg.Endpoints {
		endpoints = append(endpoints, endpointCfg.Endpoint)
	}

	return wsc.WebSocketHandler(func(request *http.Request, data map[string]interface{}) (i interface{}, err error) {
		var totals []interface{}
		for _, path := range endpoints {
			cmd := wsc.generateCommandWithSingle(`
SELECT 
max("total") AS "max_total"
FROM 
"%s"."autogen"."requests" 
WHERE 
time > %s - %s AND time < %s AND "layer"='endpoint' AND "name"='` + path + `'`)

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
			var total int64
			handler.ResultSingleDataHandler(&times, values, GetTimeFormat(), &total)
			totals = append(totals, map[string]interface{}{
				"name":  path,
				"value": total,
			})
		}

		return map[string]interface{}{
			"title": "Endpoints Pie",
			"series": []map[string]interface{}{
				{
					"name":   "Endpoints",
					"type":   "pie",
					"radius": "55%",
					"center": []string{"50%", "50%"},
					"data":   totals,
				},
			},
		}, nil
	})
}

func (wsc WebSocketClient) GetRequestsBackendsPie() http.HandlerFunc {
	var backends []string
	set := make(map[string]bool)
	for _, endpointCfg := range wsc.Cfg.Endpoints {
		for _, backend := range endpointCfg.Backends {
			if _, ok := set[backend.URLPattern]; !ok {
				backends = append(backends, backend.URLPattern)
				set[backend.URLPattern] = true
			}
		}
	}

	return wsc.WebSocketHandler(func(request *http.Request, data map[string]interface{}) (i interface{}, err error) {
		var totals []interface{}
		for _, path := range backends {
			cmd := wsc.generateCommandWithSingle(`
SELECT 
max("total") AS "max_total"
FROM 
"%s"."autogen"."requests" 
WHERE 
time > %s - %s AND time < %s AND "layer"='backend' AND "name"='` + path + `'`)

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
			var total int64
			handler.ResultSingleDataHandler(&times, values, GetTimeFormat(), &total)
			totals = append(totals, map[string]interface{}{
				"name":  path,
				"value": total,
			})
		}

		return map[string]interface{}{
			"title": "Backends Pie",
			"series": []map[string]interface{}{
				{
					"name":   "Backends",
					"type":   "pie",
					"radius": "55%",
					"center": []string{"50%", "50%"},
					"data":   totals,
				},
			},
		}, nil
	})
}
