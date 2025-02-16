package services

import (
	"context"
	"ecommerce-payments/external"
	"ecommerce-payments/helpers"
	"ecommerce-payments/internal/interfaces"
	"ecommerce-payments/internal/models"
	"encoding/json"

	"github.com/pkg/errors"
)

type PaymentService struct {
	PaymentRepo interfaces.IPaymentRepo
	External    interfaces.IExternal
}

func (s *PaymentService) PaymentMethodLink(ctx context.Context, req models.PaymentMethodLinkRequest) error {
	resp, err := s.External.PaymentLink(ctx, req)
	if err != nil {
		return errors.Wrap(err, "failed to request payment link confirm to e-wallet")
	}

	helpers.Logger.WithField("OTP", resp.Data.OTP).Info("Link response is success, need otp confirm")
	return nil
}

func (s *PaymentService) PaymentMethodLinkConfirmation(ctx context.Context, userID uint64, req models.PaymentMethodOTPRequest) error {
	_, err := s.External.PaymentLinkConfirmation(ctx, req.SourceID, req.OTP)
	if err != nil {
		return errors.Wrap(err, "failed to request payment link confirm to e-wallet")
	}

	paymentMethod := models.PaymentMethod{
		UserID:     userID,
		SourceID:   req.SourceID,
		SourceName: "fastcampus_ewallet",
	}

	return s.PaymentRepo.InsertNewPaymentMethod(ctx, &paymentMethod)
}

func (s *PaymentService) PaymentMethodUnlink(ctx context.Context, userID uint64, req models.PaymentMethodLinkRequest) error {
	_, err := s.External.PaymentUnLink(ctx, req)
	if err != nil {
		return errors.Wrap(err, "failed to request payment unlink to e-wallet")
	}

	return s.PaymentRepo.DeletePaymentMethod(ctx, userID, req.SourceID, "fastcampus_ewallet")
}

func (s *PaymentService) InitiatePayment(ctx context.Context, req models.PaymentInitiatePayload) error {
	paymentMethod, err := s.PaymentRepo.GetPaymentMethod(ctx, req.UserID, "fastcampus_ewallet")
	if err != nil {
		return errors.Wrap(err, "failed to get payment method")
	}
	trxReq := external.PaymentTransactionRequest{
		Amount:          req.TotalPrice,
		Reference:       helpers.GenerateReference(),
		TransactionType: "DEBIT",
		WalletID:        paymentMethod.SourceID,
	}
	resp, err := s.External.WalletTransaction(ctx, trxReq)
	if err != nil {
		byteReq, _ := json.Marshal(req)
		s.External.ProduceKafkaMessage(ctx, helpers.GetEnv("KAFKA_TOPIC_PAYMENT_INITIATE", "payment-initiation-topic"), byteReq)

		return errors.Wrap(err, "failed to precees to wallet transaction")
	}

	helpers.Logger.WithField("balance", resp.Data.Balance).Info("succeed to payment")
	paymentTrx := models.PaymentTransaction{
		UserID:           req.UserID,
		OrderID:          req.OrderID,
		TotalPrice:       req.TotalPrice,
		PaymentMethodID:  paymentMethod.ID,
		Status:           "SUCCESS",
		PaymentReference: trxReq.Reference,
	}

	err = s.PaymentRepo.InsertNewPaymentTransaction(ctx, &paymentTrx)
	if err != nil {
		return errors.Wrap(err, "failed to insert to payment transaction")
	}

	_, err = s.External.OrderCallback(ctx, req.OrderID, "SUCCESS")
	if err != nil {
		return errors.Wrap(err, "failed to send callback to order service")
	}

	return nil
}

func (s *PaymentService) RefundPayment(ctx context.Context, req models.RefundPayload) error {
	paymentDetail, err := s.PaymentRepo.GetPaymentByOrderID(ctx, req.OrderID)
	if err != nil {
		return errors.Wrap(err, "failed to get payment detail")
	}

	paymentMethod, err := s.PaymentRepo.GetPaymentMethodByID(ctx, paymentDetail.PaymentMethodID)
	if err != nil {
		return errors.Wrap(err, "failed to get payment method")
	}
	trxReq := external.PaymentTransactionRequest{
		Amount:          paymentDetail.TotalPrice,
		Reference:       "REFUND" + paymentDetail.PaymentReference,
		TransactionType: "CREDIT",
		WalletID:        paymentMethod.SourceID,
	}
	resp, err := s.External.WalletTransaction(ctx, trxReq)
	if err != nil {
		return errors.Wrap(err, "failed to preccess to wallet transaction")
	}

	helpers.Logger.WithField("balance", resp.Data.Balance).Info("succeed to refund")

	refund := models.PaymentRefund{
		AdminID:          req.AdminID,
		OrderID:          req.OrderID,
		Status:           "SUCCESS",
		PaymentReference: trxReq.Reference,
	}

	err = s.PaymentRepo.InsertNewPaymentRefund(ctx, &refund)
	if err != nil {
		return errors.Wrap(err, "failed to insert to payment transaction")
	}

	return nil
}
