package ws

import "fmt"

func (wsc WebSocketClient) generateCommand(cmd string) string {
	return fmt.Sprintf(cmd, wsc.DB, WsTimeControl.MinTime, WsTimeControl.TimeInterval, WsTimeControl.MaxTime, WsTimeControl.GroupTime)
}
