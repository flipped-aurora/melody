package melody

import (
	"context"
	"fmt"
	"melody/config"
	"melody/logging"
	consul "melody/middleware/melody-consul"
	etcd "melody/middleware/melody-etcd"
	"melody/sd"
	"melody/sd/dnssrv"
)

// RegisterSubscriberFactories registers all the available sd adaptors
// sd = service discovery
func RegisterSubscriberFactories(ctx context.Context, cfg config.ServiceConfig, logger logging.Logger) func(n string, p int) {
	// setup the etcd client if necessary
	// etcd Raft
	etcdClient, err := etcd.New(ctx, cfg.ExtraConfig)
	if err != nil {
		logger.Warning("building the etcd client:", err.Error())
	}
	register := sd.GetRegister()
	register.Register("etcd", etcd.SubscriberFactory(ctx, etcdClient))
	// register.Get("etcd")(backend) 会得到此backend中以第一个host为prefix的key的所有value的一个slice
	// etcdClient.GetEntries(host[0])

	// register the dns service discovery
	// 同上 register.Get("dns")(backend)
	register.Register("dns", dnssrv.SubscriberFactory)

	return func(name string, port int) {
		if err := consul.Register(ctx, cfg.ExtraConfig, port, name, logger); err != nil {
			logger.Error(fmt.Sprintf("Couldn't register %s:%d in consul: %s", name, port, err.Error()))
		}

		// TODO: add the call to the etcd service register
	}
}
