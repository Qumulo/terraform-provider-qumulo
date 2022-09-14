package qumulo

import (
	"context"
	"fmt"
	"strconv"
	"terraform-provider-qumulo/openapi"
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
			"krb5p_enabled": &schema.Schema{
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
	c := m.(*openapi.APIClient)
	settings := setNfsSettings(ctx, d, m)

	resp, r, err := c.NfsApi.V2NfsSettingsPut(context.Background()).V2NfsSettingsGet200Response(settings).Execute()
	if err != nil {
		tflog.Debug(ctx, fmt.Sprintf("Error when calling `NfsApi.V2NfsSettingsPut``: %v\n", err))
		tflog.Debug(ctx, fmt.Sprintf("Full HTTP response: %v\n", r))
		return diag.FromErr(err)
	}
	// response from `V2NfsSettingsPut`: V2NfsSettingsGet200Response
	tflog.Debug(ctx, fmt.Sprintf("Response from `NfsApi.V2NfsSettingsPut`: %v\n", resp))

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	return resourceNfsSettingsRead(ctx, d, m)
}

func resourceNfsSettingsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*openapi.APIClient)

	var errs ErrorCollection

	resp, r, err := c.NfsApi.V2NfsSettingsGet(context.Background()).Execute()
	if err != nil {
		tflog.Debug(ctx, fmt.Sprintf("Error when calling `NfsApi.V2NfsSettingsGet``: %v\n", err))
		tflog.Debug(ctx, fmt.Sprintf("Full HTTP response: %v\n", r))
		return diag.FromErr(err)
	}
	// response from `V2NfsSettingsGet`: V2NfsSettingsGet200Response
	tflog.Debug(ctx, fmt.Sprintf("Response from `NfsApi.V2NfsSettingsGet`: %v\n", resp))

	errs.addMaybeError(d.Set("v4_enabled", resp.GetV4Enabled()))
	errs.addMaybeError(d.Set("krb5_enabled", resp.GetKrb5Enabled()))
	errs.addMaybeError(d.Set("krb5p_enabled", resp.GetKrb5pEnabled()))
	errs.addMaybeError(d.Set("auth_sys_enabled", resp.GetAuthSysEnabled()))

	return errs.diags
}

func resourceNfsSettingsUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*openapi.APIClient)
	settings := setNfsSettings(ctx, d, m)

	resp, r, err := c.NfsApi.V2NfsSettingsPatch(context.Background()).V2NfsSettingsGet200Response(settings).Execute()
	if err != nil {
		tflog.Debug(ctx, fmt.Sprintf("Error when calling `NfsApi.V2NfsSettingsPatch``: %v\n", err))
		tflog.Debug(ctx, fmt.Sprintf("Full HTTP response: %v\n", r))
		return diag.FromErr(err)
	}
	// response from `V2NfsSettingsPatch`: V2NfsSettingsGet200Response
	tflog.Debug(ctx, fmt.Sprintf("Response from `NfsApi.V2NfsSettingsPatch`: %v\n", resp))

	return resourceNfsSettingsRead(ctx, d, m)
}

func resourceNfsSettingsDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "Deleting NFS settings resource")

	return nil
}

func setNfsSettings(ctx context.Context, d *schema.ResourceData, m interface{}) openapi.V2NfsSettingsGet200Response {

	settings := *openapi.NewV2NfsSettingsGet200Response()
	settings.SetV4Enabled(d.Get("v4_enabled").(bool))
	settings.SetKrb5Enabled(d.Get("krb5_enabled").(bool))
	settings.SetKrb5pEnabled(d.Get("krb5p_enabled").(bool))
	settings.SetAuthSysEnabled(d.Get("auth_sys_enabled").(bool))

	tflog.Debug(ctx, "Updating NFS Settings")

	return settings
}
