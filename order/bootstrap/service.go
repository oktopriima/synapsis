package bootstrap

import (
	"synapsis/order/app/service"

	"go.uber.org/dig"
)

func NewService(container *dig.Container) *dig.Container {
	var err error

	if err = container.Provide(service.NewCreateOrderService); err != nil {
		panic(err)
	}

	if err = container.Provide(service.NewCancelOrderService); err != nil {
		panic(err)
	}

	return container
}
