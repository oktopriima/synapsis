package repository

import (
	"context"
	"synapsis/order/app/model"
	"synapsis/order/database/connection"

	"gorm.io/gorm"
)

type OrderRepository interface {
	Create(order *model.Order, tx *gorm.DB, ctx context.Context) (*model.Order, error)
}

type orderRepositoryImpl struct {
	db *gorm.DB
}

func NewOrderRepository(instance connection.DBInstance) OrderRepository {
	return &orderRepositoryImpl{db: instance.Database()}
}

func (o *orderRepositoryImpl) Create(order *model.Order, tx *gorm.DB, ctx context.Context) (*model.Order, error) {
	stmt := tx.WithContext(ctx).Create(order)
	return order, stmt.Error
}
