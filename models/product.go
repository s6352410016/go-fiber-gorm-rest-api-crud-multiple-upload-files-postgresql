package models

import (
	"time"

	"github.com/lib/pq"
)

type Product struct {
	ID            uint           `gorm:"primaryKey"`
	ProductName   string         `form:"productName" gorm:"not null"`
	ProductPrice  float64        `form:"productPrice" gorm:"not null"`
	ProductImages pq.StringArray `gorm:"type:text[]; not null"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
