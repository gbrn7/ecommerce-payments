package models

import "github.com/go-playground/validator/v10"

type PaymentMethodLinkRequest struct {
	SourceID int `json:"source_id" validate:"required"`
}

func (l PaymentMethodLinkRequest) Validate() error {
	v := validator.New()
	return v.Struct(l)
}

type PaymentMethodOTPRequest struct {
	OTP      string `json:"otp" validate:"required"`
	SourceID int    `json:"source_id" validate:"required"`
}

func (l PaymentMethodOTPRequest) Validate() error {
	v := validator.New()
	return v.Struct(l)
}
