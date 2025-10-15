package model

import "gorm.io/gorm"

type DB struct {
	*gorm.DB
}

func (db *DB) AutoMigrate() error {
	return db.DB.AutoMigrate(
		&Product{},
		&StockMovement{},
	)
}
