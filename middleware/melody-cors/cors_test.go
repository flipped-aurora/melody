package cors

import (
	"encoding/json"
	"testing"
	"time"
)

func TestConfigGetter(t *testing.T) {
	sampleCfg := map[string]interface{}{}
	serialized := []byte(`{ "melody_cors": {
			"allow_origins": [ "http://localhost", "http://www.example.com" ],
			"allow_headers": [ "X-Test", "Content-Type"],
			"allow_methods": [ "POST", "GET" ],
			"expose_headers": [ "Content-Type" ],
			"allow_credentials": false,
			"max_age": "24h"
			}
		}`)
	json.Unmarshal(serialized, &sampleCfg)
	testCfg := GetConfig(sampleCfg).(Config)

	if len(testCfg.AllowOrigins) != 2 {
		t.Error("Should have exactly 2 allowed origins.\n")
	}
	for i, v := range []string{"http://localhost", "http://www.example.com"} {
		if testCfg.AllowOrigins[i] != v {
			t.Errorf("Invalid value %s should be %s\n", testCfg.AllowOrigins[i], v)
		}
	}
	if len(testCfg.AllowHeaders) != 2 {
		t.Error("Should have exactly 2 allowed headers.\n")
	}
	for i, v := range []string{"X-Test", "Content-Type"} {
		if testCfg.AllowHeaders[i] != v {
			t.Errorf("Invalid value %s should be %s\n", testCfg.AllowHeaders[i], v)
		}
	}
	if len(testCfg.AllowMethods) != 2 {
		t.Error("Should have exactly 2 allowed headers.\n")
	}
	for i, v := range []string{"POST", "GET"} {
		if testCfg.AllowMethods[i] != v {
			t.Errorf("Invalid value %s should be %s\n", testCfg.AllowMethods[i], v)
		}
	}
	if len(testCfg.ExposeHeaders) != 1 {
		t.Error("Should have exactly 2 allowed headers.\n")
	}
	for i, v := range []string{"Content-Type"} {
		if testCfg.ExposeHeaders[i] != v {
			t.Errorf("Invalid value %s should be %s\n", testCfg.ExposeHeaders[i], v)
		}
	}
	if testCfg.AllowCredentials {
		t.Error("Allow Credentials should be disabled.\n")
	}

	if testCfg.MaxAge != 24*time.Hour {
		t.Errorf("Unexpected collection time: %v\n", testCfg.MaxAge)
	}
}

func TestDefaultConfiguration(t *testing.T) {
	sampleCfg := map[string]interface{}{}
	serialized := []byte(`{ "melody_cors": {
			"allow_origins": [ "http://www.example.com" ]
	}}`)
	json.Unmarshal(serialized, &sampleCfg)
	defaultCfg := GetConfig(sampleCfg).(Config)
	if defaultCfg.AllowOrigins[0] != "http://www.example.com" {
		t.Error("Wrong AllowOrigin.\n")
	}
}

func TestWrongConfiguration(t *testing.T) {
	sampleCfg := map[string]interface{}{}
	if _, ok := GetConfig(sampleCfg).(Config); ok {
		t.Error("The config should be nil\n")
	}
	badCfg := map[string]interface{}{Namespace: "test"}
	if _, ok := GetConfig(badCfg).(Config); ok {
		t.Error("The config should be nil\n")
	}
}

func TestEmptyConfiguration(t *testing.T) {
	noOriginCfg := map[string]interface{}{}
	serialized := []byte(`{ "melody_cors": {
			}
		}`)
	json.Unmarshal(serialized, &noOriginCfg)
	if v, ok := GetConfig(noOriginCfg).(Config); !ok {
		t.Errorf("The configuration should not be empty: %v\n", v)
	}
}
