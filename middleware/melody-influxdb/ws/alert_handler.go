package ws

import (
	"melody/middleware/melody-alert/model"
	"net/http"
)

func (wsc WebSocketClient) GetWarnings() http.HandlerFunc {
	return wsc.WebSocketHandler(func(request *http.Request, data map[string]interface{}) (i interface{}, err error) {
		return model.WarningList, nil
	})
}
