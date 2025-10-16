package repository

import (
	"context"

	"synapsis/inventory/app/model"
	"synapsis/inventory/database/connection"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type productRepository struct {
	db *gorm.DB
}

type ProductRepository interface {
	Find(ID int64, ctx context.Context) (*model.Product, error)
	FindTransactionProduct(ID int64, tx *gorm.DB, ctx context.Context) (*model.Product, error)
	Update(product *model.Product, tx *gorm.DB, ctx context.Context) error
	Create(product *model.Product, tx *gorm.DB, ctx context.Context) (*model.Product, error)
}

func NewProductRepository(instance connection.DBInstance) ProductRepository {
	return &productRepository{db: instance.Database()}
}

func (p *productRepository) Find(ID int64, ctx context.Context) (*model.Product, error) {
	resp := new(model.Product)
	tx := p.db.WithContext(ctx).
		Where("id = ?", ID).
		First(resp)

	return resp, tx.Error
}

func (p *productRepository) Update(product *model.Product, tx *gorm.DB, ctx context.Context) error {
	return tx.WithContext(ctx).Save(product).Error
}

func (p *productRepository) FindTransactionProduct(ID int64, tx *gorm.DB, ctx context.Context) (*model.Product, error) {
	resp := new(model.Product)
	tx = tx.WithContext(ctx).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("id = ?", ID).
		First(resp)

	return resp, tx.Error
}

func (p *productRepository) Create(product *model.Product, tx *gorm.DB, ctx context.Context) (*model.Product, error) {
	stmt := tx.WithContext(ctx).Create(product)
	return product, stmt.Error
}
