package qumulo

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strconv"
	"time"
)

const ClusterSettingsEndpoint = "/v1/cluster/settings"

type ClusterSettings struct {
	ClusterName string `json:"cluster_name"`
}

func resourceClusterSettings() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceClusterSettingsCreate,
		ReadContext:   resourceClusterSettingsRead,
		UpdateContext: resourceClusterSettingsUpdate,
		DeleteContext: resourceClusterSettingsDelete,
		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceClusterSettingsCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	name := d.Get("name").(string)

	cs := ClusterSettings{
		ClusterName: name,
	}

	_, err := DoRequest[ClusterSettings, ClusterSettings](c, PUT, ClusterSettingsEndpoint, &cs)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return resourceClusterSettingsRead(ctx, d, m)
}

func resourceClusterSettingsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	cs, err := DoRequest[ClusterSettings, ClusterSettings](c, GET, ClusterSettingsEndpoint, nil)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("name", cs.ClusterName); err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func resourceClusterSettingsUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	name := d.Get("name").(string)

	cs := ClusterSettings{
		ClusterName: name,
	}

	_, err := DoRequest[ClusterSettings, ClusterSettings](c, PUT, ClusterSettingsEndpoint, &cs)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceClusterSettingsRead(ctx, d, m)
}

func resourceClusterSettingsDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	return diags
}
