package qumulo

import (
	"context"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const SmbServerEndpoint = "/v1/smb/settings"

type SmbServerRequest struct {
	SessionEncryption               string   `json:"session_encryption"`
	SupportedDialects               []string `json:"supported_dialects"`
	HideSharesFromUnauthorizedUsers bool     `json:"hide_shares_from_unauthorized_users"`
	HideSharesFromUnauthorizedHosts bool     `json:"hide_shares_from_unauthorized_hosts"`
	SnapshotDirectoryMode           string   `json:"snapshot_directory_mode"`
	BypassTraverseChecking          bool     `json:"bypass_traverse_checking"`
	SigningRequired                 bool     `json:"signing_required"`
}

var SmbEncryptionSettings = []string{"NONE", "PREFERRED", "REQUIRED"}
var SmbValidDialects = []string{"SMB2_DIALECT_2_002", "SMB2_DIALECT_2_1", "SMB2_DIALECT_3_0", "SMB2_DIALECT_3_11",
	"API_SMB2_DIALECT_2_002", "API_SMB2_DIALECT_2_1", "API_SMB2_DIALECT_3_0",
	"API_SMB2_DIALECT_3_11"}
var SmbSnapshotDirectoryMode = []string{"DISABLED", "HIDDEN", "VISIBLE"}

func resourceSmbServer() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSmbServerCreate,
		ReadContext:   resourceSmbServerRead,
		UpdateContext: resourceSmbServerUpdate,
		DeleteContext: resourceSmbServerDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"session_encryption": &schema.Schema{
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice(SmbEncryptionSettings, false)),
			},
			"supported_dialects": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice(SmbValidDialects, false)),
				},
			},
			"hide_shares_from_unauthorized_users": &schema.Schema{
				Type:     schema.TypeBool,
				Required: true,
			},
			"hide_shares_from_unauthorized_hosts": &schema.Schema{
				Type:     schema.TypeBool,
				Required: true,
			},
			"snapshot_directory_mode": &schema.Schema{
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice(SmbSnapshotDirectoryMode, false)),
			},
			"bypass_traverse_checking": &schema.Schema{
				Type:     schema.TypeBool,
				Required: true,
			},
			"signing_required": &schema.Schema{
				Type:     schema.TypeBool,
				Required: true,
			},
		},
	}
}

func resourceSmbServerCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	// convert the []interface{} into []string
	dials := d.Get("supported_dialects").([]interface{})
	dialects := make([]string, len(dials))
	for i, dial := range dials {
		dialects[i] = dial.(string)
	}

	smbServerConfig := SmbServerRequest{
		SessionEncryption:               d.Get("session_encryption").(string),
		SupportedDialects:               dialects,
		HideSharesFromUnauthorizedUsers: d.Get("hide_shares_from_unauthorized_users").(bool),
		HideSharesFromUnauthorizedHosts: d.Get("hide_shares_from_unauthorized_hosts").(bool),
		SnapshotDirectoryMode:           d.Get("snapshot_directory_mode").(string),
		BypassTraverseChecking:          d.Get("bypass_traverse_checking").(bool),
		SigningRequired:                 d.Get("signing_required").(bool),
	}

	_, err := DoRequest[SmbServerRequest, SmbServerRequest](c, PUT, SmbServerEndpoint, &smbServerConfig)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return resourceSmbServerRead(ctx, d, m)
}

func resourceSmbServerRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	var errs ErrorCollection

	smbSettings, err := DoRequest[SmbServerRequest, SmbServerRequest](c, GET, SmbServerEndpoint, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	errs.addMaybeError(d.Set("session_encryption", smbSettings.SessionEncryption))
	errs.addMaybeError(d.Set("supported_dialects", smbSettings.SupportedDialects))
	errs.addMaybeError(d.Set("hide_shares_from_unauthorized_users", smbSettings.HideSharesFromUnauthorizedUsers))
	errs.addMaybeError(d.Set("hide_shares_from_unauthorized_hosts", smbSettings.HideSharesFromUnauthorizedHosts))
	errs.addMaybeError(d.Set("snapshot_directory_mode", smbSettings.SnapshotDirectoryMode))
	errs.addMaybeError(d.Set("bypass_traverse_checking", smbSettings.BypassTraverseChecking))
	errs.addMaybeError(d.Set("signing_required", smbSettings.SigningRequired))

	return errs.diags
}

func resourceSmbServerUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceSmbServerCreate(ctx, d, m)
}

func resourceSmbServerDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	return diags
}
