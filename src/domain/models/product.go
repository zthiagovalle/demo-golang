package models

import (
	"time"

	"github.com/zthiagovalle/demo-golang/src/domain/enums"
)

type Product struct {
	ID          string              `json:"id"`
	Name        string              `json:"name"`
	Description string              `json:"description"`
	PriceCents  int64               `json:"price_cents"`
	Status      enums.ProductStatus `json:"status"`
	CreatedAt   time.Time           `json:"created_at"`
	UpdatedAt   time.Time           `json:"updated_at"`
}
