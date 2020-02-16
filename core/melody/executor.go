package melody

import (
	"context"
	"io"
	"melody/cmd"
	"melody/config"
	"melody/logging"
	gelf "melody/middleware/melody-gelf"
	gologging "melody/middleware/melody-gologging"
	logstash "melody/middleware/melody-logstash"
	metrics "melody/middleware/melody-metrics/gin"
	"os"
)

//NewExecutor return an new executor
func NewExecutor(ctx context.Context) cmd.Executor {
	return func(cfg config.ServiceConfig) {
		// 1. 确定以及初始化 log有哪些输出
		var writers []io.Writer
		// 1.1 检察是否使用Gelf
		gelfWriter, err := gelf.NewWriter(cfg.ExtraConfig)
		if err == nil {
			writers = append(writers, GelfWriter{gelfWriter})
			gologging.UpdateFormatSelector(func(w io.Writer) string {
				switch w.(type) {
				case GelfWriter:
					return "%{message}"
				default:
					return gologging.DefaultPattern
				}
			})
		}
		// 2.初始化Logger

		// 2.1 是否启用logstash
		// Logstash 是开源的服务器端数据处理管道，能够同时从多个来源采集数据，转换数据，然后将数据发送到您最喜欢的“存储库”中。
		// 所以没有logstash就没有下面其他logger
		logger, enableLogstashError := logstash.NewLogger(cfg.ExtraConfig, writers...)

		if enableLogstashError != nil {
			// 2.2 是否使用gologging
			var enableGologgingError error
			logger, enableGologgingError = gologging.NewLogger(cfg.ExtraConfig, writers...)

			if enableGologgingError != nil {
				// 2.3 默认使用基础Log  Level:Debug, Output:stdout, Prefix: ""
				logger, err := logging.NewLogger("DEBUG", os.Stdout, "")
				if err != nil {
					return
				}
				logger.Error("unable to create gologging logger")
			} else {
				logger.Debug("use gologging as logger")
			}
		} else {
			logger.Debug("use logstash as logger")
		}

		logger.Info("Melody server listening on port:", cfg.Port, "🎁")

		//TODO 3.Start Reporter (暂时不做)

		//TODO 4.加载插件 (暂时不做)

		//TODO 5.注册etcd服务发现

		// 6.创建Metrics监控
		_ = metrics.New(ctx, cfg.ExtraConfig, logger)
		//TODO 7. 集成influxdd
		//TODO 8. 集成opencensus
		//TODO 9. 集成bloomFilter
		//TODO 10. 集成JWT，注册RejecterFactory
		//TODO 11. Set up melody Router
		select{}
	}
}

type GelfWriter struct {
	io.Writer
}
