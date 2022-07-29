package qumulo

import (
	"context"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const SMBServerEndpoint = "/v1/smb/settings"

type SMBServerRequest struct {
	SessionEncryption      string   `json:"session_encryption"`
	SupportedDialects      []string `json:"supported_dialects"`
	HideSharesUsers        bool     `json:"hide_shares_from_unauthorized_users"`
	HideSharesHosts        bool     `json:"hide_shares_from_unauthorized_hosts"`
	SnapshotDirMode        string   `json:"snapshot_directory_mode"`
	BypassTraverseChecking bool     `json:"bypass_traverse_checking"`
	SigningRequired        bool     `json:"signing_required"`
}

var encryptionSettings = []string{"NONE", "PREFERRED", "REQUIRED"}
var validDialects = []string{"SMB2_DIALECT_2_002", "SMB2_DIALECT_2_1", "SMB2_DIALECT_3_0", "SMB2_DIALECT_3_11"}
var snapshotDirectoryMode = []string{"DISABLED", "HIDDEN", "VISIBLE"}

func resourceSMBServer() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSMBServerCreate,
		ReadContext:   resourceSMBServerRead,
		UpdateContext: resourceSMBServerUpdate,
		DeleteContext: resourceSMBServerDelete,
		Schema: map[string]*schema.Schema{
			"session_encryption": &schema.Schema{
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice(encryptionSettings, false)),
			},
			"supported_dialects": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice(validDialects, false)),
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
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice(snapshotDirectoryMode, false)),
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

func resourceSMBServerCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	err := setSMBServerSettings(ctx, d, m, PUT)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return resourceSMBServerRead(ctx, d, m)
}

func resourceSMBServerRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	var errs ErrorCollection
	SMBSettings, err := DoRequest[SMBServerRequest, SMBServerRequest](ctx, c, GET, SMBServerEndpoint, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	errs.addMaybeError(d.Set("session_encryption", SMBSettings.SessionEncryption))
	errs.addMaybeError(d.Set("supported_dialects", SMBSettings.SupportedDialects))
	errs.addMaybeError(d.Set("hide_shares_from_unauthorized_users", SMBSettings.HideSharesUsers))
	errs.addMaybeError(d.Set("hide_shares_from_unauthorized_hosts", SMBSettings.HideSharesHosts))
	errs.addMaybeError(d.Set("snapshot_directory_mode", SMBSettings.SnapshotDirMode))
	errs.addMaybeError(d.Set("bypass_traverse_checking", SMBSettings.BypassTraverseChecking))
	errs.addMaybeError(d.Set("signing_required", SMBSettings.SigningRequired))

	return errs.diags
}

func resourceSMBServerUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	err := setSMBServerSettings(ctx, d, m, PATCH)
	if err != nil {
		return diag.FromErr(err)
	}
	return resourceSMBServerRead(ctx, d, m)
}

func resourceSMBServerDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "Deleting SMB settings resource")
	var diags diag.Diagnostics

	return diags
}

func setSMBServerSettings(ctx context.Context, d *schema.ResourceData, m interface{}, method Method) error {
	c := m.(*Client)

	// convert the []interface{} into []string
	dials := d.Get("supported_dialects").([]interface{})
	dialects := make([]string, len(dials))
	for i, dial := range dials {
		dialects[i] = dial.(string)
	}

	SMBServerConfig := SMBServerRequest{
		SessionEncryption:      d.Get("session_encryption").(string),
		SupportedDialects:      dialects,
		HideSharesUsers:        d.Get("hide_shares_from_unauthorized_users").(bool),
		HideSharesHosts:        d.Get("hide_shares_from_unauthorized_hosts").(bool),
		SnapshotDirMode:        d.Get("snapshot_directory_mode").(string),
		BypassTraverseChecking: d.Get("bypass_traverse_checking").(bool),
		SigningRequired:        d.Get("signing_required").(bool),
	}

	tflog.Debug(ctx, "Updating SMB settings")
	_, err := DoRequest[SMBServerRequest, SMBServerRequest](ctx, c, method, SMBServerEndpoint, &SMBServerConfig)
	return err
}
