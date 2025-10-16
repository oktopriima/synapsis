package main

import (
	"log"
	"os"
	"os/signal"
	"synapsis/inventory/bootstrap"
	"synapsis/inventory/bootstrap/server"
	"synapsis/inventory/router"
	"sync"
	"syscall"
)

func main() {
	c := bootstrap.NewBootstrap()
	if err := c.Invoke(router.NewGrpcRouter); err != nil {
		panic(err)
	}

	if err := c.Invoke(router.NewApiRouter); err != nil {
		panic(err)
	}

	var wg sync.WaitGroup
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// rpc server
	if err := c.Invoke(func(srv *server.RpcInstance) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			srv.RunWithGracefulShutdown()
		}()
	}); err != nil {
		panic(err)
	}

	// http server
	if err := c.Invoke(func(srv *server.EchoInstance) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			srv.RunWithGracefullyShutdown()
		}()
	}); err != nil {
		panic(err)
	}

	<-stop
	log.Println("Shutting down gRPC and HTTP server...")

	wg.Wait()
	log.Println("Goodbye!")
}
