package qumulo

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const SMBSharesEndpoint = "/v2/smb/shares/"

var permissionTypes = []string{"ALLOWED", "DENIED"}
var rights = []string{"READ", "WRITE", "CHANGE_PERMISSIONS"}

// var userMappings = []string{"NFS_MAP_NONE", "NFS_MAP_ALL", "NFS_MAP_ROOT"}

// var fields = []string{"FILE_IDS", "FILE_SIZES", "FS_SIZE", "ALL"}

type SMBShare struct {
	Id                     string        `json:"id"`
	ExportPath             string        `json:"export_path"`
	FsPath                 string        `json:"fs_path"`
	Description            string        `json:"description"`
	Restrictions           []Permission  `json:"restrictions"`
	FieldsToPresentAs32Bit []interface{} `json:"fields_to_present_as_32_bit"`
}

type Permission struct {
	HostRestrictions      []string               `json:"host_restrictions"`
	ReadOnly              bool                   `json:"read_only"`
	RequirePrivilegedPort bool                   `json:"require_privileged_port"`
	UserMapping           string                 `json:"user_mapping"`
	MapToUser             map[string]interface{} `json:"map_to_user,omitempty"`
	MapToGroup            map[string]interface{} `json:"map_to_group,omitempty"`
}

func resourceSmbShare() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNfsExportCreate,
		ReadContext:   resourceNfsExportRead,
		UpdateContext: resourceNfsExportUpdate,
		DeleteContext: resourceNfsExportDelete,
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
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice(permissionTypes, false)),
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
										Type:     schema.TypeInt,
										Required: true,
									},
									"uid": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"gid": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"sid": {
										Type:     schema.TypeString,
										Required: true,
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
								ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice(rights, false)),
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
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice(permissionTypes, false)),
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
								ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice(rights, false)),
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
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(512, 512)),
			},
			"require_encryption": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func resourceNfsExportCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	nfsExport := NFSExport{
		ExportPath:             d.Get("export_path").(string),
		FsPath:                 d.Get("fs_path").(string),
		Description:            d.Get("description").(string),
		Restrictions:           expandRestrictions(d.Get("restrictions").([]interface{})),
		FieldsToPresentAs32Bit: d.Get("fields_to_present_as_32_bit").([]interface{}),
	}
	createNfsExportUri := NFSExportEndpoint
	if v, ok := d.Get("allow_fs_path_create").(bool); ok {
		createNfsExportUri = createNfsExportUri + "?allow-fs-path-create=" + strconv.FormatBool(v)
	}

	res, err := DoRequest[NFSExport, NFSExport](client, POST, createNfsExportUri, &nfsExport)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(res.Id)

	return resourceNfsExportRead(ctx, d, m)

}

func resourceNfsExportRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	var errs ErrorCollection
	nfsExportId := d.Id()
	getNfsExportByIdUri := NFSExportEndpoint + nfsExportId
	nfsExport, err := DoRequest[NFSExport, NFSExport](client, GET, getNfsExportByIdUri, nil)

	if err != nil {
		return diag.FromErr(err)
	}
	errs.addMaybeError(d.Set("id", nfsExportId))
	errs.addMaybeError(d.Set("export_path", nfsExport.ExportPath))
	errs.addMaybeError(d.Set("fs_path", nfsExport.FsPath))
	errs.addMaybeError(d.Set("description", nfsExport.Description))
	errs.addMaybeError(d.Set("restrictions", flattenNfsRestrictions(nfsExport.Restrictions)))
	errs.addMaybeError(d.Set("fields_to_present_as_32_bit", nfsExport.FieldsToPresentAs32Bit))

	return errs.diags
}

func resourceNfsExportUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	nfsExport := NFSExport{
		Id:                     d.Get("id").(string),
		ExportPath:             d.Get("export_path").(string),
		FsPath:                 d.Get("fs_path").(string),
		Description:            d.Get("description").(string),
		Restrictions:           expandRestrictions(d.Get("restrictions").([]interface{})),
		FieldsToPresentAs32Bit: d.Get("fields_to_present_as_32_bit").([]interface{}),
	}
	nfsExportId := d.Id()
	updateNfsExportByIdUri := NFSExportEndpoint + nfsExportId

	if v, ok := d.Get("allow_fs_path_create").(bool); ok {
		updateNfsExportByIdUri = updateNfsExportByIdUri + "?allow-fs-path-create=" + strconv.FormatBool(v)
	}

	_, err := DoRequest[NFSExport, NFSExport](client, PATCH, updateNfsExportByIdUri, &nfsExport)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceNfsExportRead(ctx, d, m)
}

func resourceNfsExportDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	var diags diag.Diagnostics
	nfsExportId := d.Id()
	_, err := DoRequest[string, NFSExport](client, DELETE, NFSExportEndpoint, &nfsExportId)
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func expandRestrictions(tfRestrictions []interface{}) []Restriction {
	var restrictions []Restriction

	if len(tfRestrictions) == 0 {
		return restrictions
	}
	for _, tfRestriction := range tfRestrictions {
		tfMap, ok := tfRestriction.(map[string]interface{})
		restriction := Restriction{}
		if !ok {
			continue
		}

		if v, ok := tfMap["host_restrictions"].([]interface{}); ok {
			expandedHostRestrictions := make([]string, len(v))
			for i, hostRestriction := range v {
				expandedHostRestrictions[i] = hostRestriction.(string)
			}
			restriction.HostRestrictions = expandedHostRestrictions
		}
		if v, ok := tfMap["read_only"].(bool); ok {
			restriction.ReadOnly = v
		}
		if v, ok := tfMap["require_privileged_port"].(bool); ok {
			restriction.RequirePrivilegedPort = v
		}
		if v, ok := tfMap["user_mapping"].(string); ok {
			restriction.UserMapping = v
		}
		if v, ok := tfMap["map_to_user"].(map[string]interface{}); ok {
			restriction.MapToUser = v
		}
		if v, ok := tfMap["map_to_group"].(map[string]interface{}); ok {
			restriction.MapToGroup = v
		}
		restrictions = append(restrictions, restriction)
	}
	return restrictions
}

func flattenNfsRestrictions(restrictions []Restriction) []interface{} {
	var tfList []interface{}

	for _, restriction := range restrictions {
		tfMap := map[string]interface{}{}

		if v := restriction.HostRestrictions; len(v) != 0 {
			tfMap["host_restrictions"] = v
		}
		tfMap["read_only"] = restriction.ReadOnly
		tfMap["require_privileged_port"] = restriction.RequirePrivilegedPort
		tfMap["user_mapping"] = restriction.UserMapping
		tfMap["map_to_user"] = restriction.MapToUser
		tfMap["map_to_group"] = restriction.MapToGroup
		tfList = append(tfList, tfMap)
	}
	return tfList
}
