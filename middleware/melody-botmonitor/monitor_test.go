package botmonitor

import (
	"fmt"
	"net/http"
	"testing"
)

func TestNew_noLRU(t *testing.T) {
	d, err := New(Config{
		Blacklist: []string{"a", "b"},
		Whitelist: []string{"c", "Pingdom.com_bot_version_1.1"},
		Patterns: []string{
			`(Pingdom.com_bot_version_)(\d+)\.(\d+)`,
			`(facebookexternalhit)/(\d+)\.(\d+)`,
		},
	})
	if err != nil {
		t.Error(err)
		return
	}

	if err := testDetection(d); err != nil {
		t.Error(err)
	}
}

func TestNew_LRU(t *testing.T) {
	d, err := New(Config{
		Blacklist: []string{"a", "b"},
		Whitelist: []string{"c", "Pingdom.com_bot_version_1.1"},
		Patterns: []string{
			`(Pingdom.com_bot_version_)(\d+)\.(\d+)`,
			`(facebookexternalhit)/(\d+)\.(\d+)`,
		},
		CacheSize: 10000,
	})
	if err != nil {
		t.Error(err)
		return
	}

	if err := testDetection(d); err != nil {
		t.Error(err)
	}
}

func testDetection(f DetectorFunc) error {
	for i, ua := range []string{
		"abcd",
		"",
		"c",
		"Pingdom.com_bot_version_1.1",
	} {
		req, _ := http.NewRequest("GET", "http://example.com", nil)
		req.Header.Add("User-Agent", ua)
		if f(req) {
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
		if !f(req) {
			return fmt.Errorf("the req #%d has not been detected as a bot: %s", i, ua)
		}
	}
	return nil
}
