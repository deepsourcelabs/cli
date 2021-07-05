// Package sdk is the sdk for DeepSource public APIs.
//
// This is the GraphQL client for DeepSource public APIs. In the future,
// we hope to publish this package as a standalone package.
package deepsource

import (
	"context"

	"github.com/deepsourcelabs/cli/deepsource/auth"
	authmut "github.com/deepsourcelabs/cli/deepsource/auth/mutations"
	"github.com/deepsourcelabs/graphql"
)

const host = "http://localhost:8000/graphql/"

type Client struct {
	gql *graphql.Client
}

func (c Client) GQL() *graphql.Client {
	return c.gql
}

func New() *Client {
	gql := graphql.NewClient(host)
	return &Client{
		gql: gql,
	}
}

func (c Client) RegisterDevice(ctx context.Context) (*auth.Device, error) {
	req := authmut.RegisterDeviceRequest{}
	res, err := req.Do(ctx, c)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c Client) Login(ctx context.Context, deviceCode string) (*auth.JWT, error) {
	req := authmut.RequestJWTRequest{
		Params: authmut.RequestJWTParams{
			DeviceCode: deviceCode,
		},
	}

	res, err := req.Do(ctx, c)
	if err != nil {
		return nil, err
	}
	return res, nil
}
