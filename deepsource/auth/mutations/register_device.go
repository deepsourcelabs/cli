package auth

import (
	"context"

	"github.com/deepsourcelabs/cli/deepsource/auth"
	"github.com/deepsourcelabs/graphql"
)

const registerDeviceMutation = `mutation register {
	registerDevice(input:{}) {
		deviceCode
		userCode
		verificationUri
		verificationUriComplete
		expiresIn
		interval
	}
}`

type RegisterDeviceRequest struct{}

type RegisterDeviceResponse struct {
	auth.Device `json:"registerDevice"`
}

type IGQLClient interface {
	GQL() *graphql.Client
}

func (r RegisterDeviceRequest) Do(ctx context.Context, client IGQLClient) (*auth.Device, error) {
	req := graphql.NewRequest(registerDeviceMutation)
	req.Header.Set("Cache-Control", "no-cache")

	var res RegisterDeviceResponse
	if err := client.GQL().Run(ctx, req, &res); err != nil {
		return nil, err
	}

	return &res.Device, nil
}
