//go:generate mockgen -source toggle_product_status_usecase.go -destination mock/toggle_product_status_usecase_mock.go -package usecasesmock
package usecases

import (
	"context"
	"errors"

	"github.com/zthiagovalle/demo-golang/src/domain/exceptions"
	"github.com/zthiagovalle/demo-golang/src/domain/models"
	"github.com/zthiagovalle/demo-golang/src/domain/repositories"
)

type IToggleProductStatusUsecase interface {
	Execute(ctx context.Context, id string) (*models.Product, error)
}

type ToggleProductStatusUsecase struct {
	ProductsRepository repositories.IProductsRepository
}

func NewToggleProductStatusUsecase(repo repositories.IProductsRepository) *ToggleProductStatusUsecase {
	return &ToggleProductStatusUsecase{ProductsRepository: repo}
}

func (u *ToggleProductStatusUsecase) Execute(ctx context.Context, id string) (*models.Product, error) {
	existing, err := u.ProductsRepository.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New(exceptions.ErrOnFindProductByID)
	}
	if existing == nil {
		return nil, errors.New(exceptions.ErrProductNotFound)
	}

	newStatus := existing.Status.Toggle()
	if err := u.ProductsRepository.UpdateStatus(ctx, id, string(newStatus)); err != nil {
		return nil, errors.New(exceptions.ErrOnUpdateProduct)
	}
	existing.Status = newStatus
	return existing, nil
}
