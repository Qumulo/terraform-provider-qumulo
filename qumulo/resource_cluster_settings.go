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

func resourceClusterSettings() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceClusterSettingsCreate,
		ReadContext:   resourceClusterSettingsRead,
		UpdateContext: resourceClusterSettingsUpdate,
		DeleteContext: resourceClusterSettingsDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"cluster_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceClusterSettingsCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	err := setClusterSettings(ctx, d, m)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	return resourceClusterSettingsRead(ctx, d, m)
}

func resourceClusterSettingsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*openapi.APIClient)

	resp, r, err := c.ClusterApi.V1ClusterSettingsGet(context.Background()).Execute()
	if err != nil {
		tflog.Debug(ctx, fmt.Sprintf("Error when calling `ClusterApi.V1ClusterSettingsGet``: %v\n", err))
		tflog.Debug(ctx, fmt.Sprintf("Full HTTP response: %v\n", r))
		return diag.FromErr(err)
	}

	// response from `V1ClusterSettingsGet`: V1ClusterSettingsGet200Response
	tflog.Debug(ctx, fmt.Sprintf("Response from `ClusterApi.V1ClusterSettingsGet`: %v\n",
		resp.GetClusterName()))

	if err := d.Set("cluster_name", resp.GetClusterName()); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceClusterSettingsUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	err := setClusterSettings(ctx, d, m)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceClusterSettingsRead(ctx, d, m)
}

func resourceClusterSettingsDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "Deleting cluster settings resource")

	return nil
}

func setClusterSettings(ctx context.Context, d *schema.ResourceData, m interface{}) error {
	c := m.(*openapi.APIClient)

	clusterName := d.Get("cluster_name").(string)

	cs := openapi.NewV1ClusterSettingsGet200Response()
	cs.SetClusterName(clusterName)

	tflog.Debug(ctx, "Updating cluster settings")

	resp, r, err := c.ClusterApi.V1ClusterSettingsPut(context.Background()).
		V1ClusterSettingsGet200Response(*cs).Execute()
	if err != nil {
		tflog.Debug(ctx, fmt.Sprintf("Error when calling `ClusterApi.V1ClusterSettingsPut``: %v\n", err))
		tflog.Debug(ctx, fmt.Sprintf("Full HTTP response: %v\n", r))
	}
	// response from `V1ClusterSettingsPut`: V1ClusterSettingsGet200Response
	tflog.Debug(ctx, fmt.Sprintf("Response from `ClusterApi.V1ClusterSettingsPut`: %v\n", resp))

	return err
}
