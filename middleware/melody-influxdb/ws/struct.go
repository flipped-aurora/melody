package ws

import (
	"github.com/gorilla/websocket"
	"github.com/influxdata/influxdb/client/v2"
	"melody/logging"
	"net/http"
)

type WebSocketClient struct {
	Client   client.Client
	Upgrader websocket.Upgrader
	Logger   logging.Logger
	DB       string
	Refresh  chan int
}

func (wsc WebSocketClient) RegisterHandleFunc() {
	http.HandleFunc("/debug/num/gc", wsc.GetDebugNumGC())
	http.HandleFunc("/runtime/num/goroutine_thread", wsc.GetGoroutineAndThreadNum())
}
