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
		start := pageIndex * 10
		end := (pageIndex + 1) * 10

		if end > len(model.WarningList.Warnings) {
			end = len(model.WarningList.Warnings)
		}
		return map[string]interface{}{
			"warnings": model.WarningList.Warnings[start:end],
			"total":    len(model.WarningList.Warnings),
		}, nil
	})
}
