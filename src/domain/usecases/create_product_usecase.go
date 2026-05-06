//go:generate mockgen -source create_product_usecase.go -destination mock/create_product_usecase_mock.go -package usecasesmock
package usecases

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/zthiagovalle/demo-golang/src/domain/enums"
	"github.com/zthiagovalle/demo-golang/src/domain/exceptions"
	"github.com/zthiagovalle/demo-golang/src/domain/models"
	"github.com/zthiagovalle/demo-golang/src/domain/producers"
	"github.com/zthiagovalle/demo-golang/src/domain/repositories"
)

type ICreateProductUsecase interface {
	Execute(ctx context.Context, input *models.ProductCreate) (*models.Product, error)
}

type CreateProductUsecase struct {
	ProductsRepository repositories.IProductsRepository
	ProductProducer    producers.IProductProducer
}

func NewCreateProductUsecase(repo repositories.IProductsRepository, producer producers.IProductProducer) *CreateProductUsecase {
	return &CreateProductUsecase{ProductsRepository: repo, ProductProducer: producer}
}

func (u *CreateProductUsecase) Execute(ctx context.Context, input *models.ProductCreate) (*models.Product, error) {
	if strings.TrimSpace(input.Name) == "" {
		return nil, errors.New(exceptions.ErrProductNameRequired)
	}
	if input.PriceCents <= 0 {
		return nil, errors.New(exceptions.ErrProductPriceInvalid)
	}

	now := time.Now().UTC()
	product := &models.Product{
		ID:          uuid.NewString(),
		Name:        input.Name,
		Description: input.Description,
		PriceCents:  input.PriceCents,
		Status:      enums.ProductStatusActive,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := u.ProductsRepository.Insert(ctx, product); err != nil {
		return nil, errors.New(exceptions.ErrOnInsertProduct)
	}

	if err := u.ProductProducer.PublishCreated(ctx, product); err != nil {
		return nil, errors.New(exceptions.ErrOnPublishProductEvent)
	}

	return product, nil
}
