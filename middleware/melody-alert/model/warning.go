package model

import "sync"

type Warning struct {
	Id          int    `gorm:"id" json:"id"`
	Description string `gorm:"description" json:"description"`
	TaskName    string `gorm:"task_name" json:"task_name"`
	CurValue    int64  `gorm:"cur_value" json:"cur_value"`
	Threshold   int64  `gorm:"threshold" json:"threshold"`
	Ctime       int64  `gorm:"ctime" json:"ctime"`
	Handled     int    `gorm:"handled" json:"handled"`
}

type Warnings struct {
	Warnings []Warning
	Lock     sync.RWMutex
}

func (ws *Warnings) Add(warning Warning) {

}
