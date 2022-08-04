package qumulo

import (
	"context"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const TimeConfigurationEndpoint = "/v1/time/settings"

type TimeConfigurationBody struct {
	UseAdForPrimary bool     `json:"use_ad_for_primary"`
	NtpServers      []string `json:"ntp_servers"`
}

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
	err := setTimeConfiguration(ctx, d, m, PUT)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return resourceTimeConfigurationRead(ctx, d, m)
}

func resourceTimeConfigurationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	var errs ErrorCollection

	timeConfig, err := DoRequest[TimeConfigurationBody, TimeConfigurationBody](ctx, c, GET, TimeConfigurationEndpoint, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	errs.addMaybeError(d.Set("use_ad_for_primary", timeConfig.UseAdForPrimary))
	errs.addMaybeError(d.Set("ntp_servers", timeConfig.NtpServers))

	return errs.diags
}

func resourceTimeConfigurationUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	err := setTimeConfiguration(ctx, d, m, PATCH)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceTimeConfigurationRead(ctx, d, m)
}

func resourceTimeConfigurationDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	return diags
}

func setTimeConfiguration(ctx context.Context, d *schema.ResourceData, m interface{}, method Method) error {
	c := m.(*Client)

	// convert the []interface{} into []string
	inputNtpServers := d.Get("ntp_servers").([]interface{})
	ntpServers := make([]string, len(inputNtpServers))
	for i, ntpServer := range inputNtpServers {
		ntpServers[i] = ntpServer.(string)
	}

	timeConfigurationRequest := TimeConfigurationBody{
		UseAdForPrimary: d.Get("use_ad_for_primary").(bool),
		NtpServers:      ntpServers,
	}

	tflog.Debug(ctx, "Updating time configuration")
	_, err := DoRequest[TimeConfigurationBody, TimeConfigurationBody](ctx, c, method, TimeConfigurationEndpoint, &timeConfigurationRequest)

	return err
}
