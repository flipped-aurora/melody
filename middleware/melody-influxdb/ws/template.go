package ws

import (
	"fmt"
	"melody/config"
	"strings"
)

func (wsc WebSocketClient) generateCommand(cmd string) string {
	return fmt.Sprintf(cmd, wsc.DB, WsTimeControl.MinTime, WsTimeControl.TimeInterval, WsTimeControl.MaxTime, WsTimeControl.GroupTime)
}

func (wsc WebSocketClient) generateCommandWithEndpoints(cmd string, cfg *config.ServiceConfig) string {
	var endpointStr []string
	for i := range cfg.Endpoints {
		endpointStr = append(endpointStr, cfg.Endpoints[i].Endpoint)
	}
	builder := strings.Builder{}
	builder.WriteString("(")
	for i := 0; i < len(endpointStr); i++ {
		builder.WriteString(`"name"='`)
		builder.WriteString(endpointStr[i])
		if i == len(endpointStr)-1 {
			builder.WriteString("'")
		} else {
			builder.WriteString("' OR ")
		}
	}
	builder.WriteString(")")
	return fmt.Sprintf(
		cmd,
		wsc.DB,
		WsTimeControl.MinTime,
		WsTimeControl.TimeInterval,
		WsTimeControl.MaxTime,
		builder.String(),
		WsTimeControl.GroupTime,
	)
}
