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
}

func (wsc WebSocketClient) RegisterHandleFunc(cfg *config.ServiceConfig) {
	http.HandleFunc("/debug/num/gc", wsc.GetDebugNumGC())
	http.HandleFunc("/debug/free-total", wsc.GetDebugFreeTotal())

	http.HandleFunc("/runtime/num/gc", wsc.GetNumGC())
	http.HandleFunc("/runtime/num/goroutine", wsc.GetNumGoroutine())
	http.HandleFunc("/runtime/num/frees", wsc.GetNumMemoryFree())
	http.HandleFunc("/runtime/num/memory", wsc.GetSysMemory())

	http.HandleFunc("/requests/complete", wsc.GetRequestsComplete())
	http.HandleFunc("/requests/error", wsc.GetRequestsError())
	http.HandleFunc("/requests/endpoints", wsc.GetRequestsEndpoints())
	http.HandleFunc("/requests/backends", wsc.GetRequestsBackends())
	http.HandleFunc("/requests/api", wsc.GetRequestsAPI())
	http.HandleFunc("/requests/endpoints/pie", wsc.GetRequestsEndpointsPie(cfg))
	http.HandleFunc("/requests/backends/pie", wsc.GetRequestsBackendsPie(cfg))

	http.HandleFunc("/router/direction", wsc.GetRouterDirection())
}
