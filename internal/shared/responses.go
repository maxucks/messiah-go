package res

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

type errorResponse struct {
	Error string `json:"error"`
}

func ErrorResponse(msg string) errorResponse {
	return errorResponse{Error: msg}
}

func Error(c *echo.Context, code int, msg string) error {
	return c.JSON(code, errorResponse{
		Error: msg,
	})
}

func Ok(c *echo.Context, data any) error {
	return c.JSON(http.StatusOK, data)
}

func NoContent(c *echo.Context) error {
	return c.NoContent(http.StatusNoContent)
}

func Created(c *echo.Context) error {
	return c.NoContent(http.StatusCreated)
}

func InternalError(c *echo.Context, msg string) error {
	return Error(c, http.StatusInternalServerError, msg)
}
