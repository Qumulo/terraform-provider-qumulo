package qumulo

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/go-cty/cty"
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

type SMBServerResponse struct {
	Placeholder string `json:"placeholder"`
}

var encryptionSettings = []string{"NONE", "PREFERRED", "REQUIRED"}
var validDialects = []string{"SMB2_DIALECT_002", "SMB2_DIALECT_2_1", "SMB2_DIALECT_3_0", "SMB2_DIALECT_3_11",
	"API_SMB2_DIALECT_2_002", "API_SMB2_DIALECT_2_1", "API_SMB2_DIALECT_3_0",
	"API_SMB2_DIALECT_3_11"}
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
				Type:             schema.TypeSet,
				Required:         true,
				ValidateDiagFunc: validateDialects,
				Elem: &schema.Schema{
					Type: schema.TypeString,
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

func validateDialects(v interface{}, p cty.Path) diag.Diagnostics {
	// Returns a set with only the items in the first set but not the second
	var diags diag.Diagnostics
	value := v.(*schema.Set).List()

	for _, val := range value {
		// Check if val is in validDialects (loop is fine, only 8 members, unlikely to scale)
		valid := false
		for _, dialect := range validDialects {
			if val == dialect {
				valid = true
				break
			}
		}
		if !valid {
			d := diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "invalid dialect",
				Detail:   fmt.Sprintf("%q is not a valid supported dialect", val),
			}
			diags = append(diags, d)
		}
	}
	return diags
}

func resourceSMBServerCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	var dialects []string
	//dialects := d.Get("supported_dialects")
	for _, dial := range d.Get("supported_dialects").(*schema.Set).List() {
		dialects = append(dialects, dial.(string))
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

	_, err := DoRequest[SMBServerRequest, SMBServerResponse](c, PUT, SMBServerEndpoint, &SMBServerConfig)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}

func resourceSMBServerRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	return diags
}

func resourceSMBServerUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceSMBServerCreate(ctx, d, m)
}

func resourceSMBServerDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	return diags
}
