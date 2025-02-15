package external

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"ecommerce-payments/helpers"
	"ecommerce-payments/internal/models"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/pkg/errors"
)

type PaymentLinkResponse struct {
	Message string `json:"message"`
	Data    struct {
		OTP string `json:"otp"`
	} `json:"data"`
}

func (e *External) generateSignature(ctx context.Context, payload, timestamp, endpoint string) string {
	secretKey := helpers.GetEnv("WALLET_CLIENT_SECRET", "")

	re := regexp.MustCompile(`[^a-zA-Z0-9]+`)
	payload = re.ReplaceAllString(payload, "")
	payload = strings.ToLower(payload) + timestamp + endpoint

	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write([]byte(payload))
	return hex.EncodeToString(h.Sum(nil))
}

func (e *External) PaymentLink(ctx context.Context, req models.PaymentMethodLinkRequest) (PaymentLinkResponse, error) {
	url := helpers.GetEnv("WALLET_HOST", "") + helpers.GetEnv("WALLET_ENDPOINT_LINK", "")

	reqMap := map[string]int{
		"wallet_id": req.SourceID,
	}

	bytePayload, err := json.Marshal(reqMap)
	if err != nil {
		return PaymentLinkResponse{}, errors.Wrap(err, "failed to marshaling request")
	}

	httpReq, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(bytePayload))
	if err != nil {
		return PaymentLinkResponse{}, errors.Wrap(err, "failed to create http request")
	}

	timestamp := time.Now().Format(time.RFC3339)
	signature := e.generateSignature(ctx, string(bytePayload), timestamp, helpers.GetEnv("WALLET_ENDPOINT_LINK", ""))

	httpReq.Header.Set("Client-Id", helpers.GetEnv("WALLET_CLIENT_ID", ""))
	httpReq.Header.Set("Timestamp", timestamp)
	httpReq.Header.Set("Signature", signature)

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return PaymentLinkResponse{}, errors.Wrap(err, "failed to call wallet link wallet")
	}

	if resp.StatusCode != http.StatusOK {
		return PaymentLinkResponse{}, fmt.Errorf("got response failed from wallet link wallet resp : %d", resp.StatusCode)
	}

	response := PaymentLinkResponse{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return response, errors.Wrap(err, "failed to decode the response")
	}

	return response, nil
}

func (e *External) PaymentUnLink(ctx context.Context, req models.PaymentMethodLinkRequest) (PaymentLinkResponse, error) {
	url := helpers.GetEnv("WALLET_HOST", "") + fmt.Sprintf(helpers.GetEnv("WALLET_ENDPOINT_UNLINK", ""), req.SourceID)

	httpReq, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return PaymentLinkResponse{}, errors.Wrap(err, "failed to create http request")
	}

	timestamp := time.Now().Format(time.RFC3339)
	signature := e.generateSignature(ctx, "", timestamp, fmt.Sprintf(helpers.GetEnv("WALLET_ENDPOINT_UNLINK", ""), req.SourceID))

	httpReq.Header.Set("Client-Id", helpers.GetEnv("WALLET_CLIENT_ID", ""))
	httpReq.Header.Set("Timestamp", timestamp)
	httpReq.Header.Set("Signature", signature)

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return PaymentLinkResponse{}, errors.Wrap(err, "failed to call wallet unlink wallet")
	}

	if resp.StatusCode != http.StatusOK {
		return PaymentLinkResponse{}, fmt.Errorf("got response failed from wallet unlink wallet resp : %d", resp.StatusCode)
	}

	response := PaymentLinkResponse{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return response, errors.Wrap(err, "failed to decode the response")
	}

	return response, nil
}

func (e *External) PaymentLinkConfirmation(ctx context.Context, walletID int, otp string) (PaymentLinkResponse, error) {
	url := helpers.GetEnv("WALLET_HOST", "") + fmt.Sprintf(helpers.GetEnv("WALLET_ENDPOINT_LINK_CONFIRM", ""), walletID)

	reqMap := map[string]string{
		"otp": otp,
	}

	bytePayload, err := json.Marshal(reqMap)

	if err != nil {
		return PaymentLinkResponse{}, errors.Wrap(err, "failed to marshaling request")
	}

	httpReq, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(bytePayload))
	if err != nil {
		return PaymentLinkResponse{}, errors.Wrap(err, "failed to create http request")
	}

	timestamp := time.Now().Format(time.RFC3339)
	signature := e.generateSignature(ctx, string(bytePayload), timestamp, fmt.Sprintf(helpers.GetEnv("WALLET_ENDPOINT_LINK_CONFIRM", ""), walletID))

	httpReq.Header.Set("Client-Id", helpers.GetEnv("WALLET_CLIENT_ID", ""))
	httpReq.Header.Set("Timestamp", timestamp)
	httpReq.Header.Set("Signature", signature)

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return PaymentLinkResponse{}, errors.Wrap(err, "failed to call wallet link confirm ")
	}

	if resp.StatusCode != http.StatusOK {
		return PaymentLinkResponse{}, fmt.Errorf("got response failed from wallet link confirm  resp : %d", resp.StatusCode)
	}

	response := PaymentLinkResponse{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return response, errors.Wrap(err, "failed to decode the response")
	}

	return response, nil
}
