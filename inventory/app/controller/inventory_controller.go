package controller

import (
	"context"
	"fmt"
	"synapsis/inventory/app/model"
	"synapsis/inventory/app/repository"
	"synapsis/inventory/database/connection"
	pb "synapsis/proto-definitions/inventory"
	"sync"

	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

type InventoryController struct {
	pb.UnimplementedInventoryServiceServer
	mu                      sync.Mutex
	nextID                  int64
	db                      *gorm.DB
	productRepository       repository.ProductRepository
	stockMovementRepository repository.StockMovementRepository
}

func NewInventoryController(
	instance connection.DBInstance,
	productRepository repository.ProductRepository,
	stockMovementRepository repository.StockMovementRepository,
) *InventoryController {
	return &InventoryController{
		nextID:                  1,
		db:                      instance.Database(),
		productRepository:       productRepository,
		stockMovementRepository: stockMovementRepository,
	}
}

func (i *InventoryController) CheckStock(ctx context.Context, req *pb.CheckStockRequest) (*pb.CheckStockResponse, error) {
	i.mu.Lock()
	defer i.mu.Unlock()

	product, err := i.productRepository.Find(req.GetProductId(), ctx)
	if err != nil {
		return nil, err
	}

	return &pb.CheckStockResponse{
		ProductId:   product.ID,
		Quantity:    product.AvailableStock,
		IsAvailable: product.AvailableStock < req.Quantity,
	}, nil
}

func (i *InventoryController) ReserveStock(ctx context.Context, req *pb.ReserveStockRequest) (*pb.ReserveStockResponse, error) {
	i.mu.Lock()
	tx := i.db.Begin()

	defer func() {
		i.mu.Unlock()
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	product, err := i.productRepository.FindTransactionProduct(req.GetProductId(), tx, ctx)
	if err != nil {
		return nil, err
	}

	// check stock once again
	if product.AvailableStock < req.Quantity {
		return nil, fmt.Errorf("product has been sold")
	}

	product.AvailableStock -= req.GetQuantity()
	product.ReservedStock += req.GetQuantity()

	if err = i.productRepository.Update(product, tx, ctx); err != nil {
		return nil, err
	}

	// insert stock movement
	if err = i.stockMovementRepository.Create(&model.StockMovement{
		Id:          i.nextID,
		ProductId:   product.ID,
		ChangeType:  "RESERVED",
		Quantity:    req.GetQuantity(),
		ReferenceId: req.GetOrderId(),
		Note:        fmt.Sprintf("reserved %d for %d", req.GetQuantity(), req.GetOrderId()),
	}, tx, ctx); err != nil {
		return nil, err
	}

	// commit transaction
	tx.Commit()

	pbProduct := &pb.Product{}
	_ = copier.Copy(pbProduct, product)
	return &pb.ReserveStockResponse{
		Success: true,
		Message: "OK",
		Product: pbProduct,
	}, nil
}
