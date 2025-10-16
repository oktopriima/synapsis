package bootstrap

import (
	"synapsis/inventory/app/service"

	"go.uber.org/dig"
)

func NewService(container *dig.Container) *dig.Container {
	var err error

	if err = container.Provide(service.NewCheckStockService); err != nil {
		panic(err)
	}

	if err = container.Provide(service.NewReserveStockService); err != nil {
		panic(err)
	}

	if err = container.Provide(service.NewReleaseStockService); err != nil {
		panic(err)
	}

	if err = container.Provide(service.NewCreateProductService); err != nil {
		panic(err)
	}

	if err = container.Provide(service.NewAddStockService); err != nil {
		panic(err)
	}

	return container
}
