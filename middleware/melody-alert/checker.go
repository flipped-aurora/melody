package alert

import (
	"errors"
	"melody/middleware/melody-alert/model"
	"strconv"
	"time"
)

type Checker func(endpoint, field string, data int64) error

func newChecker(fields map[string]interface{}) (Checker, error) {
	for field, threshold := range fields {
		if endpointFields, ok := threshold.(map[string]interface{}); ok {
			for k, v := range endpointFields {
				if _, ok := v.(string); !ok {
					return nil, errors.New("this threshold is not string")
				}
				if k == "time" {
					duration, err := time.ParseDuration(v.(string))
					if err != nil {
						return nil, err
					}
					endpointFields[k] = duration.Nanoseconds()
				} else {
					threshold, err := parseThreshold(v.(string))
					if err != nil {
						return nil, err
					}
					endpointFields[k] = threshold
				}
			}
			continue
		}

		if _, ok := threshold.(string); !ok {
			return nil, errors.New("this threshold is not string")
		}
		threshold, err := parseThreshold(threshold.(string))
		if err != nil {
			return nil, err
		}
		fields[field] = threshold
	}

	return func(endpoint, field string, data int64) error {
		fm := fields
		if endpoint != "" {
			if _, ok := fields[endpoint]; ok {
				tmp := fields[endpoint]
				fm = tmp.(map[string]interface{})
			}
		}
		if tmp, ok := fm[field]; ok {
			if data > tmp.(int64) {
				warning := model.Warning{
					Id:          model.Id.GetId(),
					Description: "警告：" + genWarningMessage(field, endpoint),
					TaskName:    field,
					CurValue:    data,
					Threshold:   tmp.(int64),
					Ctime:       time.Now().UnixNano() / 1e6,
					Handled:     0,
				}
				model.WarningList.Add(warning)
			}
		}
		return nil
	}, nil
}

func genWarningMessage(field string, endpoint string) string {
	msg := "无具体消息"
	switch field {
	case "numgc":
		msg = "GC次数超过阈值"
	case "sys":
		msg = "系统内存超过阈值"
	case "heapsys":
		msg = "堆内存超过阈值"
	case "stacksys":
		msg = "栈内存超过阈值"
	case "mcachesys":
		msg = "内存缓存池超过阈值"
	case "mspansys":
		msg = "某对象内存大小超过阈值"
	case "size":
		msg = endpoint + " 请求次数超过阈值"
	case "time":
		msg = endpoint + " 请求时间超过阈值"
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
