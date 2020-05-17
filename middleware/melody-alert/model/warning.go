package model

import (
	"fmt"
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
	Warnings []Warning
	Lock     sync.RWMutex
}

func (ws *Warnings) Add(warning Warning) {
	ws.Lock.Lock()
	ws.Warnings = append(ws.Warnings, warning)
	fmt.Printf("%+v\n", warning)
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
