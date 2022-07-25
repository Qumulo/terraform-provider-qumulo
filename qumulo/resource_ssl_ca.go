package qumulo

import (
	"context"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const SSLCAEndpoint = "/v2/cluster/settings/ssl/ca-certificate"

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

	SSLCAConfig := SSLCARequest{
		Certificate: d.Get("ca_certificate").(string),
	}

	_, err := DoRequest[SSLCARequest, SSLCAResponse](c, PUT, SSLCAEndpoint, &SSLCAConfig)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return resourceSSLCARead(ctx, d, m)
}

func resourceSSLCARead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	c := m.(*Client)

	var diags diag.Diagnostics

	cert, err := DoRequest[SSLCARequest, SSLCARequest](c, GET, SSLCAEndpoint, nil)

	if err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("ca_certificate", cert.Certificate); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceSSLCAUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceSSLCACreate(ctx, d, m)
}

func resourceSSLCADelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	return diags
}
