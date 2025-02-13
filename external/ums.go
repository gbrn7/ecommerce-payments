package external

import (
	"context"
	"ecommerce-payments/helpers"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

type Profile struct {
	Message string `json:"message"`
	Data    struct {
		ID          uint64 `json:"id"`
		Username    string `json:"username"`
		Fullname    string `json:"full_name"`
		Email       string `json:"email"`
		Phonenumber string `json:"phone_number"`
		Address     string `json:"address"`
		Dob         string `json:"dob"`
		Role        string `json:"role"`
	} `json:"data"`
}

type External struct {
}

func (ext *External) GetProfile(ctx context.Context, token string) (Profile, error) {
	url := helpers.GetEnv("UMS_HOST", "") + helpers.GetEnv("UMS_ENDPOINT_PROFILE", "")
	httpReq, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return Profile{}, errors.Wrap(err, "failed to create http request")
	}

	httpReq.Header.Set("Authorization", token)
	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return Profile{}, errors.Wrap(err, "failed to call ums get profile")
	}

	if resp.StatusCode != http.StatusOK {
		return Profile{}, fmt.Errorf("got response failed from ums get profile resp : %d", resp.StatusCode)
	}

	response := Profile{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return response, errors.Wrap(err, "failed to decode the response")
	}

	return response, nil
}
