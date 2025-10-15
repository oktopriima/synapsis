package model

import "time"

type StockMovement struct {
	Id          int64     `gorm:"primaryKey"`
	ProductId   int64     `gorm:"not null"`
	ChangeType  string    `gorm:"not null"`
	Quantity    int64     `gorm:"not null"`
	ReferenceId int64     `gorm:"not null"`
	Note        string    `gorm:"not null"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoCreateTime"`
}
