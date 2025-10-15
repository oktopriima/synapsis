package model

import "time"

type Order struct {
	ID          int64     `gorm:"primaryKey"`
	ProductID   int64     `gorm:"not null"`
	ProductName string    `gorm:"not null"`
	Quantity    int64     `gorm:"not null"`
	Price       float64   `gorm:"not null"`
	FinalPrice  float64   `gorm:"not null"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}
