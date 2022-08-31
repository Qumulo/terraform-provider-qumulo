package qumulo

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"terraform-provider-qumulo/openapi"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type Method int

const (
	GET Method = iota + 1
	PUT
	POST
	PATCH
	DELETE
)

func (m Method) String() string {
	return [...]string{"GET", "PUT", "POST", "PATCH", "DELETE"}[m-1]
}

type Client struct {
	HostURL     string
	HTTPClient  *http.Client
	BearerToken string
	Auth        AuthStruct
}

type AuthStruct struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthResponse struct {
	BearerToken string `json:"bearer_token"`
}

func NewClientGen(ctx context.Context, host, port, username, password *string) (*openapi.APIClient, error) {
	var apiClient *openapi.APIClient

	testHost := *host + ":" + *port
	const testScheme = "https"

	cfg := openapi.NewConfiguration()
	cfg.Debug = true
	cfg.Host = testHost
	cfg.Scheme = testScheme
	transCfg := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	cfg.HTTPClient = &http.Client{Timeout: 10 * time.Second, Transport: transCfg} // ignore expired SSL certificates

	apiClient = openapi.NewAPIClient(cfg)
	authReq := openapi.V1SessionLoginPostRequest{Username: username, Password: password}
	apiRes, _, err := apiClient.SessionApi.V1SessionLoginPost(context.Background()).
		V1SessionLoginPostRequest(authReq).Execute()
	if err != nil {
		return nil, err
	}
	bearerToken := "Bearer " + *apiRes.BearerToken

	cfg.AddDefaultHeader("Authorization", bearerToken)
	tflog.Info(ctx, "Qumulo client configured", map[string]interface{}{
		"host":     host,
		"port":     port,
		"username": username,
	})

	return apiClient, nil

}

func NewClient(ctx context.Context, host, port, username, password *string) (*Client, error) {
	HostURL := fmt.Sprintf("https://%s:%s", *host, *port)

	transCfg := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // ignore expired SSL certificates
	}

	c := Client{
		HTTPClient: &http.Client{Timeout: 10 * time.Second, Transport: transCfg},
		HostURL:    HostURL,
		Auth: AuthStruct{
			Username: *username,
			Password: *password,
		},
	}

	ar, err := c.SignIn(ctx)
	if err != nil {
		return nil, err
	}

	c.BearerToken = ar.BearerToken
	c.HostURL = HostURL

	tflog.Info(ctx, "Qumulo client configured", map[string]interface{}{
		"host":     host,
		"port":     port,
		"username": username,
	})
	return &c, nil
}

func DoRequest[RQ interface{}, R interface{}](ctx context.Context, client *Client, method Method, endpointUri string, reqBody *RQ) (*R, error) {
	bearerToken := "Bearer " + client.BearerToken
	HostURL := client.HostURL

	var parsedReqBody io.Reader

	if reqBody != nil {
		rb, err := json.Marshal(reqBody)
		if err != nil {
			return nil, err
		}
		parsedReqBody = strings.NewReader(string(rb))
	} else {
		parsedReqBody = nil
	}
	url := fmt.Sprintf("%s%s", HostURL, endpointUri)
	req, err := http.NewRequest(method.String(), url, parsedReqBody)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", bearerToken)
	req.Header.Add("Content-Type", "application/json")

	tflog.Trace(ctx, "Executing API request", map[string]interface{}{
		"url":    url,
		"method": method.String(),
	})

	body, err := client.makeHTTPRequest(req)
	if err != nil {
		return nil, err
	}

	var cr R
	if len(body) == 0 {
		return nil, nil
	}

	err = json.Unmarshal(body, &cr)
	if err != nil {
		return nil, err
	}
	return &cr, nil
}

func (c *Client) makeHTTPRequest(req *http.Request) ([]byte, error) {

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusAccepted {
		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, body)
	}

	return body, err
}
