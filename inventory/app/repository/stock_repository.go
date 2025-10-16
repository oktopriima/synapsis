package repository

import (
	"context"
	"synapsis/inventory/app/model"
	"synapsis/inventory/database/connection"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type stockRepository struct {
	db *gorm.DB
}

type StockRepository interface {
	Find(ID int64, ctx context.Context) (*model.Stock, error)
	FindTransaction(ID int64, tx *gorm.DB, ctx context.Context) (*model.Stock, error)
	Update(stock *model.Stock, tx *gorm.DB, ctx context.Context) error
	FindByProduct(ProductId int64, ctx context.Context) (*model.Stock, error)
	FindByProductTransaction(ProductId int64, tx *gorm.DB, ctx context.Context) (*model.Stock, error)
	Create(stock *model.Stock, tx *gorm.DB, ctx context.Context) (*model.Stock, error)
}

func NewStockRepository(instance connection.DBInstance) StockRepository {
	return &stockRepository{
		db: instance.Database(),
	}
}

func (s *stockRepository) Find(ID int64, ctx context.Context) (*model.Stock, error) {
	output := new(model.Stock)
	tx := s.db.WithContext(ctx).Where("id = ?", ID).First(output)
	return output, tx.Error
}

func (s *stockRepository) Update(stock *model.Stock, tx *gorm.DB, ctx context.Context) error {
	return tx.WithContext(ctx).Save(stock).Error
}

func (s *stockRepository) FindTransaction(ID int64, tx *gorm.DB, ctx context.Context) (*model.Stock, error) {
	output := new(model.Stock)
	tx = tx.WithContext(ctx).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("id = ?", ID).
		First(output)
	return output, tx.Error
}

func (s *stockRepository) FindByProduct(ProductId int64, ctx context.Context) (*model.Stock, error) {
	output := new(model.Stock)
	tx := s.db.WithContext(ctx).
		Where("product_id = ?", ProductId).
		First(output)
	return output, tx.Error
}

func (s *stockRepository) FindByProductTransaction(ProductId int64, tx *gorm.DB, ctx context.Context) (*model.Stock, error) {
	output := new(model.Stock)
	tx = tx.WithContext(ctx).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("product_id = ?", ProductId).
		First(output)
	return output, tx.Error
}

func (s *stockRepository) Create(stock *model.Stock, tx *gorm.DB, ctx context.Context) (*model.Stock, error) {
	stmt := tx.WithContext(ctx).Create(stock)
	return stock, stmt.Error
}
