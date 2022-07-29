package qumulo

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const AuthEndpoint = "/v1/session/login"

// SignIn - Get a new token for user
func (c *Client) SignIn(ctx context.Context) (*AuthResponse, error) {
	if c.Auth.Username == "" || c.Auth.Password == "" {
		return nil, fmt.Errorf("cannot sign in: missing username/password")
	}

	ar, err := DoRequest[AuthStruct, AuthResponse](ctx, c, POST, AuthEndpoint, &c.Auth)
	if err != nil {
		tflog.Error(ctx, "Fetching auth token failed")
		return nil, err
	}

	return ar, nil
}
