package ws

import (
	"encoding/json"
	"net/http"
	"os"
	"runtime"
	"runtime/pprof"
	"time"
)

type WebSocketHandlerFunc func(request *http.Request) (interface{}, error)

func (wsc WebSocketClient) WebSocketHandler(handler WebSocketHandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		ws, err := wsc.Upgrader.Upgrade(writer, request, nil)
		if err != nil {
			wsc.Logger.Error("websocket upgrade:", err)
		}
		defer ws.Close()
		go func() {
			for {
				mt, message, err := ws.ReadMessage()
				if err != nil {
					wsc.Logger.Error("read:", err)
					return
				}
				wsc.Logger.Debug("receive:", message, " type:", mt)
			}
		}()
		f, err := os.OpenFile("pprof.log", os.O_CREATE | os.O_APPEND | os.O_RDWR, os.ModePerm)
		if err != nil {
			wsc.Logger.Error(err)
		}
		for {
			res, err := handler(request)
			if err != nil {
				wsc.Logger.Error("websocket handler error:", err)
				errBytes, _ := json.Marshal(map[string]interface{}{"error": err})
				ws.WriteMessage(1, errBytes)
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
			wsc.Logger.Warning("goroutine nums:", runtime.NumGoroutine())
			p := pprof.Lookup("goroutine")
			p.WriteTo(f, 1)
			t := time.NewTicker(WsTimeControl.RefreshTime)
			select {
			case <-t.C:
			case <-wsc.Refresh:
			}
		}
		wsc.Logger.Debug("connect close and handler func end.")
	}
}
