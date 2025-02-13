package cmd

import (
	"ecommerce-payments/helpers"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (d *Dependency) MiddlewareValidateAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(e echo.Context) error {
		auth := e.Request().Header.Get("Authorization")
		if auth == "" {
			helpers.Logger.Errorf("authorization empty")
			return helpers.SendResponseHTTP(e, http.StatusUnauthorized, "unauthorized", nil)
		}

		profile, err := d.External.GetProfile(e.Request().Context(), auth)
		if err != nil {
			helpers.Logger.Errorf("failed to get profile: %s", err)
			return helpers.SendResponseHTTP(e, http.StatusUnauthorized, "unauthorized", nil)
		}

		e.Set("profile", profile)

		return next(e)
	}

}
