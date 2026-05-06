package models

type ProductUpdate struct {
	Name        string `json:"name" validate:"required,min=1"`
	Description string `json:"description"`
	PriceCents  int64  `json:"price_cents" validate:"required,gt=0"`
}
