package martian

import (
	"melody/middleware/melody-martian/register"

	"github.com/google/martian/parse"
)

// Register 从 krakend-martian 寄存器获取所有修饰符，并将它们注册到 martian 解析器中
func Register() {
	for k, component := range register.Get() {
		parse.Register(k, func(b []byte) (*parse.Result, error) {
			v, err := component.NewFromJSON(b)
			if err != nil {
				return nil, err
			}

			return parse.NewResult(v, toModifierType(component.Scope))
		})
	}
}

func toModifierType(scopes []register.Scope) []parse.ModifierType {
	modifierType := make([]parse.ModifierType, len(scopes))
	for k, s := range scopes {
		modifierType[k] = parse.ModifierType(s)
	}
	return modifierType
}
