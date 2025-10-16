package rpc

import (
	"context"
	"synapsis/inventory/app/service"
	pb "synapsis/proto-definitions/inventory"
	"sync"

	"github.com/jinzhu/copier"
)

type InventoryHandler struct {
	pb.UnimplementedInventoryServiceServer
	mu                  sync.Mutex
	nextID              int64
	checkStockService   service.CheckStockService
	reserveStockService service.ReserveStockService
	releaseStockService service.ReleaseStockService
}

func NewInventoryHandler(
	checkStockService service.CheckStockService,
	reserveStockService service.ReserveStockService,
	releaseStockService service.ReleaseStockService,
) *InventoryHandler {
	return &InventoryHandler{
		nextID:              1,
		checkStockService:   checkStockService,
		reserveStockService: reserveStockService,
		releaseStockService: releaseStockService,
	}
}

func (i *InventoryHandler) CheckStock(ctx context.Context, req *pb.CheckStockRequest) (*pb.CheckStockResponse, error) {
	i.mu.Lock()
	defer i.mu.Unlock()

	sReq := service.CheckStockRequest{}
	_ = copier.Copy(&sReq, req)

	resp, err := i.checkStockService.Execute(ctx, sReq)
	if err != nil {
		return nil, err
	}

	pProduct := pb.Product{}
	_ = copier.Copy(&pProduct, resp.Product)

	pStock := pb.Stock{}
	_ = copier.Copy(&pStock, resp.Stock)
	return &pb.CheckStockResponse{
		Product:     &pProduct,
		Stock:       &pStock,
		IsAvailable: resp.IsAvailable,
		Quantity:    req.Quantity,
	}, nil
}

func (i *InventoryHandler) ReserveStock(ctx context.Context, req *pb.ReserveStockRequest) (*pb.ReserveStockResponse, error) {
	i.mu.Lock()
	defer i.mu.Unlock()

	sReq := service.ReserveStockRequest{}
	_ = copier.Copy(&sReq, req)

	_, err := i.reserveStockService.Execute(ctx, sReq)
	if err != nil {
		return nil, err
	}

	return &pb.ReserveStockResponse{
		Success: true,
		Message: "success reserve stock",
	}, nil
}

func (i *InventoryHandler) ReleaseStock(ctx context.Context, req *pb.ReleaseStockRequest) (*pb.ReleaseStockResponse, error) {
	i.mu.Lock()
	defer i.mu.Unlock()

	sReq := service.ReleaseStockRequest{}
	_ = copier.Copy(&sReq, req)

	_, err := i.releaseStockService.Execute(sReq, ctx)
	if err != nil {
		return nil, err
	}

	return &pb.ReleaseStockResponse{
		Success: true,
		Message: "success release stock",
	}, nil
}
