package qumulo

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const UsersEndpoint = "/v1/users/"

type UserBody struct {
	Id                string `json:"id"`
	Name              string `json:"name"`
	PrimaryGroup      string `json:"primary_group"`
	Sid               string `json:"sid"`
	Uid               string `json:"uid"`
	HomeDirectory     string `json:"home_directory"`
	CanChangePassword bool   `json:"can_change_password"`
	Password          string `json:"password"`
}

type UserModify struct {
	Id            string `json:"id"`
	Name          string `json:"name"`
	PrimaryGroup  string `json:"primary_group"`
	Uid           string `json:"uid"`
	HomeDirectory string `json:"home_directory"`
}

func resourceUser() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceUserCreate,
		ReadContext:   resourceUserRead,
		UpdateContext: resourceUserUpdate,
		DeleteContext: resourceUserDelete,

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
			"primary_group": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"sid": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"uid": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"home_directory": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"password": &schema.Schema{
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"can_change_password": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
		},

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceUserCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	if d.Get("password") == "" {
		return diag.FromErr(fmt.Errorf("password must be set when creating user"))
	}

	userSettings := setUserSettings(ctx, d, m)

	tflog.Debug(ctx, fmt.Sprintf("Creating local user with name %q", userSettings.Name))

	user, err := DoRequest[UserBody, UserBody](ctx, c, POST, UsersEndpoint, &userSettings)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(user.Id)

	return resourceUserRead(ctx, d, m)

}

func resourceUserRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	var errs ErrorCollection

	readUserByIdUri := UsersEndpoint + d.Id()

	user, err := DoRequest[UserBody, UserBody](ctx, c, GET, readUserByIdUri, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	errs.addMaybeError(d.Set("name", user.Name))
	errs.addMaybeError(d.Set("primary_group", user.PrimaryGroup))
	errs.addMaybeError(d.Set("sid", user.Sid))
	errs.addMaybeError(d.Set("uid", user.Uid))
	errs.addMaybeError(d.Set("home_directory", user.HomeDirectory))
	errs.addMaybeError(d.Set("can_change_password", user.CanChangePassword))

	return errs.diags
}

func resourceUserUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	updateUserByNameUri := UsersEndpoint + d.Id()

	tflog.Debug(ctx, fmt.Sprintf("Updating local user with name %q", d.Get("name")))

	var err error

	if d.Get("password") != "" {
		userSettings := setUserSettings(ctx, d, m)
		userSettings.Id = d.Id()
		_, err = DoRequest[UserBody, UserBody](ctx, c, PUT, updateUserByNameUri, &userSettings)
	} else {
		userSettings := modifyUserSettings(ctx, d, m)
		userSettings.Id = d.Id()
		_, err = DoRequest[UserModify, UserModify](ctx, c, PUT, updateUserByNameUri, &userSettings)
	}

	if err != nil {
		return diag.FromErr(err)
	}

	return resourceUserRead(ctx, d, m)
}

func resourceUserDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	deleteUserByNameUri := UsersEndpoint + d.Id()

	tflog.Debug(ctx, fmt.Sprintf("Deleting local user with id %q", d.Id()))

	_, err := DoRequest[UserBody, UserBody](ctx, c, DELETE, deleteUserByNameUri, nil)

	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func setUserSettings(ctx context.Context, d *schema.ResourceData, m interface{}) UserBody {

	userConfig := UserBody{
		Name:          d.Get("name").(string),
		PrimaryGroup:  d.Get("primary_group").(string),
		Uid:           d.Get("uid").(string),
		HomeDirectory: d.Get("home_directory").(string),
		Password:      d.Get("password").(string),
	}

	return userConfig
}

func modifyUserSettings(ctx context.Context, d *schema.ResourceData, m interface{}) UserModify {

	userConfig := UserModify{
		Name:          d.Get("name").(string),
		PrimaryGroup:  d.Get("primary_group").(string),
		Uid:           d.Get("uid").(string),
		HomeDirectory: d.Get("home_directory").(string),
	}

	return userConfig
}
