package register

import "melody/register/internal"

type Untyped interface {
	Register(name string, v interface{})
	Get(name string) (interface{}, bool)
	Clone() map[string]interface{}
}

func New() Untyped {
	return internal.New()
}
