package bootstrap

import (
	"go.uber.org/dig"
	"synapsis/inventory/app/controller"
)

func NewController(container *dig.Container) *dig.Container {
	var err error

	if err = container.Provide(controller.NewInventoryController); err != nil {
		panic(err)
	}
	return container
}
