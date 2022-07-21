package qumulo

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

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
	Settings     *ActiveDirectorySettings
	JoinSettings *ActiveDirectoryJoinRequest
}

type ActiveDirectoryResponse struct {
	Settings     *ActiveDirectorySettings
	JoinResponse *ActiveDirectoryJoinResponse
}

const ADSettingsEndpoint = "/v1/ad/settings"
const ADJoinEndpoint = "/v1/ad/join"
const ADReconfigureEndpoint = "/v1/ad/reconfigure"
const ADLeaveEndpoint = "/v1/ad/leave"

var adSigningValues = []string{"NO_SIGNING", "WANT_SIGNING", "REQUIRE_SIGNING"}
var adSealingValues = []string{"NO_SEALING", "WANT_SEALING", "REQUIRE_SEALING"}
var adCryptoValues = []string{"NO_AES", "WANT_AES", "REQUIRE_AES"}

func resourceActiveDirectory() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceActiveDirectoryCreate,
		ReadContext:   resourceActiveDirectoryRead,
		UpdateContext: resourceActiveDirectoryUpdate,
		DeleteContext: resourceActiveDirectoryDelete,
		Schema: map[string]*schema.Schema{
			"domain": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"domain_netbios": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"ad_username": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"ad_password": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"ou": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"use_ad_posix_attributes": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"base_dn": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"signing": &schema.Schema{
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice(adSigningValues, false)),
			},
			"sealing": &schema.Schema{
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice(adSealingValues, false)),
			},
			"crypto": &schema.Schema{
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice(adCryptoValues, false)),
			},
		},
	}
}

func resourceActiveDirectoryCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	updatedAdSettings := ActiveDirectorySettings{
		Signing: d.Get("signing").(string),
		Sealing: d.Get("sealing").(string),
		Crypto:  d.Get("crypto").(string),
	}

	joinSettings := ActiveDirectoryJoinRequest{
		Domain:               d.Get("domain").(string),
		Domain_NetBIOS:       d.Get("domain_netbios").(string),
		User:                 d.Get("ad_username").(string),
		Password:             d.Get("ad_password").(string),
		OU:                   d.Get("ou").(string),
		UseADPosixAttributes: d.Get("use_ad_posix_attributes").(bool),
		BaseDN:               d.Get("base_dn").(string),
	}

	updatedAdRequest := ActiveDirectoryRequest{
		Settings:     &updatedAdSettings,
		JoinSettings: &joinSettings,
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

	updatedAdSettings := ActiveDirectorySettings{
		Signing: d.Get("signing").(string),
		Sealing: d.Get("sealing").(string),
		Crypto:  d.Get("crypto").(string),
	}

	updatedAdRequest := ActiveDirectoryRequest{
		Settings: &updatedAdSettings,
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

	settingsResponsePointer, err := c.UpdateActiveDirectorySettings(clusterReq.Settings)
	if err != nil {
		return nil, err
	}

	joinResponsePointer, err := c.UpdateActiveDirectoryStatus(clusterReq.JoinSettings)
	if err != nil {
		return nil, err
	}

	response := ActiveDirectoryResponse{
		Settings:     settingsResponsePointer,
		JoinResponse: joinResponsePointer,
	}

	return &response, nil
}

func (c *Client) UpdateActiveDirectorySettings(activeDirectorySettings *ActiveDirectorySettings) (*ActiveDirectorySettings, error) {
	// The AD settings API endpoint expects all of the AD settings set.
	// If the config has all settings set, use them.
	// If the config has no settings set, don't hit the endpoint.
	// If the config has SOME settings set, return an error sicne we can't apply that.

	// Get a count of all of the AD settings which are set.
	// (We have front-end validation on proper types; the field is empty if it was absent in the Terraform file.)
	adSettingsCount := 0
	if activeDirectorySettings.Signing != "" {
		adSettingsCount++
	}
	if activeDirectorySettings.Sealing != "" {
		adSettingsCount++
	}
	if activeDirectorySettings.Crypto != "" {
		adSettingsCount++
	}

	// Perform the appropriate action based on the number of fields set.
	if adSettingsCount == 3 {
		settingsResponse, err := DoRequest[ActiveDirectorySettings, ActiveDirectorySettings](c, PUT, ADSettingsEndpoint, activeDirectorySettings)
		if err != nil {
			return nil, err
		}
		return settingsResponse, nil
	} else if adSettingsCount == 1 || adSettingsCount == 2 {
		// TODO: decide if this should return an error
		log.Printf("[WARN] Incomplete Active Directory settings detected, will not apply changes. Specify all or none of Signing, Sealing, and Crypto.")
		return nil, nil
	} else {
		log.Printf("[DEBUG] No Active Directory settings detected, will not apply changes.")
		return nil, nil
	}
}

func (c *Client) UpdateActiveDirectoryStatus(joinRequest *ActiveDirectoryJoinRequest) (*ActiveDirectoryJoinResponse, error) {
	if joinRequest == nil {
		return nil, nil
	}

	joinResponse, err := DoRequest[ActiveDirectoryJoinRequest, ActiveDirectoryJoinResponse](c, POST, ADJoinEndpoint, joinRequest)
	if err != nil {
		return nil, err
	}
	return joinResponse, nil
}
