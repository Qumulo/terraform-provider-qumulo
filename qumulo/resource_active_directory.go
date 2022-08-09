package qumulo

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

type ActiveDirectorySigning int
type ActiveDirectorySealing int
type ActiveDirectoryCrypto int

const (
	NoSigning ActiveDirectorySigning = iota + 1
	WantSigning
	RequireSigning
)

const (
	NoSealing ActiveDirectorySealing = iota + 1
	WantSealing
	RequireSealing
)

const (
	NoCrypto ActiveDirectoryCrypto = iota + 1
	WantCrypto
	RequireCrypto
)

func (e ActiveDirectorySigning) String() string {
	return ActiveDirectorySigningValues[e-1]
}

func (e ActiveDirectorySealing) String() string {
	return ActiveDirectorySealingValues[e-1]
}

func (e ActiveDirectoryCrypto) String() string {
	return ActiveDirectoryCryptoValues[e-1]
}

var ActiveDirectorySigningValues = []string{"NO_SIGNING", "WANT_SIGNING", "REQUIRE_SIGNING"}
var ActiveDirectorySealingValues = []string{"NO_SEALING", "WANT_SEALING", "REQUIRE_SEALING"}
var ActiveDirectoryCryptoValues = []string{"NO_AES", "WANT_AES", "REQUIRE_AES"}

type ActiveDirectorySettingsBody struct {
	Signing string `json:"signing"`
	Sealing string `json:"sealing"`
	Crypto  string `json:"crypto"`
}

type ActiveDirectoryJoinRequest struct {
	Domain               string `json:"domain"`
	DomainNetBios        string `json:"domain_netbios"`
	User                 string `json:"user"`
	Password             string `json:"password"`
	Ou                   string `json:"ou"`
	UseAdPosixAttributes bool   `json:"use_ad_posix_attributes"`
	BaseDn               string `json:"base_dn"`
}

type ActiveDirectoryJoinResponse struct {
	MonitorUri string `json:"monitor_uri"`
}

type ActiveDirectoryStatusBody struct {
	Status               string                      `json:"status"`
	Domain               string                      `json:"domain"`
	Ou                   string                      `json:"ou"`
	UseAdPosixAttributes bool                        `json:"use_ad_posix_attributes"`
	BaseDn               string                      `json:"base_dn"`
	DomainNetBios        string                      `json:"domain_netbios"`
	Dcs                  []ActiveDirectoryDcs        `json:"dcs"`
	LdapConnectionStates []ActiveDirectoryLdapStates `json:"ldap_connection_states"`
}

type ActiveDirectoryDcs struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

type ActiveDirectoryLdapServers struct {
	BindUri    string `json:"bind_uri"`
	KdcAddress string `json:"kdc_address"`
}

type ActiveDirectoryLdapStates struct {
	NodeId      int                          `json:"node_id"`
	Servers     []ActiveDirectoryLdapServers `json:"servers"`
	BindDomain  string                       `json:"bind_domain"`
	BindAccount string                       `json:"bind_account"`
	BaseDnVec   []string                     `json:"base_dn_vec"`
	Health      string                       `json:"health"`
}

type ActiveDirectoryRequest struct {
	Settings      *ActiveDirectorySettingsBody
	JoinSettings  *ActiveDirectoryJoinRequest
	UsageSettings *ActiveDirectoryUsageSettingsRequest
}

type ActiveDirectoryResponse struct {
	Settings     *ActiveDirectorySettingsBody
	JoinResponse *ActiveDirectoryJoinResponse
}

type ActiveDirectoryMonitorLastError struct {
	Module      string `json"module"`
	ErrorClass  string `json:"error_class"`
	Description string `json:"description"`
	Stack       string `json:"stack"`
	UserVisible bool   `json:"user_visible"`
}

func (e ActiveDirectoryMonitorLastError) Error() string {
	return fmt.Sprintf("Error %s encountered in Active Directory\nDescription: %s\nStack: %s", e.ErrorClass, e.Description, e.Stack)
}

type ActiveDirectoryMonitorResponse struct {
	Status               string                          `json:"status"`
	Domain               string                          `json:"domain"`
	Ou                   string                          `json:"ou"`
	LastError            ActiveDirectoryMonitorLastError `json:"last_error"`
	LastActionTime       string                          `json:"last_action_time"`
	UseAdPosixAttributes bool                            `json:"use_ad_posix_attributes"`
	BaseDn               string                          `json:"base_dn"`
	DomainNetBios        string                          `json:"domain_netbios"`
}

type ActiveDirectoryUsageSettingsRequest struct {
	UseAdPosixAttributes bool   `json:"use_ad_posix_attributes"`
	BaseDn               string `json:"base_dn"`
}

type ActiveDirectoryLeaveRequest struct {
	Domain   string `json:"domain"`
	User     string `json:"user"`
	Password string `json:"password"`
}

const AdSettingsEndpoint = "/v1/ad/settings"
const AdStatusEndpoint = "/v1/ad/status"
const AdJoinEndpoint = "/v1/ad/join"
const AdMonitorEndpoint = "/v1/ad/monitor"
const AdReconfigureEndpoint = "/v1/ad/reconfigure"
const AdLeaveEndpoint = "/v1/ad/leave"

const AdJoinWaitTime = 1 * time.Second
const AdJoinTimeoutIterations = 60

func resourceActiveDirectory() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceActiveDirectoryCreate,
		ReadContext:   resourceActiveDirectoryRead,
		UpdateContext: resourceActiveDirectoryUpdate,
		DeleteContext: resourceActiveDirectoryDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(2 * time.Minute),
			Update: schema.DefaultTimeout(2 * time.Minute),
			Delete: schema.DefaultTimeout(2 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"domain": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"domain_netbios": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"ad_username": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"ad_password": &schema.Schema{
				Type:      schema.TypeString,
				Required:  true,
				ForceNew:  true,
				Sensitive: true,
			},
			"ou": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
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
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice(ActiveDirectorySigningValues, false)),
				Default:          WantSigning.String(),
			},
			"sealing": &schema.Schema{
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice(ActiveDirectorySealingValues, false)),
				Default:          WantSealing.String(),
			},
			"crypto": &schema.Schema{
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice(ActiveDirectoryCryptoValues, false)),
				Default:          WantCrypto.String(),
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceActiveDirectoryCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	adSettings := ActiveDirectorySettingsBody{
		Signing: d.Get("signing").(string),
		Sealing: d.Get("sealing").(string),
		Crypto:  d.Get("crypto").(string),
	}

	joinRequest := ActiveDirectoryJoinRequest{
		Domain:               d.Get("domain").(string),
		DomainNetBios:        d.Get("domain_netbios").(string),
		User:                 d.Get("ad_username").(string),
		Password:             d.Get("ad_password").(string),
		Ou:                   d.Get("ou").(string),
		UseAdPosixAttributes: d.Get("use_ad_posix_attributes").(bool),
		BaseDn:               d.Get("base_dn").(string),
	}

	adRequest := ActiveDirectoryRequest{
		Settings:     &adSettings,
		JoinSettings: &joinRequest,
	}

	tflog.Debug(ctx, "Joining Active Directory")
	_, err := createActiveDirectory(ctx, c, adRequest)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return resourceActiveDirectoryRead(ctx, d, m)
}

func resourceActiveDirectoryRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	adSettings, err := DoRequest[ActiveDirectorySettingsBody, ActiveDirectorySettingsBody](ctx, c, GET, AdSettingsEndpoint, nil)

	if err != nil {
		return diag.FromErr(err)
	}

	var errs ErrorCollection

	errs.addMaybeError(d.Set("signing", adSettings.Signing))
	errs.addMaybeError(d.Set("sealing", adSettings.Sealing))
	errs.addMaybeError(d.Set("crypto", adSettings.Crypto))

	if errs.diags != nil {
		return errs.diags
	}

	adStatus, err := DoRequest[ActiveDirectoryStatusBody, ActiveDirectoryStatusBody](ctx, c, GET, AdStatusEndpoint, nil)

	if err != nil {
		return diag.FromErr(err)
	}

	tflog.Debug(ctx, "Active Directory status:", map[string]interface{}{
		"adStatus": adStatus.Status,
	})
	errs.addMaybeError(d.Set("domain", adStatus.Domain))
	errs.addMaybeError(d.Set("ou", adStatus.Ou))
	errs.addMaybeError(d.Set("use_ad_posix_attributes", adStatus.UseAdPosixAttributes))
	errs.addMaybeError(d.Set("base_dn", adStatus.BaseDn))
	errs.addMaybeError(d.Set("domain_netbios", adStatus.DomainNetBios))
	return errs.diags
}

func resourceActiveDirectoryUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	updatedAdSettings := ActiveDirectorySettingsBody{
		Signing: d.Get("signing").(string),
		Sealing: d.Get("sealing").(string),
		Crypto:  d.Get("crypto").(string),
	}

	updatedUsageSettings := ActiveDirectoryUsageSettingsRequest{
		UseAdPosixAttributes: d.Get("use_ad_posix_attributes").(bool),
		BaseDn:               d.Get("base_dn").(string),
	}

	updatedAdRequest := ActiveDirectoryRequest{
		Settings:      &updatedAdSettings,
		UsageSettings: &updatedUsageSettings,
	}

	tflog.Debug(ctx, "Updating Active Directory settings")
	_, err := updateActiveDirectory(ctx, c, updatedAdRequest, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceActiveDirectoryRead(ctx, d, m)
}

func resourceActiveDirectoryDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	leaveAdSettings := ActiveDirectoryLeaveRequest{
		Domain:   d.Get("domain").(string),
		User:     d.Get("ad_username").(string),
		Password: d.Get("ad_password").(string),
	}

	_, err := DoRequest[ActiveDirectoryLeaveRequest, ActiveDirectoryMonitorResponse](ctx, c, POST, AdLeaveEndpoint, &leaveAdSettings)
	if err != nil {
		return diag.FromErr(err)
	}

	tflog.Info(ctx, "Leaving Active Directory")
	err = waitForADMonitorUpdate(ctx, c)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func createActiveDirectory(ctx context.Context, c *Client, clusterReq ActiveDirectoryRequest) (*ActiveDirectoryResponse, error) {

	joinResponsePointer, err := joinActiveDirectory(ctx, c, clusterReq.JoinSettings)
	if err != nil {
		return nil, err
	}

	settingsResponsePointer, err := updateActiveDirectorySettings(ctx, c, clusterReq.Settings)
	if err != nil {
		return nil, err
	}

	response := ActiveDirectoryResponse{
		Settings:     settingsResponsePointer,
		JoinResponse: joinResponsePointer,
	}

	return &response, nil
}

func updateActiveDirectory(ctx context.Context, c *Client, clusterReq ActiveDirectoryRequest, d *schema.ResourceData) (*ActiveDirectoryResponse, error) {

	var joinResponsePointer *ActiveDirectoryJoinResponse
	var err error

	if d.HasChanges("use_ad_posix_attributes", "base_dn") {
		joinResponsePointer, err = updateActiveDirectoryUsage(ctx, c, clusterReq.UsageSettings)
		if err != nil {
			return nil, err
		}
	}

	var settingsResponsePointer *ActiveDirectorySettingsBody

	if d.HasChanges("signing", "sealing", "crypto") {
		settingsResponsePointer, err = updateActiveDirectorySettings(ctx, c, clusterReq.Settings)
		if err != nil {
			return nil, err
		}
	}

	response := ActiveDirectoryResponse{
		Settings:     settingsResponsePointer,
		JoinResponse: joinResponsePointer,
	}

	return &response, nil
}

func updateActiveDirectorySettings(ctx context.Context, c *Client, activeDirectorySettings *ActiveDirectorySettingsBody) (*ActiveDirectorySettingsBody, error) {
	// XXX amanning32: The AD settings API endpoint expects all of the AD settings set.
	// If the config has all settings set, use them.
	// If the config has no settings set, don't hit the endpoint.
	// If the config has SOME settings set, return an error since we can't apply that.

	// (We have front-end validation on proper types; the field is empty if it was absent in the Terraform file.)
	if activeDirectorySettings.Signing == "" && activeDirectorySettings.Sealing == "" && activeDirectorySettings.Crypto == "" {
		tflog.Debug(ctx, "No Active Directory settings detected, will not apply changes.")
		return nil, nil
	} else if activeDirectorySettings.Signing == "" || activeDirectorySettings.Sealing == "" || activeDirectorySettings.Crypto == "" {
		// TODO: decide if this should return an error
		tflog.Warn(ctx, "Incomplete Active Directory settings detected, will not apply changes. Specify all or none of Signing, Sealing, and Crypto.")
		return nil, nil
	} else {
		settingsResponse, err := DoRequest[ActiveDirectorySettingsBody, ActiveDirectorySettingsBody](ctx, c, PUT, AdSettingsEndpoint, activeDirectorySettings)
		if err != nil {
			return nil, err
		}
		return settingsResponse, nil
	}
}

func joinActiveDirectory(ctx context.Context, c *Client, joinRequest *ActiveDirectoryJoinRequest) (*ActiveDirectoryJoinResponse, error) {
	if joinRequest == nil {
		tflog.Warn(ctx, "No Active Directory join information detected, not joining.")
		return nil, nil
	}

	joinResponse, err := DoRequest[ActiveDirectoryJoinRequest, ActiveDirectoryJoinResponse](ctx, c, POST, AdJoinEndpoint, joinRequest)
	if err != nil {
		return nil, err
	}

	err = waitForADMonitorUpdate(ctx, c)
	if err != nil {
		return nil, err
	}

	return joinResponse, nil
}

func updateActiveDirectoryUsage(ctx context.Context, c *Client, usageRequest *ActiveDirectoryUsageSettingsRequest) (*ActiveDirectoryJoinResponse, error) {
	if usageRequest == nil {
		tflog.Debug(ctx, " No updated Active Directory usage settings detected, will not apply changes.")
		return nil, nil
	}

	usageUpdateResponse, err := DoRequest[ActiveDirectoryUsageSettingsRequest, ActiveDirectoryJoinResponse](ctx, c, POST, AdReconfigureEndpoint, usageRequest)
	if err != nil {
		return nil, err
	}

	err = waitForADMonitorUpdate(ctx, c)
	if err != nil {
		return nil, err
	}

	return usageUpdateResponse, nil
}

func waitForADMonitorUpdate(ctx context.Context, c *Client) error {

	var finishedJoinStatus *ActiveDirectoryMonitorResponse

	joinCompleted := false
	numIterations := 0

	for !joinCompleted {
		joinStatus, err := DoRequest[ActiveDirectoryMonitorResponse, ActiveDirectoryMonitorResponse](ctx, c, GET, AdMonitorEndpoint, nil)
		if err != nil {
			return err
		}

		if !strings.Contains(joinStatus.Status, "IN_PROGRESS") {
			joinCompleted = true
			finishedJoinStatus = joinStatus
		} else {
			tflog.Debug(ctx, "Waiting another second for AD operation to complete.")
			numIterations++
			time.Sleep(AdJoinWaitTime)
		}

		// XXX amanning32: remove since the resource itself has a timeout?
		if numIterations > AdJoinTimeoutIterations {
			tflog.Error(ctx, "Active Directory operation timed out, exiting")
			return fmt.Errorf("ERROR: Active Directory operation timed out after %d seconds, aborting", AdJoinTimeoutIterations)
		}
	}

	if strings.Contains(finishedJoinStatus.Status, "FAILED") {
		return finishedJoinStatus.LastError
	}

	tflog.Debug(ctx, "AD join status: ", map[string]interface{}{
		"status": finishedJoinStatus.Status,
	})

	return nil
}
