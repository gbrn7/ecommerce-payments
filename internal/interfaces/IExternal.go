package interfaces

import (
	"context"
	"ecommerce-payments/external"
	"ecommerce-payments/internal/models"
)

type IExternal interface {
	GetProfile(ctx context.Context, token string) (external.Profile, error)
	ProduceKafkaMessage(ctx context.Context, topic string, data []byte) error

	PaymentLink(ctx context.Context, req models.PaymentMethodLinkRequest) (external.PaymentLinkResponse, error)
	PaymentUnLink(ctx context.Context, req models.PaymentMethodLinkRequest) (external.PaymentLinkResponse, error)
	PaymentLinkConfirmation(ctx context.Context, walletID int, otp string) (external.PaymentLinkResponse, error)
}
