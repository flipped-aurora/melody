package bloomfilter

import (
	"context"
	"encoding/json"
	"errors"
	"melody/bf"
	rpc_bf "melody/bf/rpc"
	"melody/bf/rpc/server"
	"melody/config"
	"melody/logging"
	"net/http"
)

const Namespace = "melody_bloomfilter"

var (
	errNoConfig    = errors.New("no config for the bloomfilter")
	errWrongConfig = errors.New("invalid config for the bloomfilter")
)

type Config struct {
	rpc_bf.Config
	TokenKeys []string
	Headers   []string
}

type Rejecter struct {
	BF        bf.BloomFilter
	TokenKeys []string
	Headers   []string
}

// Register registers a bloomfilter given a config and registers the service with consul
func Register(ctx context.Context, serviceName string, cfg config.ServiceConfig,
	logger logging.Logger, register func(n string, p int)) (Rejecter, error) {

	data, ok := cfg.ExtraConfig[Namespace]
	if !ok {
		logger.Debug(errNoConfig.Error())
		return nopRejecter, errNoConfig
	}

	raw, err := json.Marshal(data)
	if err != nil {
		logger.Debug(errWrongConfig.Error())
		return nopRejecter, errWrongConfig
	}

	var rpcConfig Config
	if err := json.Unmarshal(raw, &rpcConfig); err != nil {
		logger.Debug(err.Error(), string(raw))
		return nopRejecter, err
	}

	rpcBF := server.New(ctx, rpcConfig.Config)
	register(serviceName, rpcConfig.Port)

	return Rejecter{
		BF:        rpcBF.Get(),
		TokenKeys: rpcConfig.TokenKeys,
		Headers:   rpcConfig.Headers,
	}, nil
}

func (r *Rejecter) RejectToken(claims map[string]interface{}) bool {
	for _, k := range r.TokenKeys {
		v, ok := claims[k]
		if !ok {
			continue
		}
		data, ok := v.(string)
		if !ok {
			continue
		}
		if r.BF.Check([]byte(k + "-" + data)) {
			return true
		}
	}
	return false
}

func (r *Rejecter) RejectHeader(header http.Header) bool {
	for _, k := range r.Headers {
		data := header.Get(k)
		if data == "" {
			continue
		}
		if r.BF.Check([]byte(k + "-" + data)) {
			return true
		}
	}
	return false
}

var nopRejecter = Rejecter{BF: new(bf.EmptySet)}
