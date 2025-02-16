package cmd

import (
	"ecommerce-payments/external"
	"ecommerce-payments/helpers"
	"ecommerce-payments/internal/api"
	"ecommerce-payments/internal/interfaces"
	"ecommerce-payments/internal/repository"
	"ecommerce-payments/internal/services"

	"github.com/labstack/echo/v4"
)

func ServeHTTP() {
	d := dependencyInject()

	e := echo.New()
	e.GET("/healthcheck", d.HealthcheckAPI.Healthcheck)

	paymentV1 := e.Group("/payment/v1")
	paymentV1.POST("/link", d.PaymentAPI.PaymentMethodLink, d.MiddlewareValidateAuth)
	paymentV1.POST("/link/confirm", d.PaymentAPI.PaymentMethodOTP, d.MiddlewareValidateAuth)
	paymentV1.DELETE("/unlink", d.PaymentAPI.PaymentMethodUnlink, d.MiddlewareValidateAuth)

	e.Start(":" + helpers.GetEnv("PORT", "9000"))
}

type Dependency struct {
	External       interfaces.IExternal
	HealthcheckAPI *api.HealthcheckAPI

	PaymentAPI interfaces.IPaymentAPI
}

func dependencyInject() Dependency {

	external := &external.External{}

	paymentRepo := &repository.PaymentRepo{
		DB: helpers.DB,
	}

	paymentSvc := &services.PaymentService{
		PaymentRepo: paymentRepo,
		External:    external,
	}
	paymentAPI := &api.PaymentAPI{
		PaymentService: paymentSvc,
	}

	return Dependency{
		External:       external,
		HealthcheckAPI: &api.HealthcheckAPI{},
		PaymentAPI:     paymentAPI,
	}
}
