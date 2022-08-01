package qumulo

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const SmbSharesEndpoint = "/v2/smb/shares/"

var SmbPermissionTypes = []string{"ALLOWED", "DENIED"}
var SmbRights = []string{"READ", "WRITE", "CHANGE_PERMISSIONS"}

type SmbShare struct {
	Id                         string                 `json:"id"`
	ShareName                  string                 `json:"share_name"`
	FsPath                     string                 `json:"fs_path"`
	Description                string                 `json:"description"`
	Permissions                []SmbPermission        `json:"permissions"`
	NetworkPermissions         []SmbNetworkPermission `json:"network_permissions"`
	AccessBasedEnumEnabled     bool                   `json:"access_based_enumeration_enabled"`
	DefaultFileCreateMode      string                 `json:"default_file_create_mode,omitempty"`
	DefaultDirectoryCreateMode string                 `json:"default_directory_create_mode,omitempty"`
	BytesPerSector             string                 `json:"bytes_per_sector,omitempty"`
	RequireEncryption          bool                   `json:"require_encryption"`
}

type SmbPermission struct {
	Type    string     `json:"type"`
	Trustee SmbTrustee `json:"trustee"`
	Rights  []string   `json:"rights"`
}

type SmbNetworkPermission struct {
	Type          string   `json:"type"`
	AddressRanges []string `json:"address_ranges"`
	Rights        []string `json:"rights"`
}

type SmbTrustee struct {
	Domain string `json:"domain"`
	AuthId string `json:"auth_id"`
	Uid    string `json:"uid,omitempty"`
	Gid    string `json:"gid,omitempty"`
	Sid    string `json:"sid"`
	Name   string `json:"name,omitempty"`
}

func resourceSmbShare() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSmbShareCreate,
		ReadContext:   resourceSmbShareRead,
		UpdateContext: resourceSmbShareUpdate,
		DeleteContext: resourceSmbShareDelete,
		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"share_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"fs_path": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"permissions": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": &schema.Schema{
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice(SmbPermissionTypes, false)),
						},
						"trustee": &schema.Schema{
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 1,
							MinItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"domain": {
										Type:     schema.TypeString,
										Required: true,
									},
									"auth_id": {
										Type:     schema.TypeString,
										Required: true,
									},
									"uid": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"gid": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"sid": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"name": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"rights": &schema.Schema{
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Schema{
								Type:             schema.TypeString,
								ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice(SmbRights, false)),
							},
						},
					},
				},
			},
			"network_permissions": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": &schema.Schema{
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice(SmbPermissionTypes, false)),
						},
						"address_ranges": &schema.Schema{
							Type:     schema.TypeList,
							Optional: true,
							Elem:     schema.TypeString,
						},
						"rights": &schema.Schema{
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Schema{
								Type:             schema.TypeString,
								ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice(SmbRights, false)),
							},
						},
					},
				},
			},
			"access_based_enumeration_enabled": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"default_file_create_mode": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"default_directory_create_mode": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"bytes_per_sector": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"require_encryption": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"allow_fs_path_create": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func resourceSmbShareCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	smbShare := setSmbShare(d)
	createSmbSharetUri := SmbSharesEndpoint
	if v, ok := d.Get("allow_fs_path_create").(bool); ok {
		createSmbSharetUri = SmbSharesEndpoint + "?allow-fs-path-create=" + strconv.FormatBool(v)
	}

	res, err := DoRequest[SmbShare, SmbShare](ctx, client, POST, createSmbSharetUri, &smbShare)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(res.Id)

	return resourceSmbShareRead(ctx, d, m)

}

func resourceSmbShareRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	var errs ErrorCollection
	smbShareId := d.Id()
	getSmbShareByIdUri := SmbSharesEndpoint + smbShareId
	smbShare, err := DoRequest[SmbShare, SmbShare](ctx, client, GET, getSmbShareByIdUri, nil)

	if err != nil {
		return diag.FromErr(err)
	}
	errs.addMaybeError(d.Set("id", smbShareId))
	errs.addMaybeError(d.Set("share_name", smbShare.ShareName))
	errs.addMaybeError(d.Set("fs_path", smbShare.FsPath))
	errs.addMaybeError(d.Set("description", smbShare.Description))
	errs.addMaybeError(d.Set("permissions", flattenSmbPermissions(smbShare.Permissions)))
	errs.addMaybeError(d.Set("network_permissions", flattenSmbNetworkPermissions(smbShare.NetworkPermissions)))
	errs.addMaybeError(d.Set("access_based_enumeration_enabled", smbShare.AccessBasedEnumEnabled))
	errs.addMaybeError(d.Set("default_file_create_mode", smbShare.DefaultFileCreateMode))
	errs.addMaybeError(d.Set("default_directory_create_mode", smbShare.DefaultDirectoryCreateMode))
	errs.addMaybeError(d.Set("bytes_per_sector", smbShare.BytesPerSector))
	errs.addMaybeError(d.Set("require_encryption", smbShare.RequireEncryption))

	return errs.diags
}

func resourceSmbShareUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	smbShare := setSmbShare(d)
	smbShare.Id = d.Get("id").(string)

	smbShareId := d.Id()
	updateSmbShareByIdUri := SmbSharesEndpoint + smbShareId

	if v, ok := d.Get("allow_fs_path_create").(bool); ok {
		updateSmbShareByIdUri = updateSmbShareByIdUri + "?allow-fs-path-create=" + strconv.FormatBool(v)
	}

	_, err := DoRequest[SmbShare, SmbShare](ctx, client, PATCH, updateSmbShareByIdUri, &smbShare)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceSmbShareRead(ctx, d, m)
}

func resourceSmbShareDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	var diags diag.Diagnostics
	smbShareId := d.Id()
	_, err := DoRequest[string, SmbShare](ctx, client, DELETE, SmbSharesEndpoint, &smbShareId)
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func setSmbShare(d *schema.ResourceData) SmbShare {
	share := SmbShare{
		ShareName:                  d.Get("share_name").(string),
		FsPath:                     d.Get("fs_path").(string),
		Description:                d.Get("description").(string),
		Permissions:                expandPermissions(d.Get("permissions").([]interface{})),
		NetworkPermissions:         expandNetworkPermissions(d.Get("network_permissions").([]interface{})),
		AccessBasedEnumEnabled:     d.Get("access_based_enumeration_enabled").(bool),
		DefaultFileCreateMode:      d.Get("default_file_create_mode").(string),
		DefaultDirectoryCreateMode: d.Get("default_directory_create_mode").(string),
		BytesPerSector:             d.Get("bytes_per_sector").(string),
		RequireEncryption:          d.Get("require_encryption").(bool),
	}
	return share
}

func expandPermissions(tfPermissions []interface{}) []SmbPermission {
	var permissions []SmbPermission

	if len(tfPermissions) == 0 {
		return permissions
	}
	for _, tfPermission := range tfPermissions {
		tfMap, ok := tfPermission.(map[string]interface{})
		permission := SmbPermission{}
		if !ok {
			continue
		}

		if v, ok := tfMap["type"].(string); ok {
			permission.Type = v
		}

		if v, ok := tfMap["trustee"].([]interface{}); ok {
			permission.Trustee = expandTrustee(v[0])
		}

		if v, ok := tfMap["rights"].([]interface{}); ok {
			expandedRights := make([]string, len(v))
			for i, right := range v {
				expandedRights[i] = right.(string)
			}
			permission.Rights = expandedRights
		}

		permissions = append(permissions, permission)
	}

	return permissions
}

func expandNetworkPermissions(tfNetworkPermissions []interface{}) []SmbNetworkPermission {
	var networkPermissions []SmbNetworkPermission

	if len(tfNetworkPermissions) == 0 {
		return networkPermissions
	}
	for _, tfNetworkPermission := range tfNetworkPermissions {
		tfMap, ok := tfNetworkPermission.(map[string]interface{})
		networkPermission := SmbNetworkPermission{}
		if !ok {
			continue
		}

		if v, ok := tfMap["type"].(string); ok {
			networkPermission.Type = v
		}
		if v, ok := tfMap["address_ranges"].([]interface{}); ok {
			expandedAddressRanges := make([]string, len(v))
			for i, addressRange := range v {
				expandedAddressRanges[i] = addressRange.(string)
			}
			networkPermission.AddressRanges = expandedAddressRanges
		}
		if v, ok := tfMap["rights"].([]interface{}); ok {
			expandedRights := make([]string, len(v))
			for i, right := range v {
				expandedRights[i] = right.(string)
			}
			networkPermission.Rights = expandedRights
		}

		networkPermissions = append(networkPermissions, networkPermission)
	}

	return networkPermissions
}

func expandTrustee(tfTrustee interface{}) SmbTrustee {
	tfMap, _ := tfTrustee.(map[string]interface{})

	trustee := SmbTrustee{}

	if v, ok := tfMap["domain"].(string); ok {
		trustee.Domain = v
	}
	if v, ok := tfMap["auth_id"].(string); ok {
		trustee.AuthId = v
	}
	if v, ok := tfMap["uid"].(string); ok {
		trustee.Uid = v
	}
	if v, ok := tfMap["gid"].(string); ok {
		trustee.Gid = v
	}
	if v, ok := tfMap["sid"].(string); ok {
		trustee.Sid = v
	}
	if v, ok := tfMap["name"].(string); ok {
		trustee.Name = v
	}

	return trustee
}

func flattenSmbPermissions(permissions []SmbPermission) []interface{} {
	var tfList []interface{}

	for _, permission := range permissions {
		tfMap := map[string]interface{}{}

		tfMap["type"] = permission.Type

		if v := permission.Rights; len(v) != 0 {
			tfMap["rights"] = v
		}

		trusteeMapList := []map[string]interface{}{}
		trusteeMap := map[string]interface{}{}
		trustee := permission.Trustee

		trusteeMap["domain"] = trustee.Domain
		trusteeMap["auth_id"] = trustee.AuthId
		trusteeMap["uid"] = trustee.Uid
		trusteeMap["gid"] = trustee.Gid
		trusteeMap["sid"] = trustee.Sid
		trusteeMap["name"] = trustee.Name

		trusteeMapList = append(trusteeMapList, trusteeMap)

		tfMap["trustee"] = trusteeMapList

		tfList = append(tfList, tfMap)
	}
	return tfList
}

func flattenSmbNetworkPermissions(permissions []SmbNetworkPermission) []interface{} {
	var tfList []interface{}

	for _, networkPermission := range permissions {
		tfMap := map[string]interface{}{}

		tfMap["type"] = networkPermission.Type

		if v := networkPermission.AddressRanges; len(v) != 0 {
			tfMap["address_ranges"] = v
		}

		if v := networkPermission.Rights; len(v) != 0 {
			tfMap["rights"] = v
		}

		tfList = append(tfList, tfMap)
	}
	return tfList
}
