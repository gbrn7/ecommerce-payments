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
