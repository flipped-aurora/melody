package ws

import (
	"errors"
	"github.com/influxdata/influxdb/client/v2"
)

func ExecuteQuery(c client.Client, cmd, db string) ([]client.Result, error) {
	resp, err := c.Query(client.NewQuery(cmd, db, "s"))
	if err != nil || resp.Err != "" {
		err = errors.New(resp.Err)
		return nil, err
	}
	return resp.Results, nil
}
