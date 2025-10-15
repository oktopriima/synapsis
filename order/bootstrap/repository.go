package bootstrap

import (
	"go.uber.org/dig"
	"synapsis/order/app/repository"
)

func NewRepository(container *dig.Container) *dig.Container {
	var err error

	if err = container.Provide(repository.NewOrderRepository); err != nil {
		panic(err)
	}

	return container
}
