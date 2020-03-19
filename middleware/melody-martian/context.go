package martian

import "context"

// NewContext 返回一个包装了接收到的父类的上下文
func NewContext(parent context.Context) *Context {
	return &Context{
		Context: parent,
	}
}

// Context 提供单个请求/响应对的信息。
type Context struct {
	context.Context
	skipRoundTrip bool
}

// SkipRoundTrip 标记上下文以跳过往返
func (c *Context) SkipRoundTrip() {
	c.skipRoundTrip = true
}

// SkippingRoundTrip 返回跳过往返旅程的标志
func (c *Context) SkippingRoundTrip() bool {
	return c.skipRoundTrip
}

var _ context.Context = &Context{Context: context.Background()}
