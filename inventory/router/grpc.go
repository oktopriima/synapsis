package router

import (
	"synapsis/inventory/app/handler/rpc"
	pb "synapsis/proto-definitions/inventory"

	"google.golang.org/grpc"
)

func NewGrpcRouter(server *grpc.Server,
	inventoryHandler *rpc.InventoryHandler,
) {
	pb.RegisterInventoryServiceServer(server, inventoryHandler)
}
