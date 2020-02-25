package proxy

import (
	"context"
	"melody/sd"
	"net/url"
	"strings"
)

func NewLoadBalancedMiddlewareWithSubscriber(subscriber sd.Subscriber) Middleware {
	return newLoadBalancedMiddleware(sd.NewBalancer(subscriber))
}

func newLoadBalancedMiddleware(lb sd.Balancer) Middleware {
	return func(next ...Proxy) Proxy {
		if len(next) > 1 {
			panic(ErrTooManyProxies)
		}
		return func(ctx context.Context, request *Request) (*Response, error) {
			host, err := lb.Host()
			if err != nil {
				return nil, err
			}
			r := request.Clone()

			var b strings.Builder
			b.WriteString(host)
			b.WriteString(r.Path)
			r.URL, err = url.Parse(b.String())
			if err != nil {
				return nil, err
			}
			if len(r.Query) > 0 {
				r.URL.RawQuery += "&" + r.Query.Encode()
			}

			return next[0](ctx, &r)
		}
	}
}