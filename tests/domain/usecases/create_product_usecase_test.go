package usecases_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/zthiagovalle/demo-golang/src/domain/enums"
	"github.com/zthiagovalle/demo-golang/src/domain/exceptions"
	"github.com/zthiagovalle/demo-golang/src/domain/models"
	"github.com/zthiagovalle/demo-golang/src/domain/usecases"
	producersmock "github.com/zthiagovalle/demo-golang/src/domain/producers/mock"
	repositoriesmock "github.com/zthiagovalle/demo-golang/src/domain/repositories/mock"
)

func TestCreateProductUsecase_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := repositoriesmock.NewMockIProductsRepository(ctrl)
	producer := producersmock.NewMockIProductProducer(ctrl)

	repo.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(nil)
	producer.EXPECT().PublishCreated(gomock.Any(), gomock.Any()).Return(nil)

	uc := usecases.NewCreateProductUsecase(repo, producer)
	product, err := uc.Execute(context.Background(), &models.ProductCreate{
		Name:        "Cafe",
		Description: "Grao",
		PriceCents:  4990,
	})

	require.NoError(t, err)
	assert.NotEmpty(t, product.ID)
	assert.Equal(t, enums.ProductStatusActive, product.Status)
}

func TestCreateProductUsecase_ValidationErrors(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := repositoriesmock.NewMockIProductsRepository(ctrl)
	producer := producersmock.NewMockIProductProducer(ctrl)

	uc := usecases.NewCreateProductUsecase(repo, producer)

	tests := []struct {
		name    string
		input   *models.ProductCreate
		wantErr string
	}{
		{"empty name", &models.ProductCreate{Name: " ", PriceCents: 100}, exceptions.ErrProductNameRequired},
		{"zero price", &models.ProductCreate{Name: "Cafe", PriceCents: 0}, exceptions.ErrProductPriceInvalid},
		{"negative price", &models.ProductCreate{Name: "Cafe", PriceCents: -10}, exceptions.ErrProductPriceInvalid},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := uc.Execute(context.Background(), tt.input)
			require.Error(t, err)
			assert.Equal(t, tt.wantErr, err.Error())
		})
	}
}
