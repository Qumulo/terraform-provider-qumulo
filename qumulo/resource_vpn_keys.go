package qumulo

import (
	"context"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strconv"
	"time"
)

const VpnKeysEndpoint = "/v1/support/vpn-keys"

type VpnKeysBody struct {
	MqvpnClientCrt string `json:"mqvpn_client_crt"`
	MqvpnClientKey string `json:"mqvpn_client_key"`
	QumuloCaCrt    string `json:"qumulo_ca_crt"`
}

func resourceVpnKeys() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVpnKeysCreate,
		ReadContext:   resourceVpnKeysRead,
		UpdateContext: resourceVpnKeysUpdate,
		DeleteContext: resourceVpnKeysDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"mqvpn_client_crt": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"mqvpn_client_key": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"qumulo_ca_crt": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceVpnKeysCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	err := setVpnkeys(ctx, d, m, PUT)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return resourceVpnKeysRead(ctx, d, m)
}

func resourceVpnKeysRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	var errs ErrorCollection

	vpnKeys, err := DoRequest[VpnKeysBody, VpnKeysBody](ctx, c, GET, VpnKeysEndpoint, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	errs.addMaybeError(d.Set("mqvpn_client_crt", vpnKeys.MqvpnClientCrt))
	errs.addMaybeError(d.Set("mqvpn_client_key", vpnKeys.MqvpnClientKey))
	errs.addMaybeError(d.Set("qumulo_ca_crt", vpnKeys.QumuloCaCrt))

	return errs.diags
}

func resourceVpnKeysUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	err := setVpnkeys(ctx, d, m, PATCH)
	if err != nil {
		return diag.FromErr(err)
	}
	return resourceVpnKeysRead(ctx, d, m)
}

func resourceVpnKeysDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "Deleting Ftp server resource")
	return nil
}

func setVpnkeys(ctx context.Context, d *schema.ResourceData, m interface{}, method Method) error {
	c := m.(*Client)

	vpnKeys := VpnKeysBody{
		MqvpnClientCrt: d.Get("mqvpn_client_crt").(string),
		MqvpnClientKey: d.Get("mqvpn_client_key").(string),
		QumuloCaCrt:    d.Get("qumulo_ca_crt").(string),
	}
	tflog.Debug(ctx, "Setting Vpn keys")
	_, err := DoRequest[VpnKeysBody, VpnKeysBody](ctx, c, method, VpnKeysEndpoint, &vpnKeys)
	return err
}
