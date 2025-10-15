package bootstrap

import (
	"go.uber.org/dig"
	"synapsis/order/app/controller"
)

func NewController(container *dig.Container) *dig.Container {
	var err error

	if err = container.Provide(controller.NewCreateController); err != nil {
		panic(err)
	}

	return container
}
