package ws

import (
	"errors"
	"github.com/influxdata/influxdb/client/v2"
)

func (wsc WebSocketClient) executeQuery( cmd string) ([]client.Result, error) {
	resp, err := wsc.Client.Query(client.NewQuery(cmd, wsc.DB, "s"))
	if err != nil || resp.Err != "" {
		err = errors.New(resp.Err)
		return nil, err
	}
	return resp.Results, nil
}

func NormalExecuteQuery(c client.Client, cmd, db string) ([]client.Result, error) {
	resp, err := c.Query(client.NewQuery(cmd, db, "s"))
	if err != nil || resp.Err != "" {
		err = errors.New(resp.Err)
		return nil, err
	}
	return resp.Results, nil
}
