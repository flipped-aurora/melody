package gin

import (
	"io"
	"melody/config"
	"melody/encoding"
	"melody/proxy"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

// NEGOTIATE 定义了输出的编码方式
const NEGOTIATE = "negotiate"

var (
	mutex = &sync.RWMutex{}
	emptyResponse = gin.H{}
	renderRegister = map[string]Render{
		NEGOTIATE:       negotiatedRender,
		encoding.STRING: stringRender,
		encoding.JSON:   jsonRender,
		encoding.NOOP:   noopRender,
	}
)

// Render 作为Response的渲染器
type Render func(*gin.Context, *proxy.Response)

func getRender(cfg *config.EndpointConfig) Render {
	defaultRender := jsonRender
	// 如果只有一个Backends，将该Backends的编码作为响应编码
	if len(cfg.Backends) == 1{
		defaultRender = getWithDefault(cfg.Backends[0].Encoding, defaultRender)
	}

	if cfg.OutputEncoding == "" {
		return defaultRender
	}
	// Endpoint层级的编码优先级大于Backends的编码格式的优先级
	return getWithDefault(cfg.OutputEncoding, defaultRender)
}

func getWithDefault(key string, defa Render) Render {
	mutex.RLock()
	v, ok := renderRegister[key]
	mutex.RUnlock()
	if !ok {
		return defa
	}
	return v
}

func registerRender(key string, render Render) {
	mutex.Lock()
	renderRegister[key] = render
	mutex.Unlock()
}

func negotiatedRender(c *gin.Context, response *proxy.Response) {
	switch c.NegotiateFormat(gin.MIMEJSON, gin.MIMEPlain, gin.MIMEXML) {
	case gin.MIMEXML:
		xmlRender(c, response)
	case gin.MIMEPlain:
		yamlRender(c, response)
	default:
		jsonRender(c, response)
	}
}

func stringRender(c *gin.Context, response *proxy.Response) {
	status := c.Writer.Status()

	if response == nil {
		c.String(status, "")
		return
	}
	//TODO 选择string render的时候可以定制key
	d, ok := response.Data["content"]
	if !ok {
		c.String(status, "")
		return
	}
	msg, ok := d.(string)
	if !ok {
		c.String(status, "")
		return
	}
	c.String(status, msg)
}


func yamlRender(c *gin.Context, response *proxy.Response) {
	status := c.Writer.Status()
	if response == nil {
		c.YAML(status, emptyResponse)
		return
	}
	c.YAML(status, response.Data)
}

func xmlRender(c *gin.Context, response *proxy.Response) {
	status := c.Writer.Status()
	if response == nil {
		c.XML(status, nil)
		return
	}
	//TODO 选择string render的时候可以定制key
	d, ok := response.Data["content"]
	if !ok {
		c.XML(status, nil)
		return
	}
	c.XML(status, d)
}

func jsonRender(c *gin.Context, resp *proxy.Response) {
	status := c.Writer.Status()
	if resp == nil {
		c.JSON(status, emptyResponse)
		return
	}
	c.JSON(status, resp.Data)
}

func noopRender(c *gin.Context, response *proxy.Response) {
	if response == nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Status(response.Metadata.StatusCode)
	for k, vs := range response.Metadata.Headers {
		for _, v := range vs {
			c.Writer.Header().Add(k, v)
		}
	}
	if response.Io == nil {
		return
	}
	io.Copy(c.Writer, response.Io)
}

