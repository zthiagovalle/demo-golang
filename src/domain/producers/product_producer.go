//go:generate mockgen -source product_producer.go -destination mock/product_producer_mock.go -package producersmock
package producers

import (
	"context"

	"github.com/zthiagovalle/demo-golang/src/domain/models"
)

type IProductProducer interface {
	PublishCreated(ctx context.Context, p *models.Product) error
}
