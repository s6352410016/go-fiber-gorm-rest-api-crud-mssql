package models

import "time"

type Product struct {
	ProductId    uint    `gorm:"primaryKey"`
	ProductName  string  `json:"productName" gorm:"not null type:varchar(100)" validate:"required"`
	ProductPrice float64 `json:"productPrice" gorm:"not null type:float" validate:"required"`
	ProductImage string  `json:"productImage" gorm:"not null type:varchar(255)" validate:"required"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
