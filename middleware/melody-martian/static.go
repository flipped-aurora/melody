package martian

import (
	"encoding/json"
	"net/http"

	"github.com/google/martian/parse"
	"github.com/google/martian/static"
)

// StaticModifier is a martian.
type StaticModifier struct {
	*static.Modifier
}

type staticJSON struct {
	ExplicitPaths map[string]string    `json:"explicitPaths"`
	RootPath      string               `json:"rootPath"`
	Scope         []parse.ModifierType `json:"scope"`
}

// NewStaticModifier 构造一个static.Modifier，它采用一个路径来服务文件，以及一个可选的请求路径到本地文件路径的映射(仍然以rootPath为根)。
func NewStaticModifier(rootPath string) *StaticModifier {
	return &StaticModifier{
		Modifier: static.NewModifier(rootPath),
	}
}

// ModifyRequest 将上下文标记为跳过往返，并将所有https请求降级为http。
func (s *StaticModifier) ModifyRequest(req *http.Request) error {
	ctx := NewContext(req.Context())
	ctx.SkipRoundTrip()

	if req.URL.Scheme == "https" {
		req.URL.Scheme = "http"
	}

	*req = *req.WithContext(ctx)

	return nil
}

func staticModifierFromJSON(b []byte) (*parse.Result, error) {
	msg := &staticJSON{}
	if err := json.Unmarshal(b, msg); err != nil {
		return nil, err
	}

	mod := NewStaticModifier(msg.RootPath)
	mod.SetExplicitPathMappings(msg.ExplicitPaths)
	return parse.NewResult(mod, msg.Scope)
}
