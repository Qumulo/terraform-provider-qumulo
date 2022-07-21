package qumulo

import (
	"context"
	"strconv"
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

type ActiveDirectoryJoinRequest struct {
	Domain               string `json:"domain"`
	Domain_NetBIOS       string `json:"domain_netbios"`
	User                 string `json:"user"`
	Password             string `json:"password"`
	OU                   string `json:"ou"`
	UseADPosixAttributes bool   `json:"use_ad_posix_attributes"`
	BaseDN               string `json:"base_dn"`
}

type ActiveDirectoryJoinResponse struct {
	MonitorURI string `json:"monitor_uri"`
}

type ActiveDirectoryRequest struct {
	Settings     ActiveDirectorySettings
	JoinSettings ActiveDirectoryJoinRequest
}

type ActiveDirectoryResponse struct {
	Settings     ActiveDirectorySettings
	JoinResponse ActiveDirectoryJoinResponse
}

const ADSettingsEndpoint = "/v1/ad/settings"
const ADJoinEndpoint = "/v1/ad/join"

func resourceActiveDirectory() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceActiveDirectoryCreate,
		ReadContext:   resourceActiveDirectoryRead,
		UpdateContext: resourceActiveDirectoryUpdate,
		DeleteContext: resourceActiveDirectoryDelete,
		Schema: map[string]*schema.Schema{
			// "domain": &schema.Schema{
			// 	Type:     schema.TypeString,
			// 	Required: true,
			// },
			// "domain_netbios": &schema.Schema{
			// 	Type:     schema.TypeString,
			// 	Optional: true,
			// },
			// "ad_username": &schema.Schema{
			// 	Type:     schema.TypeString,
			// 	Required: true,
			// },
			// "ad_password": &schema.Schema{
			// 	Type:     schema.TypeString,
			// 	Required: true,
			// },
			// "ou": &schema.Schema{
			// 	Type:     schema.TypeString,
			// 	Optional: true,
			// },
			// "use_ad_posix_attributes": &schema.Schema{
			// 	Type:     schema.TypeBool,
			// 	Optional: true,
			// },
			// "base_dn": &schema.Schema{
			// 	Type:     schema.TypeString,
			// 	Optional: true,
			// },
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

	// TODO verify value is valid for enum

	updatedAdSettings := ActiveDirectorySettings{
		Signing: d.Get("signing").(string),
		Sealing: d.Get("sealing").(string),
		Crypto:  d.Get("crypto").(string),
	}

	// joinSettings := ActiveDirectoryJoinRequest{
	// 	Domain:               d.Get("domain").(string),
	// 	Domain_NetBIOS:       d.Get("domain_netbios").(string),
	// 	User:                 d.Get("user").(string),
	// 	Password:             d.Get("password").(string),
	// 	OU:                   d.Get("ou").(string),
	// 	UseADPosixAttributes: d.Get("use_ad_posix_attributes").(bool),
	// 	BaseDN:               d.Get("base_dn").(string),
	// }

	updatedAdRequest := ActiveDirectoryRequest{
		Settings: updatedAdSettings,
		// JoinSettings: joinSettings,
	}

	_, err := client.UpdateActiveDirectory(updatedAdRequest)
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

	cr, err := DoRequest[ActiveDirectorySettings, ActiveDirectorySettings](client, GET, ADSettingsEndpoint, nil)

	if err != nil {
		return diag.FromErr(err)
	}

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

	return diags
}

func resourceActiveDirectoryUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	// TODO verify value is valid for enum

	updatedAdSettings := ActiveDirectorySettings{
		Signing: d.Get("signing").(string),
		Sealing: d.Get("sealing").(string),
		Crypto:  d.Get("crypto").(string),
	}

	updatedAdRequest := ActiveDirectoryRequest{
		Settings: updatedAdSettings,
	}

	_, err := client.UpdateActiveDirectory(updatedAdRequest)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}

func resourceActiveDirectoryDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	return diags
}

func (c *Client) UpdateActiveDirectory(clusterReq ActiveDirectoryRequest) (*ActiveDirectoryResponse, error) {

	settingsResponse, err := DoRequest[ActiveDirectorySettings, ActiveDirectorySettings](c, PUT, ADSettingsEndpoint, &clusterReq.Settings)
	if err != nil {
		return nil, err
	}

	response := ActiveDirectoryResponse{
		Settings: *settingsResponse,
	}

	return &response, nil
}
