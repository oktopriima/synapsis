package bootstrap

import (
	"synapsis/order/app/handler/http"

	"go.uber.org/dig"
)

func NewController(container *dig.Container) *dig.Container {
	var err error

	if err = container.Provide(http.NewOrderHandler); err != nil {
		panic(err)
	}

	return container
}
