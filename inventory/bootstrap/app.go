package bootstrap

import (
	"synapsis/inventory/bootstrap/server"
	"synapsis/inventory/config"
	"synapsis/inventory/database/connection"

	"github.com/labstack/echo/v4"
	"go.uber.org/dig"
	"google.golang.org/grpc"
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

	if err = container.Provide(server.NewRpcInstance); err != nil {
		panic(err)
	}

	if err = container.Provide(func() *grpc.Server {
		return grpc.NewServer()
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
