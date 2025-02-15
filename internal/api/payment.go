package api

import (
	"ecommerce-payments/constants"
	"ecommerce-payments/external"
	"ecommerce-payments/helpers"
	"ecommerce-payments/internal/interfaces"
	"ecommerce-payments/internal/models"
	"net/http"

	"github.com/labstack/echo/v4"
)

type PaymentAPI struct {
	PaymentService interfaces.IPaymentService
}

func (api *PaymentAPI) PaymentMethodLink(e echo.Context) error {
	var (
		log = helpers.Logger
	)
	req := models.PaymentMethodLinkRequest{}
	if err := e.Bind(&req); err != nil {
		log.Error("error parse request: ", err)
		return helpers.SendResponseHTTP(e, http.StatusBadRequest, constants.ErrFailedBadRequest, nil)
	}

	if err := req.Validate(); err != nil {
		log.Error("failed to validate request", err)
		return helpers.SendResponseHTTP(e, http.StatusBadRequest, constants.ErrFailedBadRequest, nil)
	}

	err := api.PaymentService.PaymentMethodLink(e.Request().Context(), req)
	if err != nil {
		log.Error("failed to link payment method link ", err)
		return helpers.SendResponseHTTP(e, http.StatusInternalServerError, constants.ErrServerError, nil)
	}

	return helpers.SendResponseHTTP(e, http.StatusOK, constants.SuccessMessage, nil)
}

func (api *PaymentAPI) PaymentMethodOTP(e echo.Context) error {
	var (
		log = helpers.Logger
	)
	req := models.PaymentMethodOTPRequest{}
	if err := e.Bind(&req); err != nil {
		log.Error("error parse request: ", err)
		return helpers.SendResponseHTTP(e, http.StatusBadRequest, constants.ErrFailedBadRequest, nil)
	}

	if err := req.Validate(); err != nil {
		log.Error("failed to validate request", err)
		return helpers.SendResponseHTTP(e, http.StatusBadRequest, constants.ErrFailedBadRequest, nil)
	}

	profileCtx := e.Get("profile")
	profile, ok := profileCtx.(external.Profile)
	if !ok {
		log.Error("failed to get profile context, ")
		return helpers.SendResponseHTTP(e, http.StatusInternalServerError, constants.ErrServerError, nil)
	}

	err := api.PaymentService.PaymentMethodLinkConfirmation(e.Request().Context(), profile.Data.ID, req)
	if err != nil {
		log.Error("failed to link otp payment method link ", err)
		return helpers.SendResponseHTTP(e, http.StatusInternalServerError, constants.ErrServerError, nil)
	}

	return helpers.SendResponseHTTP(e, http.StatusOK, constants.SuccessMessage, nil)
}

func (api *PaymentAPI) PaymentMethodUnlink(e echo.Context) error {
	var (
		log = helpers.Logger
	)
	req := models.PaymentMethodLinkRequest{}
	if err := e.Bind(&req); err != nil {
		log.Error("error parse request: ", err)
		return helpers.SendResponseHTTP(e, http.StatusBadRequest, constants.ErrFailedBadRequest, nil)
	}

	if err := req.Validate(); err != nil {
		log.Error("failed to validate request", err)
		return helpers.SendResponseHTTP(e, http.StatusBadRequest, constants.ErrFailedBadRequest, nil)
	}

	profileCtx := e.Get("profile")
	profile, ok := profileCtx.(external.Profile)
	if !ok {
		log.Error("failed to get profile context, ")
		return helpers.SendResponseHTTP(e, http.StatusInternalServerError, constants.ErrServerError, nil)
	}

	err := api.PaymentService.PaymentMethodUnlink(e.Request().Context(), profile.Data.ID, req)
	if err != nil {
		log.Error("failed to unlink payment method link ", err)
		return helpers.SendResponseHTTP(e, http.StatusInternalServerError, constants.ErrServerError, nil)
	}

	return helpers.SendResponseHTTP(e, http.StatusOK, constants.SuccessMessage, nil)
}
