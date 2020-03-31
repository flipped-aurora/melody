package ws

import "time"

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
