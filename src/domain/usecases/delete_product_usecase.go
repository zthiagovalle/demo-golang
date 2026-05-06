//go:generate mockgen -source delete_product_usecase.go -destination mock/delete_product_usecase_mock.go -package usecasesmock
package usecases

import (
	"context"
	"errors"

	"github.com/zthiagovalle/demo-golang/src/domain/exceptions"
	"github.com/zthiagovalle/demo-golang/src/domain/repositories"
)

type IDeleteProductUsecase interface {
	Execute(ctx context.Context, id string) error
}

type DeleteProductUsecase struct {
	ProductsRepository repositories.IProductsRepository
}

func NewDeleteProductUsecase(repo repositories.IProductsRepository) *DeleteProductUsecase {
	return &DeleteProductUsecase{ProductsRepository: repo}
}

func (u *DeleteProductUsecase) Execute(ctx context.Context, id string) error {
	existing, err := u.ProductsRepository.FindByID(ctx, id)
	if err != nil {
		return errors.New(exceptions.ErrOnFindProductByID)
	}
	if existing == nil {
		return errors.New(exceptions.ErrProductNotFound)
	}
	if err := u.ProductsRepository.Delete(ctx, id); err != nil {
		return errors.New(exceptions.ErrOnDeleteProduct)
	}
	return nil
}
