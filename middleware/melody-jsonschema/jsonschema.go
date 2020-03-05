package jsonschema

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"melody/config"
	"melody/proxy"
	"strings"

	"github.com/xeipuuv/gojsonschema"
)

const Namespace = "melody_jsonschema"

var ErrEmptyBody = errors.New("could not validate an empty body")

// ProxyFactory ...
func ProxyFactory(pf proxy.Factory) proxy.FactoryFunc {
	return proxy.FactoryFunc(func(cfg *config.EndpointConfig) (proxy.Proxy, error) {
		next, err := pf.New(cfg)
		if err != nil {
			return proxy.NoopProxy, err
		}
		schemaLoader, ok := configGetter(cfg.ExtraConfig).(gojsonschema.JSONLoader)
		if !ok || schemaLoader == nil {
			return next, nil
		}
		return newProxy(schemaLoader, next), nil
	})
}

func newProxy(schemaLoader gojsonschema.JSONLoader, next proxy.Proxy) proxy.Proxy {
	return func(ctx context.Context, r *proxy.Request) (*proxy.Response, error) {
		if r.Body == nil {
			return nil, ErrEmptyBody
		}
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return nil, err
		}
		r.Body.Close()
		r.Body = ioutil.NopCloser(bytes.NewBuffer(body))

		result, err := gojsonschema.Validate(schemaLoader, gojsonschema.NewBytesLoader(body))
		if err != nil {
			return nil, err
		}
		if !result.Valid() {
			return nil, &validationError{errs: result.Errors()}
		}

		return next(ctx, r)
	}
}

func configGetter(cfg config.ExtraConfig) interface{} {
	v, ok := cfg[Namespace]
	if !ok {
		return nil
	}
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(v); err != nil {
		return nil
	}
	return gojsonschema.NewBytesLoader(buf.Bytes())
}

type validationError struct {
	errs []gojsonschema.ResultError
}

func (v *validationError) Error() string {
	errs := make([]string, len(v.errs))
	for i, desc := range v.errs {
		errs[i] = fmt.Sprintf("- %s", desc)
	}
	return strings.Join(errs, "\n")
}
