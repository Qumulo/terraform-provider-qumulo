package qumulo

import (
	"context"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const ClusterSettingsEndpoint = "/v1/cluster/settings"

type ClusterSettingsBody struct {
	ClusterName string `json:"cluster_name"`
}

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
	c := m.(*Client)

	cs, err := DoRequest[ClusterSettingsBody, ClusterSettingsBody](ctx, c, GET, ClusterSettingsEndpoint, nil)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("cluster_name", cs.ClusterName); err != nil {
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
	c := m.(*Client)

	clusterName := d.Get("cluster_name").(string)

	cs := ClusterSettingsBody{
		ClusterName: clusterName,
	}

	tflog.Debug(ctx, "Updating cluster settings")
	_, err := DoRequest[ClusterSettingsBody, ClusterSettingsBody](ctx, c, PUT, ClusterSettingsEndpoint, &cs)
	return err
}
