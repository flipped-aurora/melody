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

	http.HandleFunc("/runtime/num/gc", wsc.GetNumGC())
	http.HandleFunc("/runtime/num/goroutine", wsc.GetNumGoroutine())
	http.HandleFunc("/runtime/num/frees", wsc.GetNumMemoryFree())
	http.HandleFunc("/runtime/num/memory", wsc.GetSysMemory())

	http.HandleFunc("/test", wsc.PushTestArray())
}
