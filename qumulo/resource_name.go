package qumulo

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func resourceName() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNameCreate,
		ReadContext:   resourceNameRead,
		UpdateContext: resourceNameUpdate,
		DeleteContext: resourceNameDelete,
		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceNameCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	name := d.Get("name").(string)

	updatedCluster := ClusterRequest{
		ClusterName: name,
	}

	_, err := c.UpdateName(updatedCluster)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}

func resourceNameRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	return diags
}

func resourceNameUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	name := d.Get("name").(string)

	updatedCluster := ClusterRequest{
		ClusterName: name,
	}
	_, err := c.UpdateName(updatedCluster)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}

func resourceNameDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	return diags
}

type ClusterRequest struct {
	ClusterName string `json:"cluster_name"`
}

func (c *Client) UpdateName(clusterReq ClusterRequest) (*ClusterResponse, error) {
	bearerToken := "Bearer " + c.Bearer_Token

	HostURL := c.HostURL

	rb, err := json.Marshal(clusterReq)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/v1/cluster/settings", HostURL), strings.NewReader(string(rb)))
	req.Header.Set("Authorization", bearerToken)
	req.Header.Add("Content-Type", "application/json")
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	cr := ClusterResponse{}
	err = json.Unmarshal(body, &cr)
	if err != nil {
		return nil, err
	}

	return &cr, nil
}
