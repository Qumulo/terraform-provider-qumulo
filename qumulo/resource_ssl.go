package qumulo

import (
	"context"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const SslEndpoint = "/v2/cluster/settings/ssl/certificate"

type SslRequest struct {
	Certificate string `json:"certificate"`
	PrivateKey  string `json:"private_key"`
}

// TODO: Figure out what the proper response for an SSL update is
type SslResponse struct {
	Placeholder string `json:"placeholder"`
}

func resourceSsl() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSslCreate,
		ReadContext:   resourceSslRead,
		UpdateContext: resourceSslUpdate,
		DeleteContext: resourceSslDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
		},

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

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceSslCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	err := setSslSettings(ctx, d, m)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return resourceSslRead(ctx, d, m)
}

func resourceSslRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}

func resourceSslUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	err := setSslSettings(ctx, d, m)
	if err != nil {
		return diag.FromErr(err)
	}
	return resourceSslRead(ctx, d, m)
}

func resourceSslDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "Deleting SSL settings resource")

	return nil
}

func setSslSettings(ctx context.Context, d *schema.ResourceData, m interface{}) error {
	c := m.(*Client)

	sslConfig := SslRequest{
		Certificate: d.Get("certificate").(string),
		PrivateKey:  d.Get("private_key").(string),
	}

	tflog.Debug(ctx, "Updating SSL settings")
	_, err := DoRequest[SslRequest, SslResponse](ctx, c, PUT, SslEndpoint, &sslConfig)
	return err
}
