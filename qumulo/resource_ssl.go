package qumulo

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const SSLEndpoint = "/v2/cluster/settings/ssl/certificate"

type SSLRequest struct {
	Certificate string `json:"certificate"`
	PrivateKey  string `json:"private_key"`
}

// TODO: Figure out what the proper response for an SSL update is
type SSLResponse struct {
	Placeholder string `json:"placeholder"`
}

func resourceSSL() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSSLCreate,
		ReadContext:   resourceSSLRead,
		UpdateContext: resourceSSLUpdate,
		DeleteContext: resourceSSLDelete,
		Schema: map[string]*schema.Schema{
			"certificate": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"private_key": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceSSLCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	certificate := d.Get("certificate").(string)
	key := d.Get("private_key").(string)

	updatedSSL := SSLRequest{
		Certificate: certificate,
		PrivateKey:  key,
	}

	_, err := c.UpdateSSL(updatedSSL)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}

func resourceSSLRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	return diags
}

func resourceSSLUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceSSLCreate(ctx, d, m)
}

func resourceSSLDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	return diags
}

func (c *Client) UpdateSSL(sslReq SSLRequest) (*SSLResponse, error) {
	bearerToken := "Bearer " + c.Bearer_Token

	HostURL := c.HostURL

	rb, err := json.Marshal(sslReq)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/v2/cluster/settings/ssl/certificate", HostURL),
		strings.NewReader(string(rb)))
	req.Header.Set("Authorization", bearerToken)
	req.Header.Add("Content-Type", "application/json")
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	cr := SSLResponse{}
	err = json.Unmarshal(body, &cr)
	if err != nil {
		return nil, err
	}

	return &cr, nil
}
