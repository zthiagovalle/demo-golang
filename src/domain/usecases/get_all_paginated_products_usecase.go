//go:generate mockgen -source get_all_paginated_products_usecase.go -destination mock/get_all_paginated_products_usecase_mock.go -package usecasesmock
package usecases

import (
	"context"
	"errors"

	"github.com/zthiagovalle/demo-golang/src/domain/exceptions"
	"github.com/zthiagovalle/demo-golang/src/domain/models"
	"github.com/zthiagovalle/demo-golang/src/domain/repositories"
)

type IGetAllPaginatedProductsUsecase interface {
	Execute(ctx context.Context, params *models.ProductPageParams) (*models.ProductPage, error)
}

type GetAllPaginatedProductsUsecase struct {
	ProductsRepository repositories.IProductsRepository
}

func NewGetAllPaginatedProductsUsecase(repo repositories.IProductsRepository) *GetAllPaginatedProductsUsecase {
	return &GetAllPaginatedProductsUsecase{ProductsRepository: repo}
}

func (u *GetAllPaginatedProductsUsecase) Execute(ctx context.Context, params *models.ProductPageParams) (*models.ProductPage, error) {
	if params.Page < 1 {
		params.Page = 1
	}
	if params.PageSize < 1 || params.PageSize > 100 {
		params.PageSize = 20
	}
	page, err := u.ProductsRepository.Paginate(ctx, params)
	if err != nil {
		return nil, errors.New(exceptions.ErrOnPaginateProducts)
	}
	return page, nil
}
