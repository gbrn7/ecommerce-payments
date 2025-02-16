package external

import (
	"bytes"
	"context"
	"ecommerce-payments/helpers"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

type OrderResponse struct {
	Message string `json:"message"`
}

func (e *External) OrderCallback(ctx context.Context, orderID int, status string) (OrderResponse, error) {
	url := helpers.GetEnv("ORDER_HOST", "") + fmt.Sprintf(helpers.GetEnv("ORDER_ENDPOINT_CALLBACK", ""), orderID)

	reqMap := map[string]string{
		"status": status,
	}

	bytePayload, err := json.Marshal(reqMap)

	if err != nil {
		return OrderResponse{}, errors.Wrap(err, "failed to marshaling request")
	}

	httpReq, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(bytePayload))
	if err != nil {
		return OrderResponse{}, errors.Wrap(err, "failed to create http request")
	}

	httpReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return OrderResponse{}, errors.Wrap(err, "failed to call order callback")
	}

	if resp.StatusCode != http.StatusOK {
		return OrderResponse{}, fmt.Errorf("got response failed from order callback resp : %d", resp.StatusCode)
	}

	response := OrderResponse{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return response, errors.Wrap(err, "failed to decode the response")
	}

	return response, nil
}
