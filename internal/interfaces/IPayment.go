package interfaces

import (
	"context"
	"ecommerce-payments/internal/models"

	"github.com/labstack/echo/v4"
)

type IPaymentRepo interface {
	InsertNewPaymentMethod(ctx context.Context, req *models.PaymentMethod) error
	DeletePaymentMethod(ctx context.Context, userID uint64, sourceID int, sourceName string) error
	InsertNewPaymentTransaction(ctx context.Context, req *models.PaymentTransaction) error
	GetPaymentMethod(ctx context.Context, userID uint64, sourceName string) (models.PaymentMethod, error)
	GetPaymentByOrderID(ctx context.Context, orderID int) (models.PaymentTransaction, error)
	GetPaymentMethodByID(ctx context.Context, paymentMethodID int) (models.PaymentMethod, error)
	InsertNewPaymentRefund(ctx context.Context, req *models.PaymentRefund) error
}

type IPaymentService interface {
	PaymentMethodLink(ctx context.Context, req models.PaymentMethodLinkRequest) error
	PaymentMethodLinkConfirmation(ctx context.Context, userID uint64, req models.PaymentMethodOTPRequest) error
	PaymentMethodUnlink(ctx context.Context, userID uint64, req models.PaymentMethodLinkRequest) error
	InitiatePayment(ctx context.Context, req models.PaymentInitiatePayload) error
	RefundPayment(ctx context.Context, req models.RefundPayload) error
}

type IPaymentAPI interface {
	PaymentMethodLink(e echo.Context) error
	PaymentMethodOTP(e echo.Context) error
	PaymentMethodUnlink(e echo.Context) error
	InitiatePayment(kafkaPayload []byte) error
	RefundPayment(kafkaPayload []byte) error
}
