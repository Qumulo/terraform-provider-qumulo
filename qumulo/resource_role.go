package qumulo

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const RolesEndpoint = "/v1/auth/roles/"

type Role struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Privileges  []string `json:"privileges"`
}

func resourceRole() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRoleCreate,
		ReadContext:   resourceRoleRead,
		UpdateContext: resourceRoleUpdate,
		DeleteContext: resourceRoleDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"privileges": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceRoleCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	role := setRoleSettings(ctx, d, m)

	_, err := DoRequest[Role, Role](ctx, c, POST, RolesEndpoint, &role)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(d.Get("name").(string))

	return resourceRoleRead(ctx, d, m)
}

func resourceRoleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	var errs ErrorCollection

	readRoleByNameUri := RolesEndpoint + d.Id()

	role, err := DoRequest[Role, Role](ctx, c, GET, readRoleByNameUri, nil)
	if err != nil {
		return diag.FromErr(err)
	}
	errs.addMaybeError(d.Set("name", d.Id()))
	errs.addMaybeError(d.Set("description", role.Description))
	errs.addMaybeError(d.Set("privileges", role.Privileges))

	return errs.diags
}

func resourceRoleUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, fmt.Sprintf("Deleting role with name %q", d.Get("name").(string)))
	c := m.(*Client)

	role := setRoleSettings(ctx, d, m)
	updateRoleByNameUri := RolesEndpoint + d.Get("name").(string)

	_, err := DoRequest[Role, Role](ctx, c, PUT, updateRoleByNameUri, &role)

	if err != nil {
		return diag.FromErr(err)
	}
	return resourceRoleRead(ctx, d, m)
}

func resourceRoleDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	deleteRoleByNameUri := RolesEndpoint + d.Get("name").(string)

	_, err := DoRequest[Role, Role](ctx, c, DELETE, deleteRoleByNameUri, nil)

	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func setRoleSettings(ctx context.Context, d *schema.ResourceData, m interface{}) Role {
	var privileges = []string{}

	if v, ok := d.Get("privileges").([]interface{}); ok {
		privileges = InterfaceSliceToStringSlice(v)
	}

	roleConfig := Role{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Privileges:  privileges,
	}

	tflog.Debug(ctx, "Updating or creating Role: %v", map[string]interface{}{
		"Name": roleConfig.Name,
	})

	return roleConfig
}
