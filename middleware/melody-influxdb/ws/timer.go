package ws

import (
	"regexp"
	"time"
)

type TimeControl struct {
	MinTime      string        `json:"min_time"`
	MaxTime      string        `json:"max_time"`
	TimeInterval string        `json:"time_interval"`
	GroupTime    string        `json:"group_time"`
	RefreshTime  time.Duration
	RefreshParam string `json:"refresh_time"`
}

var WsTimeControl TimeControl

func RegisterWSTimeControl() {
	WsTimeControl = TimeControl{
		MinTime:      "now()",
		MaxTime:      "now()",
		TimeInterval: "1h",
		GroupTime:    "5m",
		RefreshTime:  5 * time.Second,
	}
}

func SetTimeControl(timer TimeControl) {
	WsTimeControl = timer
}

var r =  regexp.MustCompile(`([0-9]*)([a-z])`)

const (
	hour = "15:04"
	day = "01-02 15:04"
)

func GetTimeFormat() string {
	if WsTimeControl.TimeInterval != "" {
		match := r.FindAllStringSubmatch(WsTimeControl.TimeInterval, -1)
		m := match[0]
		if len(m) == 3 {
			unit := m[2]
			if unit == "d" {
				return day
			}
		}
	}
	return hour
}
