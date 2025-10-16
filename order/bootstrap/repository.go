package bootstrap

import (
	"synapsis/order/app/repository"

	"go.uber.org/dig"
)

func NewRepository(container *dig.Container) *dig.Container {
	var err error

	if err = container.Provide(repository.NewOrderRepository); err != nil {
		panic(err)
	}

	return container
}
