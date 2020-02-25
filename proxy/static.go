package proxy

import (
	"context"
	"melody/config"
)

const (
	staticKey = "static"

	staticAlwaysStrategy       = "always"
	staticIfSuccessStrategy    = "success"
	staticIfErroredStrategy    = "errored"
	staticIfCompleteStrategy   = "complete"
	staticIfIncompleteStrategy = "incomplete"
)

type staticConfig struct {
	Data     map[string]interface{}
	Strategy string
	Match    func(*Response, error) bool
}

func NewStaticDataMiddleware(endpoint *config.EndpointConfig) Middleware {
	v, ok := getStaticConfig(endpoint.ExtraConfig);
	if !ok {
		return EmptyMiddleware
	}

	return func(proxy ...Proxy) Proxy {
		if len(proxy) > 1 {
			panic(ErrTooManyProxies)
		}
		return func(ctx context.Context, request *Request) (response *Response, e error) {
			result, err := proxy[0](ctx, request)
			if !v.Match(result, err) {
				return result, err
			}

			if result == nil {
				result = &Response{Data: map[string]interface{}{}}
			}

			for k, v := range v.Data {
				result.Data[k] = v
			}

			return result, err
		}
	}
}

func getStaticConfig(extraConfig config.ExtraConfig) (staticConfig, bool) {
	v, ok := extraConfig[Namespace]
	if !ok {
		return staticConfig{}, ok
	}
	e, ok := v.(map[string]interface{})
	if !ok {
		return staticConfig{}, ok
	}
	v, ok = e[staticKey]
	if !ok {
		return staticConfig{}, ok
	}
	tmp, ok := v.(map[string]interface{})
	if !ok {
		return staticConfig{}, ok
	}
	data, ok := tmp["data"].(map[string]interface{})
	if !ok {
		return staticConfig{}, ok
	}

	name, ok := tmp["strategy"].(string)
	if !ok {
		name = staticAlwaysStrategy
	}
	cfg := staticConfig{
		Data:     data,
		Strategy: name,
		Match:    staticAlwaysMatch,
	}
	switch name {
	case staticAlwaysStrategy:
		cfg.Match = staticAlwaysMatch
	case staticIfSuccessStrategy:
		cfg.Match = staticIfSuccessMatch
	case staticIfErroredStrategy:
		cfg.Match = staticIfErroredMatch
	case staticIfCompleteStrategy:
		cfg.Match = staticIfCompleteMatch
	case staticIfIncompleteStrategy:
		cfg.Match = staticIfIncompleteMatch
	}
	return cfg, true
}

func staticAlwaysMatch(_ *Response, _ error) bool       { return true }
func staticIfSuccessMatch(_ *Response, err error) bool  { return err == nil }
func staticIfErroredMatch(_ *Response, err error) bool  { return err != nil }
func staticIfCompleteMatch(r *Response, err error) bool { return err == nil && r != nil && r.IsComplete }
func staticIfIncompleteMatch(r *Response, _ error) bool { return r == nil || !r.IsComplete }
