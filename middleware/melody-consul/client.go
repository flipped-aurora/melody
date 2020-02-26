package consul

import (
	"context"
	"github.com/go-contrib/uuid"
	"github.com/hashicorp/consul/api"
	"melody/config"
	"melody/logging"
)

func Register(ctx context.Context, e config.ExtraConfig, port int, serviceName string, logger logging.Logger) error {
	cfg, err := parse(e, port)
	if err != nil {
		return err
	}

	cfg.Name = serviceName

	return register(ctx, cfg, logger)
}

func register(ctx context.Context, cfg Config, logger logging.Logger) error {
	consulConfig := api.DefaultConfig()
	consulConfig.Address = cfg.Address
	c, err := api.NewClient(consulConfig)
	if err != nil {
		return err
	}

	service := &api.AgentServiceRegistration{
		Name: cfg.Name,
		Port: cfg.Port,
		Tags: cfg.Tags,
		ID:   uuid.NewV1().String(),
	}

	if err := c.Agent().ServiceRegister(service); err != nil {
		return err
	}

	go func() {
		<-ctx.Done()

		if err := c.Agent().ServiceDeregister(service.ID); err != nil {
			logger.Info("error when trying to deregister service", service.ID, ":", err)
		}
	}()

	return nil
}
