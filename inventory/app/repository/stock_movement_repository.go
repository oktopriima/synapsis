package repository

import (
	"context"
	"synapsis/inventory/app/model"
	"synapsis/inventory/database/connection"

	"gorm.io/gorm"
)

type stockMovementRepository struct {
	db *gorm.DB
}

type StockMovementRepository interface {
	Create(movement *model.StockMovement, tx *gorm.DB, ctx context.Context) error
}

func NewStockMovementRepository(instance connection.DBInstance) StockMovementRepository {
	return &stockMovementRepository{db: instance.Database()}
}

func (s *stockMovementRepository) Create(movement *model.StockMovement, tx *gorm.DB, ctx context.Context) error {
	return tx.WithContext(ctx).Create(movement).Error
}
