package alert

import (
	"errors"
	"melody/middleware/melody-alert/model"
	"strconv"
	"time"
)

type Checker func(field string, data []int64) ([]model.Warning, error)

func newChecker(fields map[string]interface{}) (Checker, error) {
	fieldsMap := make(map[string]int64)
	for field, threshold := range fields {
		if field == "api" {
			continue
		}

		if _, ok := threshold.(string); !ok {
			return nil, errors.New("this threshold is not string")
		}
		threshold, err := parseThreshold(threshold.(string))
		if err != nil {
			return nil, err
		}
		fieldsMap[field] = threshold
	}

	return func(field string, data []int64) ([]model.Warning, error) {
		warnings := make([]model.Warning, 0)
		for _, item := range data {
			if item > fieldsMap[field] {
				warning := model.Warning{
					Description: "警告: " + genWarningMessage(field, fields),
					TaskName:    field + "任务",
					CurValue:    item,
					Threshold:   fieldsMap[field],
					Ctime:       time.Now().UnixNano(),
					Handled:     0,
				}
				warnings = append(warnings, warning)
				// TODO 通知前端
			}
		}
		return warnings, nil
	}, nil
}

func genWarningMessage(field string, fields map[string]interface{}) string {
	msg := "无具体消息"
	switch field {
	case "gc_num":
		msg = "GC次数超过阈值"
	case "sys_memory":
		msg = "系统内存超过阈值"
	case "heap_memory":
		msg = "堆内存超过阈值"
	case "stack_memory":
		msg = "栈内存超过阈值"
	case "size":
		if _, ok := fields["api"].(string); !ok {
			msg = "某接口请求次数超过阈值"
		} else {
			msg = fields["api"].(string) + " 请求次数超过阈值"
		}
	case "time":
		if _, ok := fields["api"].(string); !ok {
			msg = "某接口请求次数超过阈值"
		} else {
			msg = fields["api"].(string) + " 请求时间超过阈值"
		}
	}
	return msg
}

func parseThreshold(t string) (int64, error) {
	if IsNum(string(t[len(t)-1])) {
		i, err := strconv.ParseInt(t, 10, 0)
		if err != nil {
			return 0, err
		}
		return i, nil
	} else {
		i, err := strconv.ParseInt(t[:len(t)-1], 10, 0)
		if err != nil {
			return 0, err
		}
		switch t[len(t)-1] {
		case 'k':
			return i * 1000, nil
		case 'm':
			return i * 1000 * 1000, nil
		}

		return 0, errors.New("没有这个单位: " + string(t[len(t)-1]))
	}
}

func IsNum(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}
