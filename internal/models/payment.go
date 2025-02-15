package models

import (
	"time"

	"github.com/go-playground/validator/v10"
)

type PaymentInitiatePayload struct {
	UserID     uint64  `json:"user_id"`
	OrderID    int     `json:"order_id"`
	TotalPrice float64 `json:"total_price"`
}

type RefundPayload struct {
	OrderID int    `json:"order_id"`
	AdminID uint64 `json:"admin_id"`
}

type PaymentTransaction struct {
	ID               int       `json:"id"`
	OrderID          int       `json:"order_id" validate:"required"`
	PaymentMethodID  int       `json:"payment_method_id" validate:"required"`
	TotalPrice       float64   `json:"total_price" gorm:"column:total_price;type:decimal(10,2)"`
	Status           string    `json:"status" gorm:"column:status;type:varchar(10)"`
	PaymentReference string    `json:"payment_reference" gorm:"column:payment_reference;type:varchar(100)"`
	CreatedAt        time.Time `json:"-"`
	UpdatedAt        time.Time `json:"-"`
}

func (*PaymentTransaction) TableName() string {
	return "payment_transactions"
}

func (l PaymentTransaction) Validate() error {
	v := validator.New()
	return v.Struct(l)
}

type PaymentMethod struct {
	ID         int       `json:"id"`
	UserID     uint64    `json:"user_id"`
	SourceID   int       `json:"source_id"`
	SourceName string    `json:"source_name" gorm:"column:source_name;type:varchar(50)"`
	CreatedAt  time.Time `json:"-"`
	UpdatedAt  time.Time `json:"-"`
}

func (*PaymentMethod) TableName() string {
	return "payment_methods"
}

func (l PaymentMethod) Validate() error {
	v := validator.New()
	return v.Struct(l)
}
