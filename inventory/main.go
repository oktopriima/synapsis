package main

import (
	"synapsis/inventory/bootstrap"
	"synapsis/inventory/bootstrap/server"
	"synapsis/inventory/router"
)

func main() {
	c := bootstrap.NewBootstrap()
	if err := c.Invoke(router.NewGrpcRouter); err != nil {
		panic(err)
	}

	if err := c.Invoke(func(srv *server.RpcInstance) {
		srv.RunWithGracefulShutdown()
	}); err != nil {
		panic(err)
	}
}
