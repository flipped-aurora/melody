package ws

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"melody/logging"
	"melody/middleware/melody-influxdb/ws/handler"
	"net/http"
	"time"
)

func Test(upgrader websocket.Upgrader, logger logging.Logger) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		ws, err := upgrader.Upgrade(writer, request, nil)
		if err != nil {
			logger.Error("websocket upgrade:", err)
		}
		defer ws.Close()
		done := make(chan int)
		go func(done chan<- int) {
			for {
				mt, message, err := ws.ReadMessage()
				if err != nil {
					logger.Error("read:", err)
					return
				}
				logger.Debug("receive:", message, " type:", mt)
			}
		}(done)

		for {
			bytes, err := json.Marshal(handler.PushTestArray())
			if err != nil {
				logger.Error("marshal json:", err)
				continue
			}
			err = ws.WriteMessage(1, bytes)
			if err != nil {
				log.Println("write:", err)
				break
			}
			logger.Debug("send:", string(bytes))
			time.Sleep(time.Second)
		}
		logger.Debug("connect close and handler func end.")
	}
}
