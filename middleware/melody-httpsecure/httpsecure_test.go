package httpsecure

import (
	"testing"
)

func Test_SetStrings(t *testing.T) {
	m := map[string]interface{}{"hosts": []interface{}{"127.0.0.1", "192.168.1.1"}}
	a := []string{}
	setStrings(m, "hosts", &a)
	if len(a) != 2 || a[0] != "127.0.0.1" {
		t.Error("set string array error", a)
	}
}

func Test_SetString(t *testing.T) {
	m := map[string]interface{}{"header": "yes"}
	var a string
	setString(m, "header", &a)
	if a != "yes" {
		t.Error("set string error", a)
	}
}
