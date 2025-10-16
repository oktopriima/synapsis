package bootstrap

import (
	"synapsis/inventory/app/repository"

	"go.uber.org/dig"
)

func NewRepository(container *dig.Container) *dig.Container {
	var err error

	if err = container.Provide(repository.NewProductRepository); err != nil {
		panic(err)
	}

	if err = container.Provide(repository.NewStockMovementRepository); err != nil {
		panic(err)
	}

	if err = container.Provide(repository.NewStockRepository); err != nil {
		panic(err)
	}

	return container
}
