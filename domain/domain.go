package domain

import "gorm.io/gorm"

type Topic struct {
	gorm.Model
	Name     string  `json:"name"`
	Category string  `json:"category"`
	PriceMin float64 `json:"price_min"`
	PriceMax float64 `json:"price_max"`
}

type User struct {
	gorm.Model
	Email    string  `json:"email"`
	Password string  `json:"password"`
	Topics   []Topic `gorm:"many2many:user_topics;" json:"topics"`
}
