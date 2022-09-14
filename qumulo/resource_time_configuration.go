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

func resourceTimeConfiguration() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTimeConfigurationCreate,
		ReadContext:   resourceTimeConfigurationRead,
		UpdateContext: resourceTimeConfigurationUpdate,
		DeleteContext: resourceTimeConfigurationDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"use_ad_for_primary": &schema.Schema{
				Type:     schema.TypeBool,
				Required: true,
			},
			"ntp_servers": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceTimeConfigurationCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*openapi.APIClient)

	config := setTimeConfiguration(ctx, d, m)

	resp, r, err := c.TimeApi.V1TimeSettingsPut(context.Background()).
		V1TimeSettingsGet200Response(config).Execute()
	if err != nil {
		tflog.Debug(ctx, fmt.Sprintf("Error when calling `TimeApi.V1TimeSettingsPut``: %v\n", err))
		tflog.Debug(ctx, fmt.Sprintf("Full HTTP response: %v\n", r))
		return diag.FromErr(err)
	}
	// response from `V1TimeSettingsPut`: V1TimeSettingsGet200Response
	tflog.Debug(ctx, fmt.Sprintf("Response from `TimeApi.V1TimeSettingsPut`: %v\n", resp))

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return resourceTimeConfigurationRead(ctx, d, m)
}

func resourceTimeConfigurationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*openapi.APIClient)

	var errs ErrorCollection

	resp, r, err := c.TimeApi.V1TimeSettingsGet(context.Background()).Execute()
	if err != nil {
		tflog.Debug(ctx, fmt.Sprintf("Error when calling `TimeApi.V1TimeSettingsGet``: %v\n", err))
		tflog.Debug(ctx, fmt.Sprintf("Full HTTP response: %v\n", r))
		return diag.FromErr(err)
	}
	// response from `V1TimeSettingsGet`: V1TimeSettingsGet200Response
	tflog.Debug(ctx, fmt.Sprintf("Response from `TimeApi.V1TimeSettingsGet`: %v\n", resp))

	errs.addMaybeError(d.Set("use_ad_for_primary", resp.GetUseAdForPrimary()))
	errs.addMaybeError(d.Set("ntp_servers", resp.GetNtpServers()))

	return errs.diags
}

func resourceTimeConfigurationUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*openapi.APIClient)

	config := setTimeConfiguration(ctx, d, m)

	resp, r, err := c.TimeApi.V1TimeSettingsPatch(context.Background()).
		V1TimeSettingsGet200Response(config).Execute()
	if err != nil {
		tflog.Debug(ctx, fmt.Sprintf("Error when calling `TimeApi.V1TimeSettingsPatch``: %v\n", err))
		tflog.Debug(ctx, fmt.Sprintf("Full HTTP response: %v\n", r))
		return diag.FromErr(err)
	}
	// response from `V1TimeSettingsPatch`: V1TimeSettingsGet200Response
	tflog.Debug(ctx, fmt.Sprintf("Response from `TimeApi.V1TimeSettingsPatch`: %v\n", resp))

	return resourceTimeConfigurationRead(ctx, d, m)
}

func resourceTimeConfigurationDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "Deleting time configuration resource")

	return nil
}

func setTimeConfiguration(ctx context.Context, d *schema.ResourceData, m interface{}) openapi.V1TimeSettingsGet200Response {

	config := *openapi.NewV1TimeSettingsGet200Response()

	ntpServers := InterfaceSliceToStringSlice(d.Get("ntp_servers").([]interface{}))

	config.SetUseAdForPrimary(d.Get("use_ad_for_primary").(bool))
	config.SetNtpServers(ntpServers)

	tflog.Debug(ctx, "Updating time configuration")

	return config
}
