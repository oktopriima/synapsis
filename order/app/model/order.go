package model

import "time"

type Order struct {
	ID          int64     `gorm:"primaryKey" json:"id"`
	ProductID   int64     `gorm:"not null" json:"product_id"`
	ProductName string    `gorm:"not null" json:"product_name"`
	Quantity    int64     `gorm:"not null" json:"quantity"`
	Price       float64   `gorm:"not null" json:"price"`
	FinalPrice  float64   `gorm:"not null" json:"final_price"`
	Status      string    `gorm:"not null" json:"status"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
