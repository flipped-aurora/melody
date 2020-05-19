package ws

import (
	"github.com/gorilla/websocket"
	"github.com/influxdata/influxdb/client/v2"
	"melody/config"
	"melody/logging"
	"net/http"
)

type WebSocketClient struct {
	Client   client.Client
	Upgrader websocket.Upgrader
	Logger   logging.Logger
	DB       string
	Cfg      *config.ServiceConfig
}

func (wsc WebSocketClient) RegisterHandleFunc() {
	http.HandleFunc("/debug/num/gc", wsc.GetDebugNumGC())
	http.HandleFunc("/debug/alloc", wsc.GetDebugAlloc())

	http.HandleFunc("/runtime/num/gc", wsc.GetNumGC())
	http.HandleFunc("/runtime/num/goroutine", wsc.GetNumGoroutine())
	http.HandleFunc("/runtime/num/frees", wsc.GetNumMemoryFree())
	http.HandleFunc("/runtime/num/memory", wsc.GetSysMemory())

	http.HandleFunc("/requests/complete", wsc.GetRequestsComplete())
	http.HandleFunc("/requests/error", wsc.GetRequestsError())
	http.HandleFunc("/requests/endpoints", wsc.GetRequestsEndpoints())
	http.HandleFunc("/requests/backends", wsc.GetRequestsBackends())
	http.HandleFunc("/requests/api", wsc.GetRequestsAPI())
	http.HandleFunc("/requests/endpoints/pie", wsc.GetRequestsEndpointsPie())
	http.HandleFunc("/requests/backends/pie", wsc.GetRequestsBackendsPie())

	http.HandleFunc("/router/direction", wsc.GetRouterDirection())
	http.HandleFunc("/router/size", wsc.GetRouterSize())
	http.HandleFunc("/router/time", wsc.GetRouterTime())

	http.HandleFunc("/warnings", wsc.GetWarnings())
	http.HandleFunc("/warnings/watch", wsc.WarningsWatch())
}
