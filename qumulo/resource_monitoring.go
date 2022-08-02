package qumulo

import (
	"context"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const MonitoringEndpoint = "/v1/support/settings"

type MonitoringSettings struct {
	Enabled             bool   `json:"enabled"`
	MqHost              string `json:"mq_host"`
	MqPort              int    `json:"mq_port"`
	MqProxyHost         string `json:"mq_proxy_host"`
	MqProxyPort         int    `json:"mq_proxy_port"`
	S3ProxyHost         string `json:"s3_proxy_host"`
	S3ProxyPort         int    `json:"s3_proxy_port"`
	S3ProxyDisableHttps bool   `json:"s3_proxy_disable_https"`
	VpnEnabled          bool   `json:"vpn_enabled"`
	VpnHost             string `json:"vpn_host"`
	Period              int    `json:"period"`
}

type MonitoringResponse struct {
	MonitorUri string `json:"monitor_uri"`
}

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
	err := setMonitoringSettings(ctx, d, m, PUT)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return resourceMonitoringRead(ctx, d, m)
}

func resourceMonitoringRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	var errs ErrorCollection

	settings, err := DoRequest[MonitoringSettings, MonitoringSettings](ctx, c, GET, MonitoringEndpoint, nil)
	if err != nil {
		return diag.FromErr(err)
	}
	errs.addMaybeError(d.Set("enabled", settings.Enabled))
	errs.addMaybeError(d.Set("mq_host", settings.MqHost))
	errs.addMaybeError(d.Set("mq_port", settings.MqPort))
	errs.addMaybeError(d.Set("mq_proxy_host", settings.MqProxyHost))
	errs.addMaybeError(d.Set("mq_proxy_port", settings.MqProxyPort))
	errs.addMaybeError(d.Set("s3_proxy_host", settings.S3ProxyHost))
	errs.addMaybeError(d.Set("s3_proxy_port", settings.S3ProxyPort))
	errs.addMaybeError(d.Set("s3_proxy_disable_https", settings.S3ProxyDisableHttps))
	errs.addMaybeError(d.Set("vpn_enabled", settings.VpnEnabled))
	errs.addMaybeError(d.Set("vpn_host", settings.VpnHost))
	errs.addMaybeError(d.Set("period", settings.Period))

	return errs.diags
}

func resourceMonitoringUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	err := setMonitoringSettings(ctx, d, m, PATCH)
	if err != nil {
		return diag.FromErr(err)
	}
	return resourceMonitoringRead(ctx, d, m)
}

func resourceMonitoringDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "Deleting monitor settings resource")
	var diags diag.Diagnostics

	return diags
}

func setMonitoringSettings(ctx context.Context, d *schema.ResourceData, m interface{}, method Method) error {
	c := m.(*Client)

	monitoringConfig := MonitoringSettings{
		Enabled:             d.Get("enabled").(bool),
		MqHost:              d.Get("mq_host").(string),
		MqPort:              d.Get("mq_port").(int),
		MqProxyHost:         d.Get("mq_proxy_host").(string),
		MqProxyPort:         d.Get("mq_proxy_port").(int),
		S3ProxyHost:         d.Get("s3_proxy_host").(string),
		S3ProxyPort:         d.Get("s3_proxy_port").(int),
		S3ProxyDisableHttps: d.Get("s3_proxy_disable_https").(bool),
		VpnEnabled:          d.Get("vpn_enabled").(bool),
		VpnHost:             d.Get("vpn_host").(string),
		Period:              d.Get("period").(int),
	}

	tflog.Debug(ctx, "Updating monitor settings")
	_, err := DoRequest[MonitoringSettings, MonitoringResponse](ctx, c, method, MonitoringEndpoint, &monitoringConfig)
	return err
}
