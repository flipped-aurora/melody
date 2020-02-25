package proxy

import (
	"context"
	"errors"
	"melody/config"
	"time"
)

var errNullResult = errors.New("invalid response")

func NewConcurrentCallMiddleware(backend *config.Backend) Middleware {
	// Check proxy
	if backend.ConcurrentCalls == 1 {
		panic(ErrTooManyProxies)
	}

	serviceTimeout := time.Duration(75*backend.Timeout.Nanoseconds()/100) * time.Nanosecond

	return func(proxy ...Proxy) Proxy {
		if len(proxy) > 1 {
			panic(ErrTooManyProxies)
		}

		return func(ctx context.Context, request *Request) (*Response, error) {
			totalCtx, cancel := context.WithTimeout(ctx, serviceTimeout)

			results := make(chan *Response, backend.ConcurrentCalls)
			failed := make(chan error, backend.ConcurrentCalls)

			for i := 0; i < backend.ConcurrentCalls; i++ {
				go processConcurrentCall(totalCtx, proxy[0], request, results, failed)
			}

			var response *Response
			var err error

			for i := 0; i < backend.ConcurrentCalls; i++ {
				select {
				case response = <-results:
					if response != nil && response.IsComplete {
						cancel()
						return response, nil
					}
				case err = <-failed:
				case <-ctx.Done():
				}
			}
			cancel()
			return response, err
		}
	}
}

func processConcurrentCall(ctx context.Context, proxy Proxy, request *Request, responses chan *Response, errors chan error) {
	localCtx, cancel := context.WithCancel(ctx)

	resp, err := proxy(localCtx, request)

	if err != nil {
		errors <- err
		cancel()
		return
	}

	if resp == nil {
		errors <- errNullResult
		cancel()
		return
	}

	select {
	case responses <- resp:
	case <-ctx.Done():
		errors <- ctx.Err()
	}
	cancel()
}
