package model

import (
	"sync"
)

type Warning struct {
	Id          int64  `gorm:"id" json:"id"`
	Description string `gorm:"description" json:"description"`
	TaskName    string `gorm:"task_name" json:"task_name"`
	CurValue    int64  `gorm:"cur_value" json:"cur_value"`
	Threshold   int64  `gorm:"threshold" json:"threshold"`
	Ctime       int64  `gorm:"ctime" json:"ctime"`
	Handled     int    `gorm:"handled" json:"handled"`
}

var (
	Id          = new(IdWorker)
	WarningList = new(Warnings)
)

type Warnings struct {
	Warnings []Warning    `json:"warnings"`
	Lock     sync.RWMutex `json:"-"`
}

func (ws *Warnings) Add(warning Warning) {
	ws.Lock.Lock()
	ws.Warnings = append(ws.Warnings, warning)
	ws.Lock.Unlock()
}

func (ws *Warnings) ChangeStatus(id int64) {
	ws.Lock.Lock()
	if ws.Warnings[id-1].Handled == 0 {
		ws.Warnings[id-1].Handled = 1
	} else {
		ws.Warnings[id-1].Handled = 0
	}
	ws.Lock.Unlock()
}

type IdWorker struct {
	Id   int64
	Lock sync.RWMutex
}

func (id *IdWorker) inc() {
	id.Lock.Lock()
	id.Id++
	id.Lock.Unlock()
}

func (id *IdWorker) GetId() int64 {
	id.inc()
	return id.Id
}
