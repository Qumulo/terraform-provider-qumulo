package qumulo

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceName() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNameRead,
		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceNameRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	//client := &http.Cli ent{Timeout: 10 * time.Second}
	client := m.(*Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/v1/cluster/settings", "https://10.116.100.110:24100"), nil)
	if err != nil {
		return diag.FromErr(err)
	}

	r, err := client.HTTPClient.Do(req)
	if err != nil {
		return diag.FromErr(err)
	}
	defer r.Body.Close()

	Name := make([]map[string]interface{}, 0)
	err = json.NewDecoder(r.Body).Decode(&Name)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("name", Name); err != nil {
		return diag.FromErr(err)
	}

	// always run
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}
