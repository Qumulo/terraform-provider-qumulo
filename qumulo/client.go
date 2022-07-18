package qumulo

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// Client -
type Client struct {
	HostURL      string
	HTTPClient   *http.Client
	Bearer_Token string
	Auth         AuthStruct
}

// AuthStruct -
type AuthStruct struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// AuthResponse -
type AuthResponse struct {
	Bearer_Token string `json:"bearer_token"`
}

// NewClient -
func NewClient(host, port, username, password *string) (*Client, error) {
	HostURL := fmt.Sprintf("https://%s:%s", *host, *port)

	transCfg := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // ignore expired SSL certificates
	}

	c := Client{
		HTTPClient: &http.Client{Timeout: 10 * time.Second, Transport: transCfg},
		// Default Qumulo URL
		HostURL: HostURL,
		Auth: AuthStruct{
			Username: *username,
			Password: *password,
		},
	}

	if host != nil {
		c.HostURL = *host
	}

	ar, err := c.SignIn()
	if err != nil {
		return nil, err
	}

	c.Bearer_Token = ar.Bearer_Token
	c.HostURL = HostURL

	return &c, nil
}

func (c *Client) doRequest(req *http.Request) ([]byte, error) {

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, body)
	}

	return body, err
}
