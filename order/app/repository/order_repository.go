package repository

import (
	"context"
	"synapsis/order/app/model"
	"synapsis/order/database/connection"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type OrderRepository interface {
	Create(order *model.Order, tx *gorm.DB, ctx context.Context) (*model.Order, error)
	Find(Id int64, ctx context.Context) (*model.Order, error)
	FindTransaction(Id int64, tx *gorm.DB, ctx context.Context) (*model.Order, error)
	Update(order *model.Order, tx *gorm.DB, ctx context.Context) error
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

func (o *orderRepositoryImpl) FindTransaction(Id int64, tx *gorm.DB, ctx context.Context) (*model.Order, error) {
	output := new(model.Order)
	stmt := tx.WithContext(ctx).
		Clauses(clause.Locking{
			Strength: "UPDATE",
		}).
		Where("id = ?", Id).
		First(output)

	return output, stmt.Error
}

func (o *orderRepositoryImpl) Find(Id int64, ctx context.Context) (*model.Order, error) {
	output := new(model.Order)
	stmt := o.db.WithContext(ctx).
		Where("id = ?", Id).
		First(output)

	return output, stmt.Error
}

func (o *orderRepositoryImpl) Update(order *model.Order, tx *gorm.DB, ctx context.Context) error {
	return tx.WithContext(ctx).Save(order).Error
}
