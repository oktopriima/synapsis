package server

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"synapsis/inventory/config"

	"google.golang.org/grpc"
)

type (
	RpcInstance struct {
		Rpc    *grpc.Server
		Config config.AppConfig
	}
)

func NewRpcInstance(cfg config.AppConfig, server *grpc.Server) *RpcInstance {
	return &RpcInstance{
		Rpc:    server,
		Config: cfg,
	}
}

func (r *RpcInstance) RunRpcServer() error {
	port := fmt.Sprintf(":%s", r.Config.App.Port)
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
		return err
	}

	log.Printf("🚀 gRPC server is listening on %s\n", port)

	if err = r.Rpc.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
		return err
	}

	return nil
}

func (r *RpcInstance) RunWithGracefulShutdown() {
	go func() {
		if err := r.RunRpcServer(); err != nil {
			fmt.Println("RPC server closed")
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	log.Println("🛑 Shutting down gRPC server gracefully...")
	r.Rpc.GracefulStop()
}
