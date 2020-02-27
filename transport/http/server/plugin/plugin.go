package plugin

import (
	"context"
	"melody/register"
	"net/http"
)

var serverRegister = register.NewSpaces()

// RegisterHandler ...
func RegisterHandler(
	name string,
	handler func(context.Context, map[string]interface{}, http.Handler) (http.Handler, error),
) {
	serverRegister.Register(Namespace, name, handler)
}
