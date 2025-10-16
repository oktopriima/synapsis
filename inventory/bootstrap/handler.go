package bootstrap

import (
	"synapsis/inventory/app/handler/http"
	"synapsis/inventory/app/handler/rpc"

	"go.uber.org/dig"
)

func NewHandler(container *dig.Container) *dig.Container {
	var err error

	if err = container.Provide(rpc.NewInventoryHandler); err != nil {
		panic(err)
	}

	if err = container.Provide(http.NewProductHandler); err != nil {
		panic(err)
	}

	return container
}
