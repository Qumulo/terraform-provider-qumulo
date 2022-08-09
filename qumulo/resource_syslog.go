package qumulo

import (
	"context"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const SyslogConfigEndpoint = "/v1/audit/syslog/config"

type SyslogConfigBody struct {
	Enabled       bool   `json:"enabled"`
	ServerAddress string `json:"server_address"`
	ServerPort    int    `json:"server_port"`
}

func resourceSyslog() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSyslogCreate,
		ReadContext:   resourceSyslogRead,
		UpdateContext: resourceSyslogUpdate,
		DeleteContext: resourceSyslogDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"enabled": &schema.Schema{
				Type:     schema.TypeBool,
				Required: true,
			},
			"server_address": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"server_port": &schema.Schema{
				Type:             schema.TypeInt,
				Optional:         true,
				Default:          0,
				ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(0)),
			},
		},

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceSyslogCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	err := c.modifySyslogConfig(ctx, d, PUT)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return resourceSyslogRead(ctx, d, m)
}

func resourceSyslogRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	var errs ErrorCollection

	syslogConfig, err := DoRequest[SyslogConfigBody, SyslogConfigBody](ctx, c, GET, SyslogConfigEndpoint, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	errs.addMaybeError(d.Set("enabled", syslogConfig.Enabled))
	errs.addMaybeError(d.Set("server_address", syslogConfig.ServerAddress))
	errs.addMaybeError(d.Set("server_port", syslogConfig.ServerPort))

	return errs.diags
}

func resourceSyslogUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	err := c.modifySyslogConfig(ctx, d, PATCH)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceSyslogRead(ctx, d, m)
}

func resourceSyslogDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "Deleting audit log syslog resource")

	return nil
}

func (c *Client) modifySyslogConfig(ctx context.Context, d *schema.ResourceData, method Method) error {
	var config = SyslogConfigBody{
		Enabled:       d.Get("enabled").(bool),
		ServerAddress: d.Get("server_address").(string),
		ServerPort:    d.Get("server_port").(int),
	}

	tflog.Info(ctx, "Updating audit log syslog configuration")

	_, err := DoRequest[SyslogConfigBody, SyslogConfigBody](ctx, c, method, SyslogConfigEndpoint, &config)
	return err
}
