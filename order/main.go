package main

import (
	"synapsis/order/bootstrap"
	"synapsis/order/bootstrap/server"
	"synapsis/order/router"
)

func main() {
	c := bootstrap.NewBootstrap()
	err := c.Invoke(router.NewRoute)
	if err != nil {
		panic(err)
	}

	if err = c.Invoke(func(instance *server.EchoInstance) {
		instance.RunWithGracefullyShutdown()
	}); err != nil {
		panic(err)
	}
}
