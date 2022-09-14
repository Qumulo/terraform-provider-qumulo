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

func resourceMonitoring() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMonitoringCreate,
		ReadContext:   resourceMonitoringRead,
		UpdateContext: resourceMonitoringUpdate,
		DeleteContext: resourceMonitoringDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"enabled": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"mq_host": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"mq_port": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"mq_proxy_host": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"mq_proxy_port": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"s3_proxy_host": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"s3_proxy_port": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"s3_proxy_disable_https": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"vpn_enabled": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"vpn_host": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"period": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
		},

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceMonitoringCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*openapi.APIClient)

	ms := setMonitoringSettings(ctx, d, m)
	resp, r, err := c.SupportApi.V1SupportSettingsPut(context.Background()).
		V1SupportSettingsGet200Response(*ms).Execute()
	if err != nil {
		tflog.Debug(ctx, fmt.Sprintf("Error when calling `SupportApi.V1SupportSettingsPut``: %v\n", err))
		tflog.Debug(ctx, fmt.Sprintf("Full HTTP response: %v\n", r))
		return diag.FromErr(err)
	}
	// response from `V1SupportSettingsPut`: V1AdJoinPost202Response
	tflog.Debug(ctx, fmt.Sprintf("Response from `SupportApi.V1SupportSettingsPut`: %v\n", resp))

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return resourceMonitoringRead(ctx, d, m)
}

func resourceMonitoringRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*openapi.APIClient)

	var errs ErrorCollection

	// settings, err := DoRequest[MonitoringSettings, MonitoringSettings](ctx, c, GET, MonitoringEndpoint, nil)
	resp, r, err := c.SupportApi.V1SupportSettingsGet(context.Background()).Execute()
	if err != nil {
		tflog.Debug(ctx, fmt.Sprintf("Error when calling `SupportApi.V1SupportSettingsGet``: %v\n", err))
		tflog.Debug(ctx, fmt.Sprintf("Full HTTP response: %v\n", r))
		return diag.FromErr(err)
	}
	// response from `V1SupportSettingsGet`: V1SupportSettingsGet200Response
	tflog.Debug(ctx, fmt.Sprintf("Response from `SupportApi.V1SupportSettingsGet`: %v\n", resp))

	errs.addMaybeError(d.Set("enabled", resp.GetEnabled()))
	errs.addMaybeError(d.Set("mq_host", resp.GetMqHost()))
	errs.addMaybeError(d.Set("mq_port", resp.GetMqPort()))
	errs.addMaybeError(d.Set("mq_proxy_host", resp.GetMqProxyHost()))
	errs.addMaybeError(d.Set("mq_proxy_port", resp.GetMqProxyPort()))
	errs.addMaybeError(d.Set("s3_proxy_host", resp.GetS3ProxyHost()))
	errs.addMaybeError(d.Set("s3_proxy_port", resp.GetS3ProxyPort()))
	errs.addMaybeError(d.Set("s3_proxy_disable_https", resp.GetS3ProxyDisableHttps()))
	errs.addMaybeError(d.Set("vpn_enabled", resp.GetVpnEnabled()))
	errs.addMaybeError(d.Set("vpn_host", resp.GetVpnHost()))
	errs.addMaybeError(d.Set("period", resp.GetPeriod()))

	return errs.diags
}

func resourceMonitoringUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*openapi.APIClient)

	ms := setMonitoringSettings(ctx, d, m)
	resp, r, err := c.SupportApi.V1SupportSettingsPatch(context.Background()).
		V1SupportSettingsGet200Response(*ms).Execute()
	if err != nil {
		tflog.Debug(ctx, fmt.Sprintf("Error when calling `SupportApi.V1SupportSettingsPatch``: %v\n", err))
		tflog.Debug(ctx, fmt.Sprintf("Full HTTP response: %v\n", r))
		return diag.FromErr(err)
	}
	// response from `V1SupportSettingsPatch`: V1AdJoinPost202Response
	tflog.Debug(ctx, fmt.Sprintf("Response from `SupportApi.V1SupportSettingsPatcg`: %v\n", resp))

	return resourceMonitoringRead(ctx, d, m)
}

func resourceMonitoringDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "Deleting monitor settings resource")

	return nil
}

func setMonitoringSettings(ctx context.Context, d *schema.ResourceData, m interface{}) *openapi.V1SupportSettingsGet200Response {

	// Integers are converted to float32 to match the openapi spec
	mc := openapi.NewV1SupportSettingsGet200Response()
	mc.SetEnabled(d.Get("enabled").(bool))
	mc.SetMqHost(d.Get("mq_host").(string))
	mc.SetMqPort(float32(d.Get("mq_port").(int)))
	mc.SetMqProxyHost(d.Get("mq_proxy_host").(string))
	mc.SetMqProxyPort(float32(d.Get("mq_proxy_port").(int)))
	mc.SetS3ProxyHost(d.Get("s3_proxy_host").(string))
	mc.SetS3ProxyPort(float32(d.Get("s3_proxy_port").(int)))
	mc.SetS3ProxyDisableHttps(d.Get("s3_proxy_disable_https").(bool))
	mc.SetVpnEnabled(d.Get("vpn_enabled").(bool))
	mc.SetVpnHost(d.Get("vpn_host").(string))
	mc.SetPeriod(float32(d.Get("period").(int)))

	tflog.Debug(ctx, "Updating monitor settings")

	return mc
}
