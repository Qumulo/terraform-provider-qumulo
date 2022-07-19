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

type SSLCARequest struct {
	Certificate string `json:"ca_certificate"`
}

// TODO: Figure out what the proper response for an SSL CA update is
type SSLCAResponse struct {
	Placeholder string `json:"placeholder"`
}

func resourceSSLCA() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSSLCACreate,
		ReadContext:   resourceSSLCARead,
		UpdateContext: resourceSSLCAUpdate,
		DeleteContext: resourceSSLCADelete,
		Schema: map[string]*schema.Schema{
			"ca_certificate": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceSSLCACreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	certificate := d.Get("ca_certificate").(string)

	updatedSSLCA := SSLCARequest{
		Certificate: certificate,
	}

	_, err := c.UpdateSSLCA(updatedSSLCA)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}

func resourceSSLCARead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	return diags
}

func resourceSSLCAUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	certificate := d.Get("ca_certificate").(string)

	updatedSSLCA := SSLCARequest{
		Certificate: certificate,
	}
	_, err := c.UpdateSSLCA(updatedSSLCA)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}

func resourceSSLCADelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	return diags
}

func (c *Client) UpdateSSLCA(sslReq SSLCARequest) (*SSLCAResponse, error) {
	bearerToken := "Bearer " + c.Bearer_Token

	HostURL := c.HostURL

	rb, err := json.Marshal(sslReq)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/v2/cluster/settings/ssl/ca-certificate", HostURL),
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

	cr := SSLCAResponse{}
	err = json.Unmarshal(body, &cr)
	if err != nil {
		return nil, err
	}

	return &cr, nil
}
