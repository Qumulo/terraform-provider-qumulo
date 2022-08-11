package qumulo

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// URL endpoints defined in separate files (RolesEndpoint and MembersSuffix)
const MembersEnding = "/members"

type RoleMemberAddRequest struct {
	Domain   string `json:"domain,omitempty"`
	AuthId   string `json:"auth_id,omitempty"`
	Uid      string `json:"uid,omitempty"`
	Gid      string `json:"gid,omitempty"`
	Sid      string `json:"sid,omitempty"`
	Name     string `json:"name,omitempty"`
	RoleName string `json:"role_name"`
}

// CREATE response, READ response
type RoleMemberResponse struct {
	Domain string      `json:"domain"`
	AuthId string      `json:"auth_id"`
	Uid    StringOrInt `json:"uid"`
	Gid    StringOrInt `json:"gid"`
	Sid    string      `json:"sid"`
	Name   string      `json:"name"`
}

type RoleMemberDeleteBody struct{}

func resourceRoleMember() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRoleMemberCreate,
		ReadContext:   resourceRoleMemberRead,
		DeleteContext: resourceRoleMemberDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"domain": &schema.Schema{
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice(RoleDomainValues, false)),
			},
			"auth_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"uid": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"gid": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"sid": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"role_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},

		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				// d.Id() here is the last argument passed to the
				// `terraform import RESOURCE_TYPE.RESOURCE_NAME RESOURCE_ID` command
				ids, err := ParseRoleMemberId(d.Id())
				if err != nil {
					return nil, err
				}

				// Now, ids is guaranteed to have length 2, and d.Id() is well formed
				d.Set("role_name", ids[0])
				d.Set("auth_id", ids[1])

				return []*schema.ResourceData{d}, nil
			},
		},
	}
}

func resourceRoleMemberCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	memberSettings := RoleMemberAddRequest{
		Domain:   d.Get("domain").(string),
		AuthId:   d.Get("auth_id").(string),
		Uid:      d.Get("uid").(string),
		Gid:      d.Get("gid").(string),
		Sid:      d.Get("sid").(string),
		Name:     d.Get("name").(string),
		RoleName: d.Get("role_name").(string),
	}

	tflog.Info(ctx, fmt.Sprintf("Adding member to role with name %q; member info: %q",
		memberSettings.RoleName, memberSettings))

	addMemberToRoleUri := RolesEndpoint + memberSettings.RoleName + MembersEnding

	tflog.Debug(ctx, fmt.Sprintf("Adding member with URL %s", addMemberToRoleUri))

	joinResponse, err := DoRequest[RoleMemberAddRequest, RoleMemberResponse](ctx, c, POST, addMemberToRoleUri, &memberSettings)
	if err != nil {
		return diag.FromErr(err)
	}

	// The format of a member resource ID is of the form {role_name}:{auth_id}
	ids := make([]string, 2)
	ids[0] = memberSettings.RoleName
	ids[1] = joinResponse.AuthId
	id, err := FormRoleMemberId(ids)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	return resourceGroupMemberRead(ctx, d, m)
}

func resourceRoleMemberRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	ids, err := ParseRoleMemberId(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	readRoleMemberUri := RolesEndpoint + ids[0] + MembersSuffix + ids[1]

	tflog.Debug(ctx, fmt.Sprintf("Reading member with URL %s", readRoleMemberUri))

	readResponse, err := DoRequest[RoleMemberResponse, RoleMemberResponse](ctx, c, GET, readRoleMemberUri, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	var errs ErrorCollection

	errs.addMaybeError(d.Set("domain", readResponse.Domain))
	errs.addMaybeError(d.Set("auth_id", readResponse.AuthId))
	errs.addMaybeError(SetStringOrIntValue(d, "uid", readResponse.Uid))
	errs.addMaybeError(SetStringOrIntValue(d, "gid", readResponse.Gid))
	errs.addMaybeError(d.Set("sid", readResponse.Sid))
	errs.addMaybeError(d.Set("name", readResponse.Name))
	errs.addMaybeError(d.Set("role_name", ids[0]))

	return errs.diags
}

func resourceRoleMemberDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	roleName := d.Get("role_name").(string)
	authId := d.Get("auth_id").(string)

	deleteRoleMemberUri := RolesEndpoint + roleName + MembersSuffix + authId

	tflog.Debug(ctx, fmt.Sprintf("Removing member with id %q from role with name %q", authId, roleName))

	_, err := DoRequest[RoleMemberDeleteBody, RoleMemberDeleteBody](ctx, c, DELETE, deleteRoleMemberUri, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

var RoleDomainValues = []string{
	"ACTIVE_DIRECTORY",
	"API_CREATOR_DOMAIN",
	"API_INTERNAL_DOMAIN",
	"API_INVALID_DOMAIN",
	"API_NULL_DOMAIN",
	"API_OPERATOR_DOMAIN",
	"API_RESERVED_DOMAIN",
	"LOCAL",
	"POSIX_GROUP",
	"POSIX_USER",
	"WORLD",
}
