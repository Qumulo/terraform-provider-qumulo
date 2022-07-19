package qumulo

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// TODO convert to enum to follow API specification
type ActiveDirectorySettings struct {
	Signing string `json:"signing"`
	Sealing string `json:"sealing"`
	Crypto  string `json:"crypto"`
}

type ActiveDirectoryRequest struct {
	Settings ActiveDirectorySettings
}

type ActiveDirectoryResponse struct {
	Name string `json:"test"`
}

func resourceActiveDirectory() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceActiveDirectoryCreate,
		ReadContext:   resourceActiveDirectoryRead,
		UpdateContext: resourceActiveDirectoryUpdate,
		DeleteContext: resourceActiveDirectoryDelete,
		Schema: map[string]*schema.Schema{
			"signing": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"sealing": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"crypto": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceActiveDirectoryCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	adSigning := d.Get("signing").(string)
	adSealing := d.Get("sealing").(string)
	adCrypto := d.Get("crypto").(string)

	// TODO verify value is valid for enum

	updatedAdSettings := ActiveDirectorySettings{
		Signing: adSigning,
		Sealing: adSealing,
		Crypto:  adCrypto,
	}

	_, err := client.UpdateActiveDirectory(updatedAdSettings)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}

func resourceActiveDirectoryRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/v1/ad/settings", client.HostURL), nil)
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		return diag.FromErr(err)
	}

	r, err := client.doRequest(req)
	if err != nil {
		return diag.FromErr(err)
	}

	cr := ActiveDirectorySettings{}
	err = json.Unmarshal(r, &cr)

	// TODO make Go-idiomatic
	if err := d.Set("signing", cr.Signing); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("sealing", cr.Sealing); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("crypto", cr.Crypto); err != nil {
		return diag.FromErr(err)
	}

	// always run
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}

func resourceActiveDirectoryUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// c := m.(*Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	// name := d.Get("name").(string)

	// updatedCluster := ClusterRequest{
	// 	// ClusterActiveDirectory: name,
	// }
	// _, err := c.UpdateActiveDirectory(updatedCluster)
	// if err != nil {
	// 	return diag.FromErr(err)
	// }

	// d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}

func resourceActiveDirectoryDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	return diags
}

func (c *Client) UpdateActiveDirectory(clusterReq ActiveDirectorySettings) (*ActiveDirectorySettings, error) {
	bearerToken := "Bearer " + c.Bearer_Token

	HostURL := c.HostURL

	rb, err := json.Marshal(clusterReq)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/v1/ad/settings", HostURL), strings.NewReader(string(rb)))
	req.Header.Set("Authorization", bearerToken)
	req.Header.Add("Content-Type", "application/json")
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	cr := ActiveDirectorySettings{}
	err = json.Unmarshal(body, &cr)
	if err != nil {
		return nil, err
	}

	return &cr, nil
}
