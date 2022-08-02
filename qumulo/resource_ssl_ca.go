package qumulo

import (
	"context"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const SslCaEndpoint = "/v2/cluster/settings/ssl/ca-certificate"

type SslCaBody struct {
	CaCertificate string `json:"ca_certificate"`
}

func resourceSslCa() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSslCaCreate,
		ReadContext:   resourceSslCaRead,
		UpdateContext: resourceSslCaUpdate,
		DeleteContext: resourceSslCaDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"ca_certificate": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceSslCaCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	err := setSslCaSettings(ctx, d, m)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return resourceSslCaRead(ctx, d, m)
}

func resourceSslCaRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	var diags diag.Diagnostics

	cert, err := DoRequest[SslCaBody, SslCaBody](ctx, c, GET, SslCaEndpoint, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("ca_certificate", cert.CaCertificate); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceSslCaUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	err := setSslCaSettings(ctx, d, m)
	if err != nil {
		return diag.FromErr(err)
	}
	return resourceSslCaRead(ctx, d, m)
}

func resourceSslCaDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "Deleting SSL CA Settings")
	c := m.(*Client)

	var diags diag.Diagnostics

	_, err := DoRequest[SslCaBody, SslCaBody](ctx, c, DELETE, SslCaEndpoint, nil)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}

func setSslCaSettings(ctx context.Context, d *schema.ResourceData, m interface{}) error {
	c := m.(*Client)

	sslCaConfig := SslCaBody{
		CaCertificate: d.Get("ca_certificate").(string),
	}

	tflog.Debug(ctx, "Updating SSL CA settings")
	_, err := DoRequest[SslCaBody, SslCaBody](ctx, c, PUT, SslCaEndpoint, &sslCaConfig)
	return err
}
