package ws

import (
	"encoding/json"
	influxrefresh "melody/middleware/melody-influxdb/refresh"
	"net/http"
	"time"
)

type WebSocketHandlerFunc func(request *http.Request, data map[string]interface{}) (interface{}, error)

func (wsc WebSocketClient) WebSocketHandler(handler WebSocketHandlerFunc) http.HandlerFunc {

	return func(writer http.ResponseWriter, request *http.Request) {
		ch := make(chan int)
		refresh := &influxrefresh.Refresh{Value: &ch}
		influxrefresh.RefreshList.Add(refresh)

		ws, err := wsc.Upgrader.Upgrade(writer, request, nil)
		if err != nil {
			wsc.Logger.Error("websocket upgrade:", err)
		}
		data := make(map[string]interface{})
		defer func() {
			influxrefresh.RefreshList.Remove(refresh)
			close(ch)
			ws.Close()
		}()
		go func(data map[string]interface{}) {
			for {
				mt, message, err := ws.ReadMessage()
				if err != nil {
					wsc.Logger.Error("read:", err)
					return
				}
				wsc.Logger.Debug("receive:", string(message), " type:", mt)
				data["message"] = string(message)
				ch <- 1
			}
		}(data)
		for {
			res, err := handler(request, data)
			if err != nil {
				wsc.Logger.Error("websocket handler error:", err)
				errBytes, _ := json.Marshal(map[string]interface{}{"error": err})
				_ = ws.WriteMessage(1, errBytes)
				break
			}
			bytes, err := json.Marshal(res)
			if err != nil {
				wsc.Logger.Error("marshal json:", err)
				continue
			}
			err = ws.WriteMessage(1, bytes)
			if err != nil {
				wsc.Logger.Debug("write:", err)
				break
			}
			wsc.Logger.Debug("send:", len(string(bytes)), "byte data.")
			t := time.NewTicker(WsTimeControl.RefreshTime)
			select {
			case <-t.C:
			case <-ch:
			}
		}
		wsc.Logger.Debug("connect close and handler func end.")
	}
}
