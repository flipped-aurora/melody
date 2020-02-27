package melody

import (
	"melody/logging"
)

// LoadPlugins 加载并注册插件
func LoadPlugins(folder, pattern string, logger logging.Logger) {

	// TODO load client plugin

	// load server plugin
	logger.Info("http handler plugins begin to load.")
	// n, err = server.Load(
	// 	folder,
	// 	pattern,
	// 	server.RegisterHandler,
	// )
	// if err != nil {
	// 	logger.Warning("loading plugins:", err)
	// }
	// logger.Info("total http handler plugins loaded:", n)
}
