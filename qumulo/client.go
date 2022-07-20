package qumulo

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type Method int

const (
	GET Method = iota + 1
	PUT
	POST
)

func (m Method) String() string {
	return [...]string{"GET", "PUT", "POST"}[m-1]
}

// Client -
type Client struct {
	HostURL     string
	HTTPClient  *http.Client
	BearerToken string
	Auth        AuthStruct
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
	BearerToken string `json:"bearer_token"`
}

// NewClient -
func NewClient(host, port, username, password *string) (*Client, error) {
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

	ar, err := c.SignIn()
	if err != nil {
		return nil, err
	}

	c.BearerToken = ar.BearerToken
	c.HostURL = HostURL

	return &c, nil
}

func (c *Client) MakeHTTPRequest(req *http.Request) ([]byte, error) {

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

func DoRequest[RQ interface{}, R interface{}](client *Client, method Method, endpointUri string, reqBody *RQ) (*R, error) {
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

	req, err := http.NewRequest(method.String(), fmt.Sprintf("%s%s", HostURL, endpointUri), parsedReqBody)
	req.Header.Set("Authorization", bearerToken)
	req.Header.Add("Content-Type", "application/json")
	if err != nil {
		return nil, err
	}

	body, err := client.MakeHTTPRequest(req)
	if err != nil {
		return nil, err
	}

	var cr R
	err = json.Unmarshal(body, &cr)
	if err != nil {
		return nil, err
	}
	return &cr, nil
}
