package qumulo

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const MembersSuffix = "/members/"

type MemberRequest struct {
	MemberId string `json:"member_id"`
	GroupId  string `json:"group_id"`
}

func resourceGroupMember() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGroupMemberCreate,
		ReadContext:   resourceGroupMemberRead,
		DeleteContext: resourceGroupMemberDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"member_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"group_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},

		// Note that the import function does not verify with the API whether this user is actually a member
		// of this group
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				// d.Id() here is the last argument passed to the
				// `terraform import RESOURCE_TYPE.RESOURCE_NAME RESOURCE_ID` command
				ids, err := ParseLocalGroupMemberId(d.Id())

				if err != nil {
					return nil, err
				}

				// Now, ids is guaranteed to have length 2, and d.Id() is well formed
				d.Set("group_id", ids[0])
				d.Set("member_id", ids[1])

				return []*schema.ResourceData{d}, nil
			},
		},
	}
}

func resourceGroupMemberCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	memberSettings := MemberRequest{
		MemberId: d.Get("member_id").(string),
		GroupId:  d.Get("group_id").(string),
	}

	tflog.Debug(ctx, fmt.Sprintf("Adding member with id %q to group with id %q",
		memberSettings.MemberId, memberSettings.GroupId))

	addMemberUri := GroupsEndpoint + memberSettings.GroupId + MembersSuffix

	_, err := DoRequest[MemberRequest, MemberRequest](ctx, c, POST, addMemberUri, &memberSettings)
	if err != nil {
		return diag.FromErr(err)
	}

	// The format of a member resource ID is of the form {group_id}:{user_id}
	ids := make([]string, 2)
	ids[0] = memberSettings.GroupId
	ids[1] = memberSettings.MemberId
	id, err := FormLocalGroupMemberId(ids)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	return resourceGroupMemberRead(ctx, d, m)
}

func resourceGroupMemberRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	tflog.Debug(ctx, "Reading group member; This is a no-op")
	return nil
}

func resourceGroupMemberDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	deleteGroupMemberByIdUri := GroupsEndpoint + d.Get("group_id").(string) +
		MembersSuffix + d.Get("member_id").(string)

	tflog.Debug(ctx, fmt.Sprintf("Removing member with id %q to group with id %q",
		d.Get("member_id").(string), d.Get("group_id").(string)))

	_, err := DoRequest[MemberRequest, MemberRequest](ctx, c, DELETE, deleteGroupMemberByIdUri, nil)

	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
