package gin

import (
	"context"
	"fmt"
	"melody/config"
	"melody/core"
	"melody/proxy"
	"melody/router"
	"net/textproto"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	passAllRequestHeaders = "*"
	passAllQueryParams    = "*"
)

// HandlerFactory 返回Endpoint层的Handler工厂
type HandlerFactory func(*config.EndpointConfig, proxy.Proxy) gin.HandlerFunc

func EndpointHandler(config *config.EndpointConfig, proxy proxy.Proxy) gin.HandlerFunc {
	return CustomErrorEndpointHandler(config, proxy, router.DefaultToHTTPError)
}

// CustomErrorEndpointHandler 作为HandlerFactory接口的默认实现
func CustomErrorEndpointHandler(config *config.EndpointConfig, proxy proxy.Proxy, errF router.ToHTTPError) gin.HandlerFunc {
	cacheControlHeader := fmt.Sprintf("public, max-age=%d", int(config.CacheTTL.Seconds()))
	isCacheEnable := config.CacheTTL.Seconds() != 0
	request := NewRequest(config.HeadersToPass)
	responseRender := getRender(config)

	return func(c *gin.Context) {
		reqCtx, cancel := context.WithTimeout(c, config.Timeout)
		c.Header(core.MelodyHeaderKey, core.MelodyHeaderValue)
		// 执行代理 *
		response, err := proxy(reqCtx, request(c, config.QueryString))

		select {
		case <-reqCtx.Done():
			if err == nil {
				err = router.ErrorInternalError
			}
		default:
		}

		// 默认将complete状态置为false
		complete := router.HeaderInCompleteResponseValue

		if response != nil && len(response.Data) > 0 {
			if response.IsComplete {
				complete = router.HeaderCompleteResponseValue
				if isCacheEnable {
					c.Header("Cache-Control", cacheControlHeader)
				}
			}

			// 将Backends层代理回来的请求头，写入Endpoints层的response
			for k, vs := range response.Metadata.Headers {
				for _, v := range vs {
					c.Writer.Header().Add(k, v)
				}
			}
		}


		// 设置最终代理是否完成校验响应头
		c.Header(router.HeaderCompleteKey, complete)

		// 校验响应是否发生err
		if err != nil {
			c.Error(err)
			// 校验响应是否为nil
			if response == nil {
				if t, ok := err.(responseError); ok {
					c.Status(t.StatusCode())
				} else {
					c.Status(errF(err))
				}
				cancel()
				return
			}
		}
		// 去render成最终的编码格式
		responseRender(c, response)
		// call cancel去关闭本次req context
		cancel()
	}
}

// NewRequest 从context中提取新的请求
func NewRequest(passHeaders []string) func(*gin.Context, []string) *proxy.Request {
	if len(passHeaders) == 0 {
		passHeaders = router.PassHeaders
	}

	return func(ctx *gin.Context, queryString []string) *proxy.Request {
		// handle url params
		params := make(map[string]string, len(ctx.Params))
		for _, param := range ctx.Params {
			params[strings.Title(param.Key)] = param.Value
		}
		// handle request headers
		headers := make(map[string][]string, len(passHeaders)+2)
		for _, k := range passHeaders {
			if k == passAllRequestHeaders {
				headers = ctx.Request.Header
				break
			}

			if h, ok := ctx.Request.Header[textproto.CanonicalMIMEHeaderKey(k)]; ok {
				headers[k] = h
			}
		}
		// 添加自定义请求头
		headers["X-Forwarded-For"] = []string{ctx.ClientIP()}
		if _, ok := headers["User-Agent"]; !ok {
			headers["User-Agent"] = router.UserAgentHeaderValue
		} else {
			headers["X-Forwarded-Via"] = router.UserAgentHeaderValue
		}
		// 封装 query map
		query := make(map[string][]string, len(queryString))
		queryValues := ctx.Request.URL.Query()
		for i := range queryString {
			if queryString[i] == passAllQueryParams {
				query = ctx.Request.URL.Query()
			}

			if v, ok := queryValues[queryString[i]]; ok && len(v) > 0 {
				query[queryString[i]] = v
			}
		}

		return &proxy.Request{
			Method:  ctx.Request.Method,
			Query:   query,
			Body:    ctx.Request.Body,
			Params:  params,
			Headers: headers,
		}
	}
}

type responseError interface {
	error
	StatusCode() int
}
