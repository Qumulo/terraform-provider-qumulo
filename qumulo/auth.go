package qumulo

import (
	"fmt"
)

const AuthEndpoint = "/v1/session/login"

// SignIn - Get a new token for user
func (c *Client) SignIn() (*AuthResponse, error) {
	if c.Auth.Username == "" || c.Auth.Password == "" {
		return nil, fmt.Errorf("define username and password")
	}

	ar, err := DoRequest[AuthStruct, AuthResponse](c, POST, AuthEndpoint, &c.Auth)
	if err != nil {
		return nil, err
	}

	return ar, nil
}
