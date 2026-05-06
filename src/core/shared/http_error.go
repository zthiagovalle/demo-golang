package shared

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func ErrorJSON(c *echo.Context, status int, code string) error {
	return c.JSON(status, ErrorResponse{Code: code, Message: code})
}

func InternalErrorJSON(c *echo.Context, err error) error {
	return c.JSON(http.StatusInternalServerError, ErrorResponse{
		Code:    "errInternal",
		Message: err.Error(),
	})
}
