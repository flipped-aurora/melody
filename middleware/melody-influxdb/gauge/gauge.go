package gauge

import (
	"melody/logging"
	alert "melody/middleware/melody-alert"
	"time"

	"github.com/influxdata/influxdb/client/v2"
)

func Points(hostname string, now time.Time, counters map[string]int64, logger logging.Logger, checker alert.Checker) []*client.Point {
	res := make([]*client.Point, 4)

	in := map[string]interface{}{
		"gauge": int(counters["melody.router.connected-gauge"]),
	}
	incoming, err := client.NewPoint("router", map[string]string{"host": hostname, "direction": "in"}, in, now)
	if err != nil {
		logger.Error("creating incoming connection counters point:", err.Error())
		return res
	}
	res[0] = incoming

	out := map[string]interface{}{
		"gauge": int(counters["melody.router.disconnected-gauge"]),
	}
	outgoing, err := client.NewPoint("router", map[string]string{"host": hostname, "direction": "out"}, out, now)
	if err != nil {
		logger.Error("creating outgoing connection counters point:", err.Error())
		return res
	}
	res[1] = outgoing

	debug := map[string]interface{}{}
	runtime := map[string]interface{}{}

	for k, v := range counters {
		if k == "melody.router.connected-gauge" || k == "melody.router.disconnected-gauge" {
			continue
		}
		if k[:21] == "melody.service.debug." {
			debug[k[21:]] = int(v)
			continue
		}
		if k[:23] == "melody.service.runtime." {
			runtime[k[23:]] = int(v)
			continue
		}
		logger.Debug("unknown gauge key:", k)
	}

	debugPoint, err := client.NewPoint("debug", map[string]string{"host": hostname}, debug, now)
	if err != nil {
		logger.Error("creating debug counters point:", err.Error())
		return res
	}
	res[2] = debugPoint

	runtimePoint, err := client.NewPoint("runtime", map[string]string{"host": hostname}, runtime, now)
	if err != nil {
		logger.Error("creating runtime counters point:", err.Error())
		return res
	}
	res[3] = runtimePoint

	_ = checker("", "numgc", int64(debug["GCStats.NumGC"].(int)))
	_ = checker("", "sys", int64(runtime["MemStats.Sys"].(int)))
	_ = checker("", "heapsys", int64(runtime["MemStats.HeapSys"].(int)))
	_ = checker("", "stacksys", int64(runtime["MemStats.StackSys"].(int)))
	_ = checker("", "mcachesys", int64(runtime["MemStats.MCacheSys"].(int)))
	_ = checker("", "mspansys", int64(runtime["MemStats.MSpanSys"].(int)))

	return res
}
