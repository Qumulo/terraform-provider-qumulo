package qumulo

import (
	"context"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const NfsSettingsEndpoint = "/v2/nfs/settings"

type NfsSettingsBody struct {
	V4Enabled      bool `json:"v4_enabled"`
	Krb5Enabled    bool `json:"krb5_enabled"`
	AuthSysEnabled bool `json:"auth_sys_enabled"`
}

func resourceNfsSettings() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNfsSettingsCreate,
		ReadContext:   resourceNfsSettingsRead,
		UpdateContext: resourceNfsSettingsUpdate,
		DeleteContext: resourceNfsSettingsDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"v4_enabled": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"krb5_enabled": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"auth_sys_enabled": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
		},

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceNfsSettingsCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	err := setNfsSettings(ctx, d, m, PUT)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	return resourceNfsSettingsRead(ctx, d, m)
}

func resourceNfsSettingsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	var errs ErrorCollection
	s, err := DoRequest[NfsSettingsBody, NfsSettingsBody](ctx, c, GET, NfsSettingsEndpoint, nil)
	if err != nil {
		return diag.FromErr(err)
	}
	errs.addMaybeError(d.Set("v4_enabled", s.V4Enabled))
	errs.addMaybeError(d.Set("krb5_enabled", s.Krb5Enabled))
	errs.addMaybeError(d.Set("auth_sys_enabled", s.AuthSysEnabled))

	return errs.diags
}

func resourceNfsSettingsUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	err := setNfsSettings(ctx, d, m, PATCH)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	return resourceNfsSettingsRead(ctx, d, m)
}

func resourceNfsSettingsDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "Deleting NFS settings resource")

	return nil
}

func setNfsSettings(ctx context.Context, d *schema.ResourceData, m interface{}, method Method) error {
	c := m.(*Client)

	nfsSettings := NfsSettingsBody{
		V4Enabled:      d.Get("v4_enabled").(bool),
		Krb5Enabled:    d.Get("krb5_enabled").(bool),
		AuthSysEnabled: d.Get("auth_sys_enabled").(bool),
	}

	tflog.Debug(ctx, "Updating NFS Settings")
	_, err := DoRequest[NfsSettingsBody, NfsSettingsBody](ctx, c, method, NfsSettingsEndpoint, &nfsSettings)
	return err
}
