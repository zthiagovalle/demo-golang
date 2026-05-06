//go:generate mockgen -source update_product_usecase.go -destination mock/update_product_usecase_mock.go -package usecasesmock
package usecases

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/zthiagovalle/demo-golang/src/domain/exceptions"
	"github.com/zthiagovalle/demo-golang/src/domain/models"
	"github.com/zthiagovalle/demo-golang/src/domain/repositories"
)

type IUpdateProductUsecase interface {
	Execute(ctx context.Context, id string, input *models.ProductUpdate) (*models.Product, error)
}

type UpdateProductUsecase struct {
	ProductsRepository repositories.IProductsRepository
}

func NewUpdateProductUsecase(repo repositories.IProductsRepository) *UpdateProductUsecase {
	return &UpdateProductUsecase{ProductsRepository: repo}
}

func (u *UpdateProductUsecase) Execute(ctx context.Context, id string, input *models.ProductUpdate) (*models.Product, error) {
	if strings.TrimSpace(input.Name) == "" {
		return nil, errors.New(exceptions.ErrProductNameRequired)
	}
	if input.PriceCents <= 0 {
		return nil, errors.New(exceptions.ErrProductPriceInvalid)
	}

	existing, err := u.ProductsRepository.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New(exceptions.ErrOnFindProductByID)
	}
	if existing == nil {
		return nil, errors.New(exceptions.ErrProductNotFound)
	}

	existing.Name = input.Name
	existing.Description = input.Description
	existing.PriceCents = input.PriceCents
	existing.UpdatedAt = time.Now().UTC()

	if err := u.ProductsRepository.Update(ctx, existing); err != nil {
		return nil, errors.New(exceptions.ErrOnUpdateProduct)
	}
	return existing, nil
}
