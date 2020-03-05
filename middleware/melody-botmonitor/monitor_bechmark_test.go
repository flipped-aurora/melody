package botmonitor

import (
	"net/http"
	"testing"
)

func BenchmarkDetector(b *testing.B) {
	d, err := New(Config{
		Blacklist: []string{"a", "b"},
		Whitelist: []string{"c", "Pingdom.com_bot_version_1.1"},
		Patterns: []string{
			`(Pingdom.com_bot_version_)(\d+)\.(\d+)`,
			`(facebookexternalhit)/(\d+)\.(\d+)`,
		},
	})
	if err != nil {
		b.Error(err)
		return
	}

	becnhDetection(b, d)
}

func BenchmarkLRUDetector(b *testing.B) {
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
		b.Error(err)
		return
	}

	becnhDetection(b, d)
}

func becnhDetection(b *testing.B, f DetectorFunc) {
	for _, tc := range []struct {
		name string
		ua   string
	}{
		{"ok_1", "abcd"},
		{"ok_2", ""},
		{"ok_3", "c"},
		{"ok_4", "Pingdom.com_bot_version_1.1"},
		{"ko_1", "a"},
		{"ko_2", "b"},
		{"ko_3", "facebookexternalhit/1.1"},
		{"ko_4", "Pingdom.com_bot_version_1.2"},
	} {

		req, _ := http.NewRequest("GET", "http://example.com", nil)
		req.Header.Add("User-Agent", tc.ua)
		b.Run(tc.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				f(req)
			}
		})
	}
}
