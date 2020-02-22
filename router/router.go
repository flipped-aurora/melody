package router

import (
	"context"
	"melody/config"
	http "melody/transport/http/server"
)

const (
	HeaderCompleteResponseValue = http.HeaderCompleteResponseValue
	HeaderInCompleteResponseValue = http.HeaderIncompleteResponseValue
	HeaderCompleteKey = http.HeaderCompleteKey
)

var (
	// DefaultRunServer default run server func
	DefaultRunServer = http.RunServer
	// PassHeaders 默认放行的请求头
	PassHeaders = http.HeadersToSend
	// UserAgentHeaderValue 代理的请求头标签值
	UserAgentHeaderValue = http.UserAgentHeaderValue
	ErrorInternalError = http.ErrorInternalError
	DefaultToHTTPError = http.DefaultToHTTPError
)

// Router 暴露出去的接口
type Router interface {
	Run(config.ServiceConfig)
}

// ToHTTPError change  error -> http status code
type ToHTTPError http.ToHTTPError

// Factory 暴露出去的接口
type Factory interface {
	New() Router
	NewWithContext(context.Context) Router
}

//(¬‿¬)
