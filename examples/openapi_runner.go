package main

import (
	"context"
	"crypto/tls"
	"net/http"
	"terraform-provider-qumulo/openapi"
	"time"
)

func main() {
	var client *openapi.APIClient

	const testHost = "10.116.100.137:22582"
	const testScheme = "https"

	cfg := openapi.NewConfiguration()
	cfg.Debug = true
	cfg.Host = testHost
	cfg.Scheme = testScheme
	transCfg := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	cfg.HTTPClient = &http.Client{Timeout: 10 * time.Second, Transport: transCfg} // ignore expired SSL certificates

	client = openapi.NewAPIClient(cfg)
	userName := "admin"
	password := "Admin123"
	authReq := openapi.V1SessionLoginPostRequest{Username: &userName, Password: &password}
	apiRes, _, err := client.SessionApi.V1SessionLoginPost(context.Background()).V1SessionLoginPostRequest(authReq).Execute()
	if err != nil {
		print(err)
	}
	print(*apiRes.BearerToken)
}
