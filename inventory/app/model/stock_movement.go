package model

import "time"

type StockMovement struct {
	Id          int64     `gorm:"primaryKey" json:"id"`
	ProductId   int64     `gorm:"not null" json:"product_id"`
	ChangeType  string    `gorm:"not null" json:"change_type"`
	Quantity    int64     `gorm:"not null" json:"quantity"`
	ReferenceId int64     `gorm:"not null" json:"reference_id"`
	Note        string    `gorm:"not null" json:"note"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoCreateTime" json:"updated_at"`
}
