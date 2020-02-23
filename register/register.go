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

func NewSpaces() *Namespaced {
	return &Namespaced{New()}
}

type Namespaced struct {
	data Untyped
}

func (n *Namespaced) Get(namespace string) (Untyped, bool) {
	v, ok := n.data.Get(namespace)
	if !ok {
		return nil, ok
	}
	register, ok := v.(Untyped)
	return register, ok
}

func (n *Namespaced) Register(namespace, name string, v interface{}) {
	if register, ok := n.Get(namespace); ok {
		register.Register(name, v)
		return
	}

	register := New()
	register.Register(name, v)
	n.data.Register(namespace, register)
}

func (n *Namespaced) AddNamespace(namespace string) {
	if _, ok := n.Get(namespace); ok {
		return
	}
	n.data.Register(namespace, New())
}
