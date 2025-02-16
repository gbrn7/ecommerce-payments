package repository

import (
	"context"
	"ecommerce-payments/internal/models"

	"gorm.io/gorm"
)

type PaymentRepo struct {
	DB *gorm.DB
}

func (r *PaymentRepo) InsertNewPaymentMethod(ctx context.Context, req *models.PaymentMethod) error {
	return r.DB.Create(req).Error
}

func (r *PaymentRepo) DeletePaymentMethod(ctx context.Context, userID uint64, sourceID int, sourceName string) error {
	return r.DB.Exec("DELETE FROM payment_methods WHERE source_id = ? AND source_name = ? AND user_id = ?", sourceID, sourceName, userID).Error
}

func (r *PaymentRepo) GetPaymentMethod(ctx context.Context, userID uint64, sourceName string) (models.PaymentMethod, error) {
	var (
		resp models.PaymentMethod
		err  error
	)

	err = r.DB.Where("user_id = ?", userID).Where("source_name = ?", sourceName).First(&resp).Error
	return resp, err
}

func (r *PaymentRepo) InsertNewPaymentTransaction(ctx context.Context, req *models.PaymentTransaction) error {
	return r.DB.Create(req).Error
}

func (r *PaymentRepo) InsertNewPaymentRefund(ctx context.Context, req *models.PaymentRefund) error {
	return r.DB.Create(req).Error
}

func (r *PaymentRepo) GetPaymentByOrderID(ctx context.Context, orderID int) (models.PaymentTransaction, error) {
	var (
		resp models.PaymentTransaction
		err  error
	)

	err = r.DB.Where("order_id = ?", orderID).First(&resp).Error
	return resp, err
}

func (r *PaymentRepo) GetPaymentMethodByID(ctx context.Context, paymentMethodID int) (models.PaymentMethod, error) {
	var (
		resp models.PaymentMethod
		err  error
	)

	err = r.DB.Where("id = ?", paymentMethodID).First(&resp).Error
	return resp, err
}
