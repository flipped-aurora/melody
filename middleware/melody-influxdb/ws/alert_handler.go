package ws

import (
	"melody/middleware/melody-alert/model"
	"net/http"
	"strconv"
)

func (wsc WebSocketClient) GetWarnings() http.HandlerFunc {
	return wsc.WebSocketHandler(func(request *http.Request, data map[string]interface{}) (i interface{}, err error) {
		pageIndex := 0
		if tmp, ok := data["message"]; ok {
			if tmp2, ok := tmp.(string); ok {
				atoi, _ := strconv.Atoi(tmp2)
				pageIndex = atoi - 1
			}
		}

		total := len(model.WarningList.Warnings)

		start := total - (pageIndex+1)*10
		end := total - pageIndex*10

		if start < 0 {
			start = 0
		}

		result := make([]model.Warning, 0)

		for i := end - 1; i >= start; i-- {
			result = append(result, model.WarningList.Warnings[i])
		}

		return map[string]interface{}{
			"warnings": result,
			"total":    len(result),
		}, nil
	})
}

func (wsc WebSocketClient) WarningsWatch() http.HandlerFunc {
	return wsc.WebSocketWatchHandler(func(request *http.Request, data map[string]interface{}) (i interface{}, err error) {
		warning := data["warning"]
		return map[string]interface{}{
			"warning": warning,
		}, nil
	})
}
