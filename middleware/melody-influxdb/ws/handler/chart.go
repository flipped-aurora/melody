package handler

import "melody/middleware/melody-influxdb/ws/convert"

/**
 将InfluxDB返回的结果集映射到数组中
times: 时间数组
data: InfluxDB返回的结果集
format: 时间的格式
lines: 其他数据数组 (一个或以上)
*/
func ResultDataHandler(times *[]string, data [][]interface{}, format string, lines ...*[]int64) {
	for _, v := range data { // v:[ time , data1 ,data2]
		if t, ok := convert.ObjectToStringTime(v[0], format); ok {
			*times = append(*times, t)
		}
		for index:= range lines {
			if g, ok := convert.ObjectToInt(v[index+1]); ok {
				*lines[index] = append(*lines[index], g)
			} else {
				*lines[index] = append(*lines[index], 0)
			}
		}
	}
}
