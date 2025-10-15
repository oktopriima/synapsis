package router

import (
	"synapsis/inventory/app/controller"
	pb "synapsis/proto-definitions/inventory"

	"google.golang.org/grpc"
)

func NewGrpcRouter(server *grpc.Server,
	inventoryController *controller.InventoryController,
) {
	pb.RegisterInventoryServiceServer(server, inventoryController)
}
