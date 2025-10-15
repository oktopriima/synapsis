package model

import "time"

type Product struct {
	ID             int64     `gorm:"primaryKey"`
	Name           string    `gorm:"not null"`
	Sku            string    `gorm:"not null,unique"`
	Description    string    `gorm:"not null"`
	Price          float64   `gorm:"not null"`
	TotalStock     int64     `gorm:"not null"`
	AvailableStock int64     `gorm:"not null"`
	ReservedStock  int64     `gorm:"not null"`
	CreatedAt      time.Time `gorm:"autoCreateTime"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime"`
}

func (p *Product) TableName() string {
	return "products"
}
