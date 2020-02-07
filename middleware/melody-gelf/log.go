package gelf

import (
	"errors"
	"fmt"
	"gopkg.in/Graylog2/go-gelf.v2/gelf"
	"io"
	"melody/config"
)

const Namespace = "melody_gelf"

var (
	tcpWriter = gelf.NewTCPWriter
	udpWriter = gelf.NewUDPWriter
	ErrorConfig = fmt.Errorf("not found extra config about melody-gelf module")
	ErrorMissAddress = errors.New("miss gelf address to send log")
)

type Config struct {
	Addr string
	TCPEnable bool
}



func NewWriter(cfg config.ExtraConfig) (io.Writer, error) {
	g, ok := GetConfig(cfg).(Config)
	if !ok {
		return nil, ErrorConfig
	}

	if g.Addr == "" {
		return nil, ErrorMissAddress
	}

	if g.TCPEnable {
		return tcpWriter(g.Addr)
	}

	return udpWriter(g.Addr)
}

func GetConfig(cfg config.ExtraConfig) interface{} {
	v, ok := cfg[Namespace]

	if !ok {
		return nil
	}

	t, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}

	addr, ok := t["addr"]
	if !ok {
		return nil
	}

	enable, ok := t["tcp_enable"]
	if !ok {
		return nil
	}

	return Config{
		Addr:      addr.(string),
		TCPEnable: enable.(bool),
	}
}



