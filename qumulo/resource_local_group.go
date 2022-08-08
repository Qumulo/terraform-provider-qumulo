package qumulo

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const GroupsEndpoint = "/v1/groups/"

type CreateGroupRequest struct {
	Name string `json:"name"`
	Gid  string `json:"gid"`
}

type UpdateGroupRequest struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Gid  string `json:"gid"`
}

type GroupResponse struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Sid  string `json:"sid"`
	Gid  string `json:"gid"`
}

type DeleteGroupBody struct{}

func resourceGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGroupCreate,
		ReadContext:   resourceGroupRead,
		UpdateContext: resourceGroupUpdate,
		DeleteContext: resourceGroupDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"sid": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"gid": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
		},

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceGroupCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	groupSettings := CreateGroupRequest{
		Name: d.Get("name").(string),
		Gid:  d.Get("gid").(string),
	}

	tflog.Debug(ctx, fmt.Sprintf("Creating local group with name %q", groupSettings.Name))

	group, err := DoRequest[CreateGroupRequest, GroupResponse](ctx, c, POST, GroupsEndpoint, &groupSettings)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(group.Id)

	return resourceGroupRead(ctx, d, m)
}

func resourceGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	var errs ErrorCollection

	readGroupByIdUri := GroupsEndpoint + d.Id()

	group, err := DoRequest[GroupResponse, GroupResponse](ctx, c, GET, readGroupByIdUri, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	errs.addMaybeError(d.Set("name", group.Name))
	errs.addMaybeError(d.Set("sid", group.Sid))
	errs.addMaybeError(d.Set("gid", group.Gid))

	return errs.diags
}

func resourceGroupUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	groupSettings := UpdateGroupRequest{
		Id:   d.Id(),
		Name: d.Get("name").(string),
		Gid:  d.Get("gid").(string),
	}

	updateGroupByIdUri := GroupsEndpoint + d.Id()

	tflog.Debug(ctx, fmt.Sprintf("Updating local group with name %q", groupSettings.Name))

	_, err := DoRequest[UpdateGroupRequest, GroupResponse](ctx, c, PUT, updateGroupByIdUri, &groupSettings)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceGroupRead(ctx, d, m)
}

func resourceGroupDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	deleteGroupByIdUri := GroupsEndpoint + d.Id()

	tflog.Debug(ctx, fmt.Sprintf("Deleting local group with id %q", d.Id()))

	_, err := DoRequest[DeleteGroupBody, DeleteGroupBody](ctx, c, DELETE, deleteGroupByIdUri, nil)

	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
