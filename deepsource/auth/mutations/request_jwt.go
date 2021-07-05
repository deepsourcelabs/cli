package auth

import (
	"context"

	"github.com/deepsourcelabs/cli/deepsource/auth"
	"github.com/deepsourcelabs/graphql"
)

type RequestJWTParams struct {
	DeviceCode string `json:"deviceCode"`
}
type RequestJWTRequest struct {
	Params RequestJWTParams
}

const requestJWTMutation = `
mutation request($input:RequestJWTInput!) {
	requestJwt(input:$input) {
		payload
		token
		refreshToken
		refreshExpiresIn
	}
}`

type RequestJWTResponse struct {
	auth.JWT `json:"requestJwt"`
}

func (r RequestJWTRequest) Do(ctx context.Context, client IGQLClient) (*auth.JWT, error) {
	req := graphql.NewRequest(requestJWTMutation)
	req.Header.Set("Cache-Control", "no-cache")
	req.Var("input", r.Params)

	var res RequestJWTResponse
	if err := client.GQL().Run(ctx, req, &res); err != nil {
		return nil, err
	}

	return &res.JWT, nil
}
