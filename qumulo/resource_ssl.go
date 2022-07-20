package qumulo

import (
	"context"
	"strconv"
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

	sslr := SSLRequest{
		Certificate: certificate,
		PrivateKey:  key,
	}

	_, err := DoRequest[SSLRequest, SSLResponse](c, PUT, SSLEndpoint, &sslr)
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
