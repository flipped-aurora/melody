package register

import "sync"

const (
	// ScopeRequest 修改HTTP请求。
	ScopeRequest Scope = "request"
	// ScopeResponse 修改HTTP响应。
	ScopeResponse Scope = "response"
)

// Register 是包含所有martian组件的结构体
type Register map[string]Component

// Scope 定义组件的范围
type Scope string

// Component 组件包含范围和模块工厂
type Component struct {
	Scope       []Scope
	NewFromJSON func(b []byte) (interface{}, error)
}

var (
	register = Register{}
	mutex    = &sync.RWMutex{}
)

// Set 将接收到的数据添加到寄存器中
func Set(name string, scope []Scope, f func(b []byte) (interface{}, error)) {
	mutex.Lock()
	register[name] = Component{
		Scope:       scope,
		NewFromJSON: f,
	}
	mutex.Unlock()
}

// Get 获取寄存器的副本
func Get() Register {
	mutex.RLock()
	r := make(Register, len(register))
	for k, v := range register {
		r[k] = v
	}
	mutex.RUnlock()
	return r
}
