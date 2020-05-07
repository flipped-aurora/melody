package proxy

import (
	"context"
	"melody/encoding"
	"net/http"
)

// HTTPResponseParserConfig 封装了解码器和格式化器，作用于proxy.Response
type HTTPResponseParserConfig struct {
	Decoder         encoding.Decoder
	EntityFormatter EntityFormatter
}

// NoOpHTTPResponseParser 不对http.Response做任何format，直接封装成proxy.Response
func NoOpHTTPResponseParser(ctx context.Context, resp *http.Response) (*Response, error) {
	return &Response{
		Data:       map[string]interface{}{},
		IsComplete: true,
		Io:         NewReadCloserWrapper(ctx, resp.Body),
		Metadata: Metadata{
			Headers:    resp.Header,
			StatusCode: resp.StatusCode,
		},
	}, nil
}

func DefaultHTTPResponseParserFactory(cfg HTTPResponseParserConfig) HTTPResponseParser {
	return func(ctx context.Context, resp *http.Response) (*Response, error) {
		defer resp.Body.Close()

		var data map[string]interface{}
		err := cfg.Decoder(resp.Body, &data)
		if err != nil {
			return nil, err
		}
		response := Response{
			Data:       data,
			IsComplete: true,
		}
		response = cfg.EntityFormatter.Format(response)
		return &response, nil
	}
}
