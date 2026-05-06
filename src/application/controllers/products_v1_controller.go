package controllers

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v5"

	"github.com/zthiagovalle/demo-golang/src/core/shared"
	"github.com/zthiagovalle/demo-golang/src/domain/exceptions"
	"github.com/zthiagovalle/demo-golang/src/domain/models"
	"github.com/zthiagovalle/demo-golang/src/domain/usecases"
)

type ProductsV1Controller struct {
	GetAllPaginatedProductsUsecase usecases.IGetAllPaginatedProductsUsecase
	CreateProductUsecase           usecases.ICreateProductUsecase
	UpdateProductUsecase           usecases.IUpdateProductUsecase
	DeleteProductUsecase           usecases.IDeleteProductUsecase
	ToggleProductStatusUsecase     usecases.IToggleProductStatusUsecase
}

func NewProductsV1Controller(
	getAllPaginated usecases.IGetAllPaginatedProductsUsecase,
	create usecases.ICreateProductUsecase,
	update usecases.IUpdateProductUsecase,
	del usecases.IDeleteProductUsecase,
	toggle usecases.IToggleProductStatusUsecase,
) *ProductsV1Controller {
	return &ProductsV1Controller{
		GetAllPaginatedProductsUsecase: getAllPaginated,
		CreateProductUsecase:           create,
		UpdateProductUsecase:           update,
		DeleteProductUsecase:           del,
		ToggleProductStatusUsecase:     toggle,
	}
}

func (c *ProductsV1Controller) Routes() []shared.Route {
	const basePath = "/v1/products"
	const basePathWithId = basePath + "/:id"
	const basePathWithIdAndToggleStatus = basePathWithId + "/toggle-status"

	return []shared.Route{
		{URI: basePath, Method: http.MethodGet, Handler: c.GetAllPaginatedProducts},
		{URI: basePath, Method: http.MethodPost, Handler: c.CreateProduct},
		{URI: basePathWithId, Method: http.MethodPut, Handler: c.UpdateProduct},
		{URI: basePathWithId, Method: http.MethodDelete, Handler: c.DeleteProduct},
		{URI: basePathWithIdAndToggleStatus, Method: http.MethodPatch, Handler: c.ToggleProductStatus},
	}
}

// @Summary Retornar lista paginada de produtos
// @Tags produtos
// @Accept json
// @Produce json
// @Success 200 {object} models.ProductPage
// @Failure 400
// @Failure 500
// @Param page query uint16 true "página" minimum(1)
// @Param pageSize query uint16 true "tamanho da página" minimum(1)
// @Param status query enums.ProductStatus false "status do produto"
// @Router /v1/products [get]
func (c *ProductsV1Controller) GetAllPaginatedProducts(ec *echo.Context) error {
	var params models.ProductPageParams
	if err := ec.Bind(&params); err != nil {
		return shared.ErrorJSON(ec, http.StatusBadRequest, exceptions.ErrInvalidBody)
	}
	if err := ec.Validate(&params); err != nil {
		return shared.ErrorJSON(ec, http.StatusBadRequest, exceptions.ErrValidation)
	}

	result, err := c.GetAllPaginatedProductsUsecase.Execute(ec.Request().Context(), &params)
	if err != nil {
		return shared.InternalErrorJSON(ec, err)
	}
	return ec.JSON(http.StatusOK, result)
}

// @Summary Criar produto
// @Tags produtos
// @Accept json
// @Produce json
// @Success 201 {object} models.Product
// @Failure 400
// @Failure 500
// @Param product body models.ProductCreate true "Dados do produto"
// @Router /v1/products [post]
func (c *ProductsV1Controller) CreateProduct(ec *echo.Context) error {
	var input models.ProductCreate
	if err := ec.Bind(&input); err != nil {
		return shared.ErrorJSON(ec, http.StatusBadRequest, exceptions.ErrInvalidBody)
	}
	if err := ec.Validate(&input); err != nil {
		return shared.ErrorJSON(ec, http.StatusBadRequest, exceptions.ErrValidation)
	}

	result, err := c.CreateProductUsecase.Execute(ec.Request().Context(), &input)
	if err != nil {
		switch err.Error() {
		case exceptions.ErrProductNameRequired, exceptions.ErrProductPriceInvalid:
			return shared.ErrorJSON(ec, http.StatusBadRequest, err.Error())
		}
		return shared.InternalErrorJSON(ec, err)
	}
	return ec.JSON(http.StatusCreated, result)
}

// @Summary Atualizar produto
// @Tags produtos
// @Accept json
// @Produce json
// @Success 200 {object} models.Product
// @Failure 400
// @Failure 404
// @Failure 500
// @Param id path string true "ID do produto"
// @Param product body models.ProductUpdate true "Dados do produto"
// @Router /v1/products/{id} [put]
func (c *ProductsV1Controller) UpdateProduct(ec *echo.Context) error {
	id, err := uuid.Parse(ec.Param("id"))
	if err != nil {
		return shared.ErrorJSON(ec, http.StatusBadRequest, exceptions.ErrInvalidUuid)
	}

	var input models.ProductUpdate
	if err := ec.Bind(&input); err != nil {
		return shared.ErrorJSON(ec, http.StatusBadRequest, exceptions.ErrInvalidBody)
	}
	if err := ec.Validate(&input); err != nil {
		return shared.ErrorJSON(ec, http.StatusBadRequest, exceptions.ErrValidation)
	}

	result, err := c.UpdateProductUsecase.Execute(ec.Request().Context(), id.String(), &input)
	if err != nil {
		return mapProductError(ec, err)
	}
	return ec.JSON(http.StatusOK, result)
}

// @Summary Deletar produto
// @Tags produtos
// @Accept json
// @Produce json
// @Success 204
// @Failure 400
// @Failure 404
// @Failure 500
// @Param id path string true "ID do produto"
// @Router /v1/products/{id} [delete]
func (c *ProductsV1Controller) DeleteProduct(ec *echo.Context) error {
	id, err := uuid.Parse(ec.Param("id"))
	if err != nil {
		return shared.ErrorJSON(ec, http.StatusBadRequest, exceptions.ErrInvalidUuid)
	}

	if err := c.DeleteProductUsecase.Execute(ec.Request().Context(), id.String()); err != nil {
		return mapProductError(ec, err)
	}
	return ec.NoContent(http.StatusNoContent)
}

// @Summary Alternar status do produto
// @Tags produtos
// @Accept json
// @Produce json
// @Success 200 {object} models.Product
// @Failure 400
// @Failure 404
// @Failure 500
// @Param id path string true "ID do produto"
// @Router /v1/products/{id}/toggle-status [patch]
func (c *ProductsV1Controller) ToggleProductStatus(ec *echo.Context) error {
	id, err := uuid.Parse(ec.Param("id"))
	if err != nil {
		return shared.ErrorJSON(ec, http.StatusBadRequest, exceptions.ErrInvalidUuid)
	}

	result, err := c.ToggleProductStatusUsecase.Execute(ec.Request().Context(), id.String())
	if err != nil {
		return mapProductError(ec, err)
	}
	return ec.JSON(http.StatusOK, result)
}

func mapProductError(ec *echo.Context, err error) error {
	if errors.Is(err, nil) {
		return nil
	}
	switch err.Error() {
	case exceptions.ErrProductNotFound:
		return shared.ErrorJSON(ec, http.StatusNotFound, err.Error())
	case exceptions.ErrProductNameRequired, exceptions.ErrProductPriceInvalid:
		return shared.ErrorJSON(ec, http.StatusBadRequest, err.Error())
	}
	return shared.InternalErrorJSON(ec, err)
}
