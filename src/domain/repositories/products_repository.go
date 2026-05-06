//go:generate mockgen -source products_repository.go -destination mock/products_repository_mock.go -package repositoriesmock
package repositories

import (
	"context"

	"github.com/zthiagovalle/demo-golang/src/domain/models"
)

type IProductsRepository interface {
	Insert(ctx context.Context, p *models.Product) error
	Update(ctx context.Context, p *models.Product) error
	Delete(ctx context.Context, id string) error
	FindByID(ctx context.Context, id string) (*models.Product, error)
	UpdateStatus(ctx context.Context, id, status string) error
	Paginate(ctx context.Context, params *models.ProductPageParams) (*models.ProductPage, error)
}
