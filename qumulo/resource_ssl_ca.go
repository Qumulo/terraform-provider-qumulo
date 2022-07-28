package qumulo

import (
	"context"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const SslCaEndpoint = "/v2/cluster/settings/ssl/ca-certificate"

type SslCaRequest struct {
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
	}
}

func resourceSslCaCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	sslCaConfig := SslCaRequest{
		CaCertificate: d.Get("ca_certificate").(string),
	}

	_, err := DoRequest[SslCaRequest, SslCaRequest](c, PUT, SslCaEndpoint, &sslCaConfig)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return resourceSslCaRead(ctx, d, m)
}

func resourceSslCaRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	var diags diag.Diagnostics

	cert, err := DoRequest[SslCaRequest, SslCaRequest](c, GET, SslCaEndpoint, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("ca_certificate", cert.CaCertificate); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceSslCaUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceSslCaCreate(ctx, d, m)
}

func resourceSslCaDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	var diags diag.Diagnostics

	_, err := DoRequest[SslCaRequest, SslCaRequest](c, DELETE, SslCaEndpoint, nil)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}
