package counter

import (
	"github.com/influxdata/influxdb/client/v2"
	"melody/logging"
	alert "melody/middleware/melody-alert"
	"regexp"
	"strings"
	"sync"
	"time"
)

var (
	lastRequestCount  = map[string]int{}
	lastResponseCount = map[string]int{}
	mu                = new(sync.Mutex)
)

func Points(hostname string, now time.Time, counters map[string]int64, logger logging.Logger, checker alert.Checker) []*client.Point {
	points := requestPoints(hostname, now, counters, logger, checker)
	points = append(points, responsePoints(hostname, now, counters, logger, checker)...)
	points = append(points, connectionPoints(hostname, now, counters, logger, checker)...)
	return points
}

var (
	requestCounterPattern = `melody\.proxy\.requests\.layer\.([a-zA-Z]+)\.name\.(.*)\.complete\.(true|false)\.error\.(true|false)`
	requestCounterRegexp  = regexp.MustCompile(requestCounterPattern)

	responseCounterPattern = `melody\.router\.response\.(.*)\.stagauge.gotus\.([\d]{3})\.count`
	responseCounterRegexp  = regexp.MustCompile(responseCounterPattern)
)

func requestPoints(hostname string, now time.Time, counters map[string]int64, logger logging.Logger, checker alert.Checker) []*client.Point {
	res := []*client.Point{}
	mu.Lock()
	for k, count := range counters {
		if !requestCounterRegexp.MatchString(k) {
			continue
		}
		params := requestCounterRegexp.FindAllStringSubmatch(k, -1)[0][1:]
		tags := map[string]string{
			"host":     hostname,
			"layer":    params[0],
			"name":     params[1],
			"complete": params[2],
			"error":    params[3],
		}
		last, ok := lastRequestCount[strings.Join(params, ".")]
		if !ok {
			last = 0
		}
		fields := map[string]interface{}{
			"total": int(count),
			"count": int(count) - last,
		}
		lastRequestCount[strings.Join(params, ".")] = int(count)

		countersPoint, err := client.NewPoint("requests", tags, fields, now)
		if err != nil {
			logger.Error("creating request counters point:", err.Error())
			continue
		}
		res = append(res, countersPoint)

		if tags["layer"] == "endpoint" && count > 0 {
			err = checker(tags["name"], "error", count)
			if err != nil {
				logger.Error("creating request counters point:", err.Error())
				continue
			}
		}
	}
	mu.Unlock()
	return res
}

func responsePoints(hostname string, now time.Time, counters map[string]int64, logger logging.Logger, checker alert.Checker) []*client.Point {
	res := []*client.Point{}
	mu.Lock()
	for k, count := range counters {
		if !responseCounterRegexp.MatchString(k) {
			continue
		}
		params := responseCounterRegexp.FindAllStringSubmatch(k, -1)[0][1:]
		tags := map[string]string{
			"host":   hostname,
			"name":   params[0],
			"status": params[1],
		}
		last, ok := lastResponseCount[strings.Join(params, ".")]
		if !ok {
			last = 0
		}
		fields := map[string]interface{}{
			"total": int(count),
			"count": int(count) - last,
		}
		lastResponseCount[strings.Join(params, ".")] = int(count)

		countersPoint, err := client.NewPoint("responses", tags, fields, now)
		if err != nil {
			logger.Error("creating response counters point:", err.Error())
			continue
		}
		res = append(res, countersPoint)
	}
	mu.Unlock()
	return res
}

func connectionPoints(hostname string, now time.Time, counters map[string]int64, logger logging.Logger, checker alert.Checker) []*client.Point {
	res := make([]*client.Point, 2)

	in := map[string]interface{}{
		"current": int(counters["melody.router.connected"]),
		"total":   int(counters["melody.router.connected-total"]),
	}
	incoming, err := client.NewPoint("router", map[string]string{"host": hostname, "direction": "in"}, in, now)
	if err != nil {
		logger.Error("creating incoming connection counters point:", err.Error())
		return res
	}
	res[0] = incoming

	out := map[string]interface{}{
		"current": int(counters["melody.router.disconnected"]),
		"total":   int(counters["melody.router.disconnected-total"]),
	}
	outgoing, err := client.NewPoint("router", map[string]string{"host": hostname, "direction": "out"}, out, now)
	if err != nil {
		logger.Error("creating outgoing connection counters point:", err.Error())
		return res
	}
	res[1] = outgoing

	return res
}
