package gin

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"melody/config"
	"melody/logging"
	melody "melody/middleware/melody-botmonitor/melody"
	"melody/proxy"
	melodygin "melody/router/gin"
)

func TestRegister(t *testing.T) {
	gin.SetMode(gin.TestMode)

	engine := gin.New()

	cfg := config.ServiceConfig{
		ExtraConfig: config.ExtraConfig{
			melody.Namespace: map[string]interface{}{
				"blacklist": []string{"a", "b"},
				"whitelist": []string{"c", "Pingdom.com_bot_version_1.1"},
				"patterns": []string{
					`(Pingdom.com_bot_version_)(\d+)\.(\d+)`,
					`(facebookexternalhit)/(\d+)\.(\d+)`,
				},
			},
		},
	}

	Register(cfg, logging.NoOp, engine)

	engine.GET("/", func(c *gin.Context) {
		c.String(200, "hi!")
	})

	if err := testDetection(engine); err != nil {
		t.Error(err)
	}
}

func TestNew(t *testing.T) {
	gin.SetMode(gin.TestMode)

	engine := gin.New()

	cfg := &config.EndpointConfig{
		ExtraConfig: config.ExtraConfig{
			melody.Namespace: map[string]interface{}{
				"blacklist": []string{"a", "b"},
				"whitelist": []string{"c", "Pingdom.com_bot_version_1.1"},
				"patterns": []string{
					`(Pingdom.com_bot_version_)(\d+)\.(\d+)`,
					`(facebookexternalhit)/(\d+)\.(\d+)`,
				},
			},
		},
	}

	proxyfunc := func(_ context.Context, _ *proxy.Request) (*proxy.Response, error) {
		return &proxy.Response{IsComplete: true}, nil
	}

	engine.GET("/", New(melodygin.EndpointHandler, logging.NoOp)(cfg, proxyfunc))

	if err := testDetection(engine); err != nil {
		t.Error(err)
	}
}

func testDetection(engine *gin.Engine) error {
	for i, ua := range []string{
		"abcd",
		"",
		"c",
		"Pingdom.com_bot_version_1.1",
	} {
		req, _ := http.NewRequest("GET", "http://example.com", nil)
		req.Header.Add("User-Agent", ua)

		w := httptest.NewRecorder()
		engine.ServeHTTP(w, req)

		if w.Result().StatusCode != 200 {
			return fmt.Errorf("the req #%d has been detected as a bot: %s", i, ua)
		}
	}

	for i, ua := range []string{
		"a",
		"b",
		"facebookexternalhit/1.1",
		"Pingdom.com_bot_version_1.2",
	} {
		req, _ := http.NewRequest("GET", "http://example.com", nil)
		req.Header.Add("User-Agent", ua)

		w := httptest.NewRecorder()
		engine.ServeHTTP(w, req)

		if w.Result().StatusCode != http.StatusForbidden {
			return fmt.Errorf("the req #%d has not been detected as a bot: %s", i, ua)
		}
	}
	return nil
}
