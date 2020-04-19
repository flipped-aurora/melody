package ws

import (
	"errors"
	"github.com/influxdata/influxdb/client/v2"
)

func (wsc WebSocketClient) executeQuery(cmd string) ([]client.Result, error) {
	resp, err := wsc.Client.Query(client.NewQuery(cmd, wsc.DB, "s"))
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, errors.New("error: influx query no resp")
	}
	if resp.Err != "" {
		return nil, errors.New(resp.Err)
	}
	if len(resp.Results) == 0 {
		return nil, errors.New("error: influx query no resp")
	}
	return resp.Results, nil
}

func NormalExecuteQuery(c client.Client, cmd, db string) ([]client.Result, error) {
	resp, err := c.Query(client.NewQuery(cmd, db, "s"))
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, errors.New("error: influx query no resp")
	}
	if resp.Err != "" {
		return nil, errors.New(resp.Err)
	}
	if len(resp.Results) == 0 {
		return nil, errors.New("error: influx query no resp")
	}
	return resp.Results, nil
}
