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

const NfsExportsEndpoint = "/v2/nfs/exports/"

var NfsExportsUserMappingsValues = []string{"NFS_MAP_NONE", "NFS_MAP_ALL", "NFS_MAP_ROOT"}

var NfsExportsFieldToPresentAs32BitValues = []string{"FILE_IDS", "FILE_SIZES", "FS_SIZE", "ALL"}

type NfsExport struct {
	Id                     string           `json:"id"`
	ExportPath             string           `json:"export_path"`
	FsPath                 string           `json:"fs_path"`
	Description            string           `json:"description"`
	Restrictions           []NfsRestriction `json:"restrictions"`
	FieldsToPresentAs32Bit []interface{}    `json:"fields_to_present_as_32_bit"`
}

type NfsRestriction struct {
	HostRestrictions      []string               `json:"host_restrictions"`
	ReadOnly              bool                   `json:"read_only"`
	RequirePrivilegedPort bool                   `json:"require_privileged_port"`
	UserMapping           string                 `json:"user_mapping"`
	MapToUser             map[string]interface{} `json:"map_to_user,omitempty"`
	MapToGroup            map[string]interface{} `json:"map_to_group,omitempty"`
}

func resourceNfsExport() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNfsExportCreate,
		ReadContext:   resourceNfsExportRead,
		UpdateContext: resourceNfsExportUpdate,
		DeleteContext: resourceNfsExportDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"export_path": &schema.Schema{
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
			"restrictions": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"host_restrictions": {
							Type:     schema.TypeList,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"read_only": {
							Type:     schema.TypeBool,
							Required: true,
						},
						"require_privileged_port": {
							Type:     schema.TypeBool,
							Required: true,
						},
						"user_mapping": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice(NfsExportsUserMappingsValues, false)),
						},
						"map_to_user": {
							Type:     schema.TypeMap,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"map_to_group": {
							Type:     schema.TypeMap,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"fields_to_present_as_32_bit": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice(NfsExportsFieldToPresentAs32BitValues, false)),
				},
			},
			"allow_fs_path_create": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func resourceNfsExportCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	res, err := createOrUpdateNfsExport(ctx, d, m, POST, NfsExportsEndpoint)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(res.Id)

	return resourceNfsExportRead(ctx, d, m)

}

func resourceNfsExportRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	var errs ErrorCollection

	nfsExportId := d.Id()
	getNfsExportByIdUri := NfsExportsEndpoint + nfsExportId
	nfsExport, err := DoRequest[NfsExport, NfsExport](ctx, c, GET, getNfsExportByIdUri, nil)

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
	nfsExportId := d.Id()
	updateNfsExportByIdUri := NfsExportsEndpoint + nfsExportId

	_, err := createOrUpdateNfsExport(ctx, d, m, PATCH, updateNfsExportByIdUri)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceNfsExportRead(ctx, d, m)
}

func resourceNfsExportDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "Deleting NFS Export")
	c := m.(*Client)

	var diags diag.Diagnostics

	nfsExportId := d.Id()
	_, err := DoRequest[string, NfsExport](ctx, c, DELETE, NfsExportsEndpoint, &nfsExportId)
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func createOrUpdateNfsExport(ctx context.Context, d *schema.ResourceData, m interface{}, method Method, url string) (*NfsExport, error) {
	c := m.(*Client)

	nfsExport := NfsExport{
		Id:                     d.Get("id").(string),
		ExportPath:             d.Get("export_path").(string),
		FsPath:                 d.Get("fs_path").(string),
		Description:            d.Get("description").(string),
		Restrictions:           expandRestrictions(ctx, d.Get("restrictions").([]interface{})),
		FieldsToPresentAs32Bit: d.Get("fields_to_present_as_32_bit").([]interface{}),
	}

	if v, ok := d.Get("allow_fs_path_create").(bool); ok {
		url = url + "?allow-fs-path-create=" + strconv.FormatBool(v)
	}

	tflog.Debug(ctx, "Creating/Updating NFS Export")
	res, err := DoRequest[NfsExport, NfsExport](ctx, c, method, url, &nfsExport)
	return res, err
}

func expandRestrictions(ctx context.Context, tfRestrictions []interface{}) []NfsRestriction {
	var restrictions []NfsRestriction

	if len(tfRestrictions) == 0 {
		tflog.Warn(ctx, "No restrictions found for the NFS Export")
		return restrictions
	}
	for _, tfRestriction := range tfRestrictions {
		tfMap, ok := tfRestriction.(map[string]interface{})
		restriction := NfsRestriction{}
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

func flattenNfsRestrictions(restrictions []NfsRestriction) []interface{} {
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
