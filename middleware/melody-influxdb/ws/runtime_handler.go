package ws

import (
	"errors"
	"melody/middleware/melody-influxdb/ws/handler"
	"net/http"
)

func (wsc WebSocketClient) GetNumGoroutine() http.HandlerFunc {
	return wsc.WebSocketHandler(func(request *http.Request, data map[string]interface{}) (i interface{}, err error) {
		cmd := wsc.generateCommand(`SELECT mean("NumGoroutine") AS "mean_NumGoroutine" 
		FROM "%s"."autogen"."runtime" WHERE time > %s - %s AND time < %s GROUP BY time(%s) 
		FILL(null)`)
		re, err := wsc.executeQuery(cmd)
		if err != nil {
			return nil, err
		}
		result := re[0]
		if result.Err != "" {
			return nil, errors.New(result.Err)
		}
		if len(result.Series) == 0 {
			return map[string]interface{}{
				"title": "NumGoroutine",
			}, nil
		}
		values := result.Series[0].Values

		var times []string
		var numGoroutine []int64
		handler.ResultDataHandler(&times, values, GetTimeFormat(), &numGoroutine)

		return map[string]interface{}{
			"title": "NumGoroutine",
			"times": times,
			"series": []map[string]interface{}{
				{
					"data":   numGoroutine,
					"name":   "goroutine num",
					"type":   "line",
					"smooth": true,
				},
			},
		}, nil
	})
}

func (wsc WebSocketClient) GetNumGC() http.HandlerFunc {
	return wsc.WebSocketHandler(func(request *http.Request, data map[string]interface{}) (i interface{}, err error) {
		cmd := wsc.generateCommand(`SELECT sum("MemStats.NumGC") AS "sum_MemStats.NumGC" FROM "%s"."autogen"."runtime" WHERE 
		time > %s - %s AND time < %s GROUP BY time(%s) FILL(null)`)

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
				"title": "NumGC",
			}, nil
		}
		values := result.Series[0].Values
		var times []string
		var numGC []int64
		handler.ResultDataHandler(&times, values, GetTimeFormat(), &numGC)

		return map[string]interface{}{
			"title": "NumGC",
			"times": times,
			"series": []map[string]interface{}{
				{
					"data": numGC,
					"name": "GC num",
					"type": "line",
				},
			},
		}, nil
	})
}

func (wsc WebSocketClient) GetNumMemoryFree() http.HandlerFunc {
	return wsc.WebSocketHandler(func(request *http.Request, data map[string]interface{}) (i interface{}, err error) {
		cmd := wsc.generateCommand(`SELECT sum("MemStats.Frees") AS "sum_MemStats.Frees" FROM "%s"."autogen"."runtime" 
		WHERE time > %s - %s AND time < %s GROUP BY time(%s) FILL(null)`)
		re, err := wsc.executeQuery(cmd)

		if err != nil {
			return nil, err
		}
		result := re[0]
		if result.Err != "" {
			return nil, errors.New(result.Err)
		}
		if len(result.Series) == 0 {
			return map[string]interface{}{
				"title": "NumGCFrees",
			}, nil
		}
		values := result.Series[0].Values

		var times []string
		var numGCFrees []int64
		handler.ResultDataHandler(&times, values, GetTimeFormat(), &numGCFrees)

		return map[string]interface{}{
			"title": "NumGCFrees",
			"times": times,
			"series": []map[string]interface{}{
				{
					"data":      numGCFrees,
					"name":      "free memory",
					"type":      "line",
					"areaStyle": map[string]interface{}{},
				},
			},
		}, nil
	})
}

func (wsc WebSocketClient) GetSysMemory() http.HandlerFunc {
	return wsc.WebSocketHandler(func(request *http.Request, data map[string]interface{}) (i interface{}, err error) {
		cmd := wsc.generateCommand(`SELECT mean("MemStats.HeapSys") AS "sum_MemStats.HeapSys", mean("MemStats.MCacheSys")
		AS "sum_MemStats.MCacheSys", mean("MemStats.MSpanSys") AS "sum_MemStats.MSpanSys", mean("MemStats.Sys") AS "sum_MemStats.Sys",
		mean("MemStats.StackSys") AS "sum_MemStats.StackSys" FROM "%s"."autogen"."runtime" WHERE time > %s - %s AND time <
		%s GROUP BY time(%s) FILL(null)`)
		resp, err := wsc.executeQuery(cmd)

		if err != nil {
			return nil, err
		}
		result := resp[0]
		if result.Err != "" {
			return nil, errors.New(result.Err)
		}
		if len(result.Series) == 0 {
			return map[string]interface{}{
				"title": "Memory",
			}, nil
		}
		values := result.Series[0].Values

		var times []string
		var sys, heap, stack, mspan, mcache []int64

		handler.ResultDataHandler(&times, values, GetTimeFormat(), &heap, &mcache, &mspan, &sys, &stack)

		return map[string]interface{}{
			"title": "Memory",
			"times": times,
			"series": []map[string]interface{}{
				{
					"data":      sys,
					"name":      "sys",
					"type":      "line",
					"areaStyle": map[string]interface{}{},
				},
				{
					"data":      heap,
					"name":      "heap",
					"type":      "line",
					"areaStyle": map[string]interface{}{},
				},
				{
					"data":      stack,
					"name":      "stack",
					"type":      "line",
					"areaStyle": map[string]interface{}{},
				},
				{
					"data":      mcache,
					"name":      "mcache",
					"type":      "line",
					"areaStyle": map[string]interface{}{},
				},
				{
					"data":      mspan,
					"name":      "mspan",
					"type":      "line",
					"areaStyle": map[string]interface{}{},
				},
			},
		}, nil
	})
}
