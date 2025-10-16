package bootstrap

import (
	grpcclient "synapsis/order/bootstrap/grpc-client"
	"synapsis/order/bootstrap/server"
	"synapsis/order/config"
	"synapsis/order/database/connection"

	"github.com/labstack/echo/v4"
	"go.uber.org/dig"
)

func NewApplication(container *dig.Container) *dig.Container {
	var err error

	// provide config
	if err = container.Provide(func() config.AppConfig {
		return config.NewAppConfig()
	}); err != nil {
		panic(err)
	}

	// provide postgres connection
	if err = container.Provide(func(cfg config.AppConfig) connection.DBInstance {
		return connection.NewDatabaseInstance(cfg)
	}); err != nil {
		panic(err)
	}

	if err = container.Provide(func(cfg config.AppConfig) grpcclient.GrpcClientInstance {
		return grpcclient.NewGrpcClient(cfg)
	}); err != nil {
		panic(err)
	}

	// provide echo instance
	if err = container.Provide(server.NewEchoInstance); err != nil {
		panic(err)
	}

	// provide router
	if err = container.Provide(func() *echo.Echo {
		e := echo.New()
		return e
	}); err != nil {
		panic(err)
	}

	return container
}
