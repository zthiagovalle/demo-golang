package enums

import (
	"slices"

	"github.com/go-playground/validator/v10"
)

// @description Status do produto
type ProductStatus string // @name StatusProduto

const (
	ProductStatusActive   ProductStatus = "ACTIVE"   // Ativo
	ProductStatusInactive ProductStatus = "INACTIVE" // Inativo
)

var productStatusValues = []string{
	ProductStatusActive.String(),
	ProductStatusInactive.String(),
}

func (s ProductStatus) String() string {
	return string(s)
}

func (s ProductStatus) Toggle() ProductStatus {
	if s == ProductStatusActive {
		return ProductStatusInactive
	}
	return ProductStatusActive
}

func ProductStatusValidator(fl validator.FieldLevel) bool {
	return slices.Contains(productStatusValues, fl.Field().String())
}
