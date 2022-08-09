package qumulo

import (
	"context"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const CloudWatchConfigEndpoint = "/v1/audit/cloudwatch/config"

type CloudWatchConfigBody struct {
	Enabled      bool   `json:"enabled"`
	LogGroupName string `json:"log_group_name"`
	Region       string `json:"region"`
}

func resourceCloudWatch() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCloudWatchCreate,
		ReadContext:   resourceCloudWatchRead,
		UpdateContext: resourceCloudWatchUpdate,
		DeleteContext: resourceCloudWatchDelete,

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
			"log_group_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"region": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
		},

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceCloudWatchCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	err := c.modifyCloudWatchConfig(ctx, d, PUT)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return resourceCloudWatchRead(ctx, d, m)
}

func resourceCloudWatchRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	var errs ErrorCollection

	cloudWatchConfig, err := DoRequest[CloudWatchConfigBody, CloudWatchConfigBody](ctx, c, GET, CloudWatchConfigEndpoint, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	errs.addMaybeError(d.Set("enabled", cloudWatchConfig.Enabled))
	errs.addMaybeError(d.Set("log_group_name", cloudWatchConfig.LogGroupName))
	errs.addMaybeError(d.Set("region", cloudWatchConfig.Region))

	return errs.diags
}

func resourceCloudWatchUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	err := c.modifyCloudWatchConfig(ctx, d, PATCH)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceCloudWatchRead(ctx, d, m)
}

func resourceCloudWatchDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "Deleting CloudWatch resource")

	return nil
}

func (c *Client) modifyCloudWatchConfig(ctx context.Context, d *schema.ResourceData, method Method) error {
	var config = CloudWatchConfigBody{
		Enabled:      d.Get("enabled").(bool),
		LogGroupName: d.Get("log_group_name").(string),
		Region:       d.Get("region").(string),
	}

	tflog.Info(ctx, "Updating audit log CloudWatch configuration")

	_, err := DoRequest[CloudWatchConfigBody, CloudWatchConfigBody](ctx, c, method, CloudWatchConfigEndpoint, &config)
	return err
}
