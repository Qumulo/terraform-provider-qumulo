package qumulo

import (
	"context"
	"fmt"
	"strconv"
	"terraform-provider-qumulo/openapi"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var SmbEncryptionSettings = []string{"NONE", "PREFERRED", "REQUIRED"}
var SmbValidDialects = []string{"SMB2_DIALECT_2_002", "SMB2_DIALECT_2_1", "SMB2_DIALECT_3_0", "SMB2_DIALECT_3_11"}
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
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceSmbServerCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*openapi.APIClient)

	settings := setSmbServerSettings(ctx, d, m)

	resp, r, err := c.SmbApi.V1SmbSettingsPut(context.Background()).
		V1SmbSettingsGet200Response(*settings).Execute()
	if err != nil {
		tflog.Debug(ctx, fmt.Sprintf("Error when calling `SmbApi.V1SmbSettingsPut``: %v\n", err))
		tflog.Debug(ctx, fmt.Sprintf("Full HTTP response: %v\n", r))
		return diag.FromErr(err)
	}
	// response from `V1SmbSettingsPut`: V1SmbSettingsGet200Response
	tflog.Debug(ctx, fmt.Sprintf("Response from `SmbApi.V1SmbSettingsPut`: %v\n", resp))

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return resourceSmbServerRead(ctx, d, m)
}

func resourceSmbServerRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*openapi.APIClient)

	var errs ErrorCollection

	resp, r, err := c.SmbApi.V1SmbSettingsGet(context.Background()).Execute()
	if err != nil {
		tflog.Debug(ctx, fmt.Sprintf("Error when calling `SmbApi.V1SmbSettingsGet``: %v\n", err))
		tflog.Debug(ctx, fmt.Sprintf("Full HTTP response: %v\n", r))
		return diag.FromErr(err)
	}
	// response from `V1SmbSettingsGet`: V1SmbSettingsGet200Response
	tflog.Debug(ctx, fmt.Sprintf("Response from `SmbApi.V1SmbSettingsGet`: %v\n", resp))

	errs.addMaybeError(d.Set("session_encryption", resp.GetSessionEncryption()))
	errs.addMaybeError(d.Set("supported_dialects", resp.GetSupportedDialects()))
	errs.addMaybeError(d.Set("hide_shares_from_unauthorized_users", resp.GetHideSharesFromUnauthorizedUsers()))
	errs.addMaybeError(d.Set("hide_shares_from_unauthorized_hosts", resp.GetHideSharesFromUnauthorizedHosts()))
	errs.addMaybeError(d.Set("snapshot_directory_mode", resp.GetSnapshotDirectoryMode()))
	errs.addMaybeError(d.Set("bypass_traverse_checking", resp.GetBypassTraverseChecking()))
	errs.addMaybeError(d.Set("signing_required", resp.GetSigningRequired()))

	return errs.diags
}

func resourceSmbServerUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*openapi.APIClient)

	settings := setSmbServerSettings(ctx, d, m)

	resp, r, err := c.SmbApi.V1SmbSettingsPatch(context.Background()).
		V1SmbSettingsGet200Response(*settings).Execute()
	if err != nil {
		tflog.Debug(ctx, fmt.Sprintf("Error when calling `SmbApi.V1SmbSettingsPatch``: %v\n", err))
		tflog.Debug(ctx, fmt.Sprintf("Full HTTP response: %v\n", r))
		return diag.FromErr(err)
	}
	// response from `V1SmbSettingsPatch`: V1SmbSettingsGet200Response
	tflog.Debug(ctx, fmt.Sprintf("Response from `SmbApi.V1SmbSettingsPatch`: %v\n", resp))

	return resourceSmbServerRead(ctx, d, m)
}

func resourceSmbServerDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "Deleting SMB settings resource")

	return nil
}

func setSmbServerSettings(ctx context.Context, d *schema.ResourceData, m interface{}) *openapi.V1SmbSettingsGet200Response {

	dialects := InterfaceSliceToStringSlice(d.Get("supported_dialects").([]interface{}))

	settings := openapi.NewV1SmbSettingsGet200Response()

	settings.SetSessionEncryption(d.Get("session_encryption").(string))
	settings.SetSupportedDialects(dialects)
	settings.SetHideSharesFromUnauthorizedUsers(d.Get("hide_shares_from_unauthorized_users").(bool))
	settings.SetHideSharesFromUnauthorizedHosts(d.Get("hide_shares_from_unauthorized_hosts").(bool))
	settings.SetSnapshotDirectoryMode(d.Get("snapshot_directory_mode").(string))
	settings.SetBypassTraverseChecking(d.Get("bypass_traverse_checking").(bool))
	settings.SetSigningRequired(d.Get("signing_required").(bool))

	tflog.Debug(ctx, "Updating SMB settings")

	return settings
}
