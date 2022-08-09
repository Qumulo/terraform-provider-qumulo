package qumulo

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

type LdapSchema int

const (
	Rfc2307 LdapSchema = iota + 1
	Custom
)

func (e LdapSchema) String() string {
	return LdapSchemaValues[e-1]
}

const LdapServerEndpoint = "/v2/ldap/settings"

type LdapServerSettingsBody struct {
	UseLdap                bool                  `json:"use_ldap"`
	BindUri                string                `json:"bind_uri"`
	User                   string                `json:"user"`
	Password               string                `json:"password"`
	BaseDistinguishedNames string                `json:"base_distinguished_names"`
	LdapSchema             string                `json:"ldap_schema"`
	LdapSchemaDescription  LdapSchemaDescription `json:"ldap_schema_description"`
	EncryptConnection      bool                  `json:"encrypt_connection"`
}

type LdapSchemaDescription struct {
	GroupMemberAttribute         string `json:"group_member_attribute"`
	UserGroupIdentifierAttribute string `json:"user_group_identifier_attribute"`
	LoginNameAttribute           string `json:"login_name_attribute"`
	GroupNameAttribute           string `json:"group_name_attribute"`
	UserObjectClass              string `json:"user_object_class"`
	GroupObjectClass             string `json:"group_object_class"`
	UidNumberAttribute           string `json:"uid_number_attribute"`
	GidNumberAttribute           string `json:"gid_number_attribute"`
}

var LdapSchemaValues = []string{"RFC2307", "CUSTOM"}

func resourceLdapServer() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLdapServerCreate,
		ReadContext:   resourceLdapServerRead,
		UpdateContext: resourceLdapServerUpdate,
		DeleteContext: resourceLdapServerDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"use_ldap": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"bind_uri": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"user": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"password": &schema.Schema{
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"base_distinguished_names": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"ldap_schema": &schema.Schema{
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice(LdapSchemaValues, false)),
				Default:          Rfc2307.String(),
			},
			"ldap_schema_description": {
				Type:     schema.TypeList,
				MaxItems: 1,
				// API applies a default config for ldap schema description if ldap_schema = RFC2307
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"group_member_attribute": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"user_group_identifier_attribute": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"login_name_attribute": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"group_name_attribute": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"user_object_class": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"group_object_class": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"uid_number_attribute": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"gid_number_attribute": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"encrypt_connection": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
		},

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceLdapServerCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	err := setLdapServerSettings(ctx, d, m, PUT)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return resourceLdapServerRead(ctx, d, m)
}

func resourceLdapServerRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	var errs ErrorCollection

	ls, err := DoRequest[LdapServerSettingsBody, LdapServerSettingsBody](ctx, c, GET, LdapServerEndpoint, nil)
	if err != nil {
		return diag.FromErr(err)
	}
	errs.addMaybeError(d.Set("use_ldap", ls.UseLdap))
	errs.addMaybeError(d.Set("bind_uri", ls.BindUri))
	errs.addMaybeError(d.Set("user", ls.User))
	errs.addMaybeError(d.Set("password", ls.Password))
	errs.addMaybeError(d.Set("base_distinguished_names", ls.BaseDistinguishedNames))
	errs.addMaybeError(d.Set("ldap_schema", ls.LdapSchema))
	errs.addMaybeError(d.Set("encrypt_connection", ls.EncryptConnection))

	err = d.Set("ldap_schema_description", flattenLdapSchemaDescription(
		ls.LdapSchemaDescription))
	if err != nil {
		return diag.FromErr(fmt.Errorf("error setting Ldap schema description: %w", err))
	}
	return errs.diags
}

func resourceLdapServerUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	err := setLdapServerSettings(ctx, d, m, PATCH)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceLdapServerRead(ctx, d, m)
}

func resourceLdapServerDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "Deleting LDAP settings resource")

	return nil
}

func setLdapServerSettings(ctx context.Context, d *schema.ResourceData, m interface{}, method Method) error {
	c := m.(*Client)

	ldapServerSettings := LdapServerSettingsBody{
		UseLdap:                d.Get("use_ldap").(bool),
		BindUri:                d.Get("bind_uri").(string),
		BaseDistinguishedNames: d.Get("base_distinguished_names").(string),
		LdapSchemaDescription:  expandLdapSchemaDescription(d.Get("ldap_schema_description").([]interface{})),
		LdapSchema:             d.Get("ldap_schema").(string),
		EncryptConnection:      d.Get("encrypt_connection").(bool),
	}

	if v := d.Get("user").(string); v != "" {
		ldapServerSettings.User = v
	}

	if v := d.Get("password").(string); v != "" {
		ldapServerSettings.Password = v
	}

	tflog.Debug(ctx, "Updating LDAP settings")
	_, err := DoRequest[LdapServerSettingsBody, LdapServerSettingsBody](ctx, c, method, LdapServerEndpoint, &ldapServerSettings)
	return err
}

func expandLdapSchemaDescription(tfLdapSchemaDescriptions []interface{}) LdapSchemaDescription {
	apiObject := LdapSchemaDescription{}

	if len(tfLdapSchemaDescriptions) == 0 {
		return apiObject
	}
	tfLdapSchemaDescription := tfLdapSchemaDescriptions[0]
	tfMap, ok := tfLdapSchemaDescription.(map[string]interface{})
	if !ok {
		return apiObject
	}

	if v, ok := tfMap["group_member_attribute"].(string); ok && v != "" {
		apiObject.GroupMemberAttribute = v
	}

	if v, ok := tfMap["user_group_identifier_attribute"].(string); ok && v != "" {
		apiObject.UserGroupIdentifierAttribute = v
	}

	if v, ok := tfMap["login_name_attribute"].(string); ok && v != "" {
		apiObject.LoginNameAttribute = v
	}

	if v, ok := tfMap["group_name_attribute"].(string); ok && v != "" {
		apiObject.GroupNameAttribute = v
	}
	if v, ok := tfMap["user_object_class"].(string); ok && v != "" {
		apiObject.UserObjectClass = v
	}
	if v, ok := tfMap["group_object_class"].(string); ok && v != "" {
		apiObject.GroupObjectClass = v
	}
	if v, ok := tfMap["uid_number_attribute"].(string); ok && v != "" {
		apiObject.UidNumberAttribute = v
	}
	if v, ok := tfMap["gid_number_attribute"].(string); ok && v != "" {
		apiObject.GidNumberAttribute = v
	}
	return apiObject
}

func flattenLdapSchemaDescription(apiObject LdapSchemaDescription) []interface{} {
	tfMap := map[string]interface{}{}

	if v := apiObject.GroupMemberAttribute; v != "" {
		tfMap["group_member_attribute"] = v
	}

	if v := apiObject.UserGroupIdentifierAttribute; v != "" {
		tfMap["user_group_identifier_attribute"] = v
	}

	if v := apiObject.LoginNameAttribute; v != "" {
		tfMap["login_name_attribute"] = v
	}

	if v := apiObject.GroupNameAttribute; v != "" {
		tfMap["group_name_attribute"] = v
	}

	if v := apiObject.UserObjectClass; v != "" {
		tfMap["user_object_class"] = v
	}

	if v := apiObject.GroupObjectClass; v != "" {
		tfMap["group_object_class"] = v
	}

	if v := apiObject.UidNumberAttribute; v != "" {
		tfMap["uid_number_attribute"] = v
	}

	if v := apiObject.GidNumberAttribute; v != "" {
		tfMap["gid_number_attribute"] = v
	}

	return []interface{}{tfMap}
}
