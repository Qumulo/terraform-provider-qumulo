package qumulo

import (
	"context"
	"github.com/hashicorp/terraform-plugin-log/tflog"
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
	err := setSSLSettings(ctx, d, m)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	return resourceSSLRead(ctx, d, m)
}

func resourceSSLRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	return diags
}

func resourceSSLUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	err := setSSLSettings(ctx, d, m)
	if err != nil {
		return diag.FromErr(err)
	}
	return resourceSSLRead(ctx, d, m)
}

func resourceSSLDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "Deleting SSL settings resource")
	var diags diag.Diagnostics

	return diags
}

func setSSLSettings(ctx context.Context, d *schema.ResourceData, m interface{}) error {
	c := m.(*Client)

	SSLConfig := SSLRequest{
		Certificate: d.Get("certificate").(string),
		PrivateKey:  d.Get("private_key").(string),
	}

	tflog.Debug(ctx, "Updating SSL settings")
	_, err := DoRequest[SSLRequest, SSLResponse](ctx, c, PUT, SSLEndpoint, &SSLConfig)
	return err
}
