package graphic

import (
	"github.com/labstack/echo"

	"net/http"
)

// Gets is get lesson graphic.
func Gets(c echo.Context) error {
	return c.JSON(http.StatusOK, "dummy")
}
