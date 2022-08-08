package qumulo

import (
	"context"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const FileSystemPermissionsEndpoint = "/v1/file-system/settings/permissions"
const FileSystemAtimeEndpoint = "/v1/file-system/settings/atime"

type FileSystemPermissionsModeSetting int
type FileSystemAtimeGranularitySetting int

const (
	CrossProtocol FileSystemPermissionsModeSetting = iota + 1
	Native
)

const (
	Day FileSystemAtimeGranularitySetting = iota + 1
	Hour
	Week
)

func (e FileSystemPermissionsModeSetting) String() string {
	return FileSystemPermissionsModeValues[e-1]
}

func (e FileSystemAtimeGranularitySetting) String() string {
	return FileSystemAtimeGranularityValues[e-1]
}

var FileSystemPermissionsModeValues = []string{"CROSS_PROTOCOL", "NATIVE"}
var FileSystemAtimeGranularityValues = []string{"DAY", "HOUR", "WEEK"}

type FileSystemPermissionsSettingsBody struct {
	Mode string `json:"mode"`
}

type FileSystemAtimeSettingsBody struct {
	Enabled     bool   `json:"enabled"`
	Granularity string `json:"granularity"`
}

type FileSystemSettingsBody struct {
	Permissions   *FileSystemPermissionsSettingsBody
	AtimeSettings *FileSystemAtimeSettingsBody
}

func resourceFileSystemSettings() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFileSystemSettingsCreate,
		ReadContext:   resourceFileSystemSettingsRead,
		UpdateContext: resourceFileSystemSettingsUpdate,
		DeleteContext: resourceFileSystemSettingsDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"permissions_mode": &schema.Schema{
				Type:             schema.TypeString,
				Optional:         true,
				Default:          CrossProtocol.String(),
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice(FileSystemPermissionsModeValues, false)),
			},
			"atime_enabled": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"atime_granularity": &schema.Schema{
				Type:             schema.TypeString,
				Optional:         true,
				Default:          Hour.String(),
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice(FileSystemAtimeGranularityValues, false)),
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceFileSystemSettingsCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	err := c.createOrUpdateFileSystemSettings(ctx, d, PUT)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return resourceFileSystemSettingsRead(ctx, d, m)
}

func resourceFileSystemSettingsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)
	var errs ErrorCollection

	fileSystemPermissions, err := DoRequest[FileSystemPermissionsSettingsBody, FileSystemPermissionsSettingsBody](ctx,
		c, GET, FileSystemPermissionsEndpoint, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	errs.addMaybeError(d.Set("permissions_mode", fileSystemPermissions.Mode))

	fileSystemAtime, err := DoRequest[FileSystemAtimeSettingsBody, FileSystemAtimeSettingsBody](ctx, c, GET,
		FileSystemAtimeEndpoint, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	errs.addMaybeError(d.Set("atime_enabled", fileSystemAtime.Enabled))
	errs.addMaybeError(d.Set("atime_granularity", fileSystemAtime.Granularity))

	return errs.diags
}

func resourceFileSystemSettingsUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	err := c.createOrUpdateFileSystemSettings(ctx, d, PATCH)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceFileSystemSettingsRead(ctx, d, m)
}

func resourceFileSystemSettingsDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "Deleting file system settomgs resource")

	return nil
}

func (c *Client) createOrUpdateFileSystemSettings(ctx context.Context, d *schema.ResourceData, atimeMethod Method) error {

	permissionsSettings := FileSystemPermissionsSettingsBody{
		Mode: d.Get("permissions_mode").(string),
	}

	atimeSettings := FileSystemAtimeSettingsBody{
		Enabled:     d.Get("atime_enabled").(bool),
		Granularity: d.Get("atime_granularity").(string),
	}

	var err error

	if d.HasChanges("permissions_mode") {
		tflog.Info(ctx, "Updating file system permission settings")

		_, err = DoRequest[FileSystemPermissionsSettingsBody, FileSystemPermissionsSettingsBody](ctx,
			c, PUT, FileSystemPermissionsEndpoint, &permissionsSettings)
		if err != nil {
			return err
		}
	}

	if d.HasChanges("atime_enabled", "atime_granularity") {
		tflog.Info(ctx, "Updating file system atime settings")

		_, err = DoRequest[FileSystemAtimeSettingsBody, FileSystemAtimeSettingsBody](ctx, c, atimeMethod,
			FileSystemAtimeEndpoint, &atimeSettings)
		if err != nil {
			return err
		}
	}

	return nil
}
