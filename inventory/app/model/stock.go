package model

import "time"

type Stock struct {
	Id             int64     `gorm:"primary_key;auto_increment" json:"id"`
	ProductId      int64     `gorm:"not null" json:"product_id"`
	TotalStock     int64     `gorm:"not null" json:"total_stock"`
	AvailableStock int64     `gorm:"not null" json:"available_stock"`
	ReservedStock  int64     `gorm:"not null" json:"reserved_stock"`
	CreatedAt      time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
