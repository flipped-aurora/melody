package melody

import (
	"context"
	"io"
	"melody/cmd"
	"melody/config"
	"melody/logging"
	gologging "melody/middleware/melody-gologging"
	logstash "melody/middleware/melody-logstash"
	"os"
)

//NewExecutor return an new executor
func NewExecutor(ctx context.Context) cmd.Executor {
	return func(cfg config.ServiceConfig) {
		//TODO 1. 确定以及初始化 log有哪些输出
		var writers []io.Writer
		//TODO 1.1 检察是否使用Gelf作为输出

		//TODO 2.初始化Logger

		//TODO 2.1 是否启用logstash
		logger, enableLogstashError := logstash.NewLogger(cfg.ExtraConfig, writers...)

		if enableLogstashError != nil {
			//TODO 2.2 是否使用gologging
			var enableGologgingError error
			logger, enableGologgingError = gologging.NewLogger(cfg.ExtraConfig, writers...)

			if enableGologgingError != nil {
				//TODO 2.3 默认使用基础Log  Level:Debug, Output:stdout, Prefix: ""
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

		//TODO Start Reporter (目前还不知道这在干什么)

		//TODO 加载插件
		//TODO ...
	}
}
