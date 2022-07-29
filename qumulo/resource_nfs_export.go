package qumulo

import (
	"context"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const NFSExportEndpoint = "/v2/nfs/exports/"

var userMappings = []string{"NFS_MAP_NONE", "NFS_MAP_ALL", "NFS_MAP_ROOT"}

var fields = []string{"FILE_IDS", "FILE_SIZES", "FS_SIZE", "ALL"}

type NFSExport struct {
	Id                     string        `json:"id"`
	ExportPath             string        `json:"export_path"`
	FsPath                 string        `json:"fs_path"`
	Description            string        `json:"description"`
	Restrictions           []Restriction `json:"restrictions"`
	FieldsToPresentAs32Bit []interface{} `json:"fields_to_present_as_32_bit"`
}

type Restriction struct {
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
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice(userMappings, false)),
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
					ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice(fields, false)),
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
	res, err := createOrUpdateNFSExport(ctx, d, m, POST, NFSExportEndpoint)
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
	nfsExport, err := DoRequest[NFSExport, NFSExport](ctx, client, GET, getNfsExportByIdUri, nil)

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
	updateNfsExportByIdUri := NFSExportEndpoint + nfsExportId

	_, err := createOrUpdateNFSExport(ctx, d, m, PATCH, updateNfsExportByIdUri)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceNfsExportRead(ctx, d, m)
}

func resourceNfsExportDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "Deleting NFS Export")
	client := m.(*Client)
	var diags diag.Diagnostics

	nfsExportId := d.Id()
	_, err := DoRequest[string, NFSExport](ctx, client, DELETE, NFSExportEndpoint, &nfsExportId)
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func createOrUpdateNFSExport(ctx context.Context, d *schema.ResourceData, m interface{}, method Method, url string) (*NFSExport, error) {
	client := m.(*Client)
	nfsExport := NFSExport{
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
	res, err := DoRequest[NFSExport, NFSExport](ctx, client, method, url, &nfsExport)
	return res, err
}

func expandRestrictions(ctx context.Context, tfRestrictions []interface{}) []Restriction {
	var restrictions []Restriction

	if len(tfRestrictions) == 0 {
		tflog.Warn(ctx, "No restrictions found for the NFS Export")
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
