package internal

import "sync"

func New() *Untyped {
	return &Untyped{
		data:  map[string]interface{}{},
		mutex: &sync.RWMutex{},
	}
}

type Untyped struct {
	data  map[string]interface{}
	mutex *sync.RWMutex
}

func (u *Untyped) Register(name string, v interface{}) {
	u.mutex.Lock()
	u.data[name] = v
	u.mutex.Unlock()
}

func (u *Untyped) Get(name string) (interface{}, bool) {
	u.mutex.RLock()
	v, ok := u.data[name]
	u.mutex.RUnlock()
	return v, ok
}

func (u *Untyped) Clone() map[string]interface{} {
	u.mutex.RLock()
	clone := make(map[string]interface{}, len(u.data))
	for k, v := range u.data {
		clone[k] = v
	}
	u.mutex.RUnlock()
	return clone
}
