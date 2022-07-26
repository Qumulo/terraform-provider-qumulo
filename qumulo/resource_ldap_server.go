package qumulo

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

type LDAPSchema int

const (
	RFC2307 LDAPSchema = iota + 1
	CUSTOM
)

func (e LDAPSchema) String() string {
	return ldapSchemaList[e-1]
}

const LdapServerEndpoint = "/v2/ldap/settings"

type LdapServerSettings struct {
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

var ldapSchemaList = []string{"RFC2307", "CUSTOM"}

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
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice(ldapSchemaList, false)),
				Default:          RFC2307,
			},
			"ldap_schema_description": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Required: true,
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
	}
}

func resourceLdapServerCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	ldapSettings := LdapServerSettings{
		UseLdap:                d.Get("use_ldap").(bool),
		BindUri:                d.Get("bind_uri").(string),
		BaseDistinguishedNames: d.Get("base_distinguished_names").(string),
		LdapSchemaDescription:  expandLdapDescription(d.Get("ldap_schema_description").([]interface{})),
		LdapSchema:             d.Get("ldap_schema").(string),
		EncryptConnection:      d.Get("encrypt_connection").(bool),
	}

	if d.Get("user").(string) != "" {
		ldapSettings.User = d.Get("user").(string)
	}

	if d.Get("password").(string) != "" {
		ldapSettings.Password = d.Get("password").(string)
	}

	_, err := DoRequest[LdapServerSettings, LdapServerSettings](client, PUT, LdapServerEndpoint, &ldapSettings)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return resourceLdapServerRead(ctx, d, m)
}

func resourceLdapServerRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	var diags diag.Diagnostics

	ls, err := DoRequest[LdapServerSettings, LdapServerSettings](c, GET, LdapServerEndpoint, nil)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("use_ldap", ls.UseLdap)
	d.Set("bind_uri", ls.BindUri)
	d.Set("user", ls.User)

	d.Set("password", ls.Password)
	d.Set("base_distinguished_names", ls.BaseDistinguishedNames)
	d.Set("ldap_schema", ls.LdapSchema)
	err = d.Set("ldap_schema_description", flattenLdapSchemaDescription(
		ls.LdapSchemaDescription))
	if err != nil {
		return diag.FromErr(fmt.Errorf("error setting Ldap schema description: %w", err))
	}
	d.Set("encrypt_connection", ls.EncryptConnection)
	return diags
}

func resourceLdapServerUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	ldapSettings := LdapServerSettings{
		UseLdap:                d.Get("use_ldap").(bool),
		BindUri:                d.Get("bind_uri").(string),
		BaseDistinguishedNames: d.Get("base_distinguished_names").(string),
		LdapSchemaDescription:  expandLdapDescription(d.Get("ldap_schema_description").([]interface{})),
		LdapSchema:             d.Get("ldap_schema").(string),
		EncryptConnection:      d.Get("encrypt_connection").(bool),
	}

	if d.Get("user").(string) != "" {
		ldapSettings.User = d.Get("user").(string)
	}

	if d.Get("password").(string) != "" {
		ldapSettings.Password = d.Get("password").(string)
	}

	_, err := DoRequest[LdapServerSettings, LdapServerSettings](client, PATCH, LdapServerEndpoint, &ldapSettings)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceLdapServerRead(ctx, d, m)
}

func resourceLdapServerDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	return diags
}

func expandLdapDescription(tfLdapSchemaDescriptions []interface{}) LdapSchemaDescription {
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
		apiObject.GroupMemberAttribute = tfMap["group_member_attribute"].(string)
	}

	if v, ok := tfMap["user_group_identifier_attribute"].(string); ok && v != "" {
		apiObject.UserGroupIdentifierAttribute = tfMap["user_group_identifier_attribute"].(string)
	}

	if v, ok := tfMap["login_name_attribute"].(string); ok && v != "" {
		apiObject.LoginNameAttribute = tfMap["login_name_attribute"].(string)
	}

	if v, ok := tfMap["group_name_attribute"].(string); ok && v != "" {
		apiObject.GroupNameAttribute = tfMap["group_name_attribute"].(string)
	}
	if v, ok := tfMap["user_object_class"].(string); ok && v != "" {
		apiObject.UserObjectClass = tfMap["user_object_class"].(string)
	}
	if v, ok := tfMap["group_object_class"].(string); ok && v != "" {
		apiObject.GroupObjectClass = tfMap["group_object_class"].(string)
	}
	if v, ok := tfMap["uid_number_attribute"].(string); ok && v != "" {
		apiObject.UidNumberAttribute = tfMap["uid_number_attribute"].(string)
	}
	if v, ok := tfMap["gid_number_attribute"].(string); ok && v != "" {
		apiObject.GidNumberAttribute = tfMap["gid_number_attribute"].(string)
	}
	return apiObject
}

func flattenLdapSchemaDescription(apiObject LdapSchemaDescription) []interface{} {
	tfMap := map[string]interface{}{}

	if v := apiObject.GroupMemberAttribute; v != "" {
		tfMap["group_member_attribute"] = apiObject.GroupMemberAttribute
	}

	if v := apiObject.UserGroupIdentifierAttribute; v != "" {
		tfMap["user_group_identifier_attribute"] = apiObject.UserGroupIdentifierAttribute
	}

	if v := apiObject.LoginNameAttribute; v != "" {
		tfMap["login_name_attribute"] = apiObject.LoginNameAttribute
	}

	if v := apiObject.GroupNameAttribute; v != "" {
		tfMap["group_name_attribute"] = apiObject.GroupNameAttribute
	}

	if v := apiObject.UserObjectClass; v != "" {
		tfMap["user_object_class"] = apiObject.UserObjectClass
	}

	if v := apiObject.GroupObjectClass; v != "" {
		tfMap["group_object_class"] = apiObject.GroupObjectClass
	}

	if v := apiObject.UidNumberAttribute; v != "" {
		tfMap["uid_number_attribute"] = apiObject.UidNumberAttribute
	}

	if v := apiObject.GidNumberAttribute; v != "" {
		tfMap["gid_number_attribute"] = apiObject.GidNumberAttribute
	}

	return []interface{}{tfMap}
}
