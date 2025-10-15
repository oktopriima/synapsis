package grpc_client

import (
	"fmt"
	"synapsis/order/config"
	pb "synapsis/proto-definitions/inventory"

	"google.golang.org/grpc"
)

type GrpcClient struct {
	inventoryConn pb.InventoryServiceClient
}

type GrpcClientInstance interface {
	InventoryConnection() pb.InventoryServiceClient
}

func NewGrpcClient(config config.AppConfig) GrpcClientInstance {
	cl := new(GrpcClient)

	inventoryAddr := fmt.Sprintf("%s:%s", config.Rpc.Inventory.Address, config.Rpc.Inventory.Port)
	inventoryConn, err := grpc.Dial(inventoryAddr, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	cl.inventoryConn = pb.NewInventoryServiceClient(inventoryConn)

	return cl
}

func (i *GrpcClient) InventoryConnection() pb.InventoryServiceClient {
	return i.inventoryConn
}
