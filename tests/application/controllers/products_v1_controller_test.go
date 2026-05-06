package controllers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/zthiagovalle/demo-golang/src/application/controllers"
	"github.com/zthiagovalle/demo-golang/src/core/shared"
	"github.com/zthiagovalle/demo-golang/src/domain/enums"
	"github.com/zthiagovalle/demo-golang/src/domain/models"
	usecasesmock "github.com/zthiagovalle/demo-golang/src/domain/usecases/mock"
)

func TestCreateProduct_ReturnsCreated(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	createUC := usecasesmock.NewMockICreateProductUsecase(ctrl)
	createUC.EXPECT().
		Execute(gomock.Any(), gomock.Any()).
		Return(&models.Product{
			ID:         "abc-123",
			Name:       "Cafe",
			PriceCents: 4990,
			Status:     enums.ProductStatusActive,
		}, nil)

	c := controllers.NewProductsV1Controller(
		usecasesmock.NewMockIGetAllPaginatedProductsUsecase(ctrl),
		createUC,
		usecasesmock.NewMockIUpdateProductUsecase(ctrl),
		usecasesmock.NewMockIDeleteProductUsecase(ctrl),
		usecasesmock.NewMockIToggleProductStatusUsecase(ctrl),
	)

	e := echo.New()
	e.Validator = shared.NewCustomValidator()

	body := `{"name":"Cafe","description":"Grao","price_cents":4990}`
	req := httptest.NewRequest(http.MethodPost, "/v1/products", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	ec := e.NewContext(req, rec)

	require.NoError(t, c.CreateProduct(ec))
	assert.Equal(t, http.StatusCreated, rec.Code)

	var got models.Product
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &got))
	assert.Equal(t, "abc-123", got.ID)
}
