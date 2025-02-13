package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type HealthcheckAPI struct {
}

func (api *HealthcheckAPI) Healthcheck(e echo.Context) error {
	return e.JSON(http.StatusOK, nil)
}
