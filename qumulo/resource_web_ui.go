package qumulo

import (
	"context"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const WebUiEndpoint = "/v1/web-ui/settings"

type WebUiBody struct {
	InactivityTimeout WebUiTimeout `json:"inactivity_timeout"`
}

type WebUiEmpty struct {
	InactivityTimeout *string `json:"inactivity_timeout"`
}

type WebUiTimeout struct {
	Nanoseconds string `json:"nanoseconds"`
}

func resourceWebUi() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceWebUiCreate,
		ReadContext:   resourceWebUiRead,
		UpdateContext: resourceWebUiUpdate,
		DeleteContext: resourceWebUiDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"inactivity_timeout": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"nanoseconds": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		},

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceWebUiCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	err := setWebUi(ctx, d, m, PUT)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return resourceWebUiRead(ctx, d, m)
}

func resourceWebUiRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	var errs ErrorCollection

	uiConfig, err := DoRequest[WebUiBody, WebUiBody](ctx, c, GET, WebUiEndpoint, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	var tfList []interface{}

	tfMap := map[string]interface{}{}
	tfMap["nanoseconds"] = uiConfig.InactivityTimeout.Nanoseconds

	tfList = append(tfList, tfMap)

	errs.addMaybeError(d.Set("inactivity_timeout", tfList))

	return errs.diags
}

func resourceWebUiUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	err := setWebUi(ctx, d, m, PATCH)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceWebUiRead(ctx, d, m)
}

func resourceWebUiDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "Deleting web UI resource")

	c := m.(*Client)

	// Using a different struct here is a workaround. We have to pass in something that will be represented
	// as null in json, and a nil string pointer is the simplest way to do that.
	nullTimeout := WebUiEmpty{}
	nullTimeout.InactivityTimeout = nil

	_, err := DoRequest[WebUiEmpty, WebUiEmpty](ctx, c, PATCH, WebUiEndpoint, &nullTimeout)

	return diag.FromErr(err)
}

func setWebUi(ctx context.Context, d *schema.ResourceData, m interface{}, method Method) error {
	c := m.(*Client)

	tfMap, ok := d.Get("inactivity_timeout").([]interface{})[0].(map[string]interface{})
	if !ok {
		tflog.Debug(ctx, "Error getting web UI resource")
	}

	timeout := WebUiTimeout{}
	if v, ok := tfMap["nanoseconds"].(string); ok {
		timeout.Nanoseconds = v
	} else {
		tflog.Debug(ctx, "Error getting web UI timeout")
	}

	webUiRequest := WebUiBody{
		InactivityTimeout: timeout,
	}

	tflog.Debug(ctx, "Updating web UI")
	_, err := DoRequest[WebUiBody, WebUiBody](ctx, c, method, WebUiEndpoint, &webUiRequest)

	return err
}
