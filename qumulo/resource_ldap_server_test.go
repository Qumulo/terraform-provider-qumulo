package qumulo

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"reflect"
	"strings"
	"testing"
)

func TestAccChangeLdapServer(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{ // Reset state to default
				Config:             testAccDefaultLdapServerConfig(defaultLdapServerConfig),
				Check:              testAccCheckLdapServerSettings(defaultLdapServerConfigApplied),
				ExpectNonEmptyPlan: true,
			},
			{
				Config: testAccLdapServerConfig(testingLdapServerConfig),
				Check: resource.ComposeTestCheckFunc(
					testAccCompareLdapServerSettings(testingLdapServerConfig),
					testAccCheckLdapServerSettings(testingLdapServerConfig),
				),
			},
		},
	})
}

var defaultLdapServerConfig = LdapServerSettingsBody{
	UseLdap:           false,
	LdapSchema:        Rfc2307.String(),
	EncryptConnection: true,
}

var defaultLdapServerConfigApplied = LdapServerSettingsBody{
	UseLdap:    false,
	LdapSchema: Rfc2307.String(),
	// Default schema description that's applied by the API
	LdapSchemaDescription: LdapSchemaDescription{
		GroupMemberAttribute:         "memberUid",
		UserGroupIdentifierAttribute: "uid",
		LoginNameAttribute:           "uid",
		GroupNameAttribute:           "cn",
		UserObjectClass:              "posixAccount",
		GroupObjectClass:             "posixGroup",
		UidNumberAttribute:           "uidNumber",
		GidNumberAttribute:           "gidNumber",
	},
	EncryptConnection: true,
}

var testingLdapServerConfig = LdapServerSettingsBody{
	UseLdap:                true,
	BindUri:                "ldap://ldap.denvrdata.com",
	BaseDistinguishedNames: "dc=cloud,dc=denvrdata,dc=com",
	LdapSchema:             Custom.String(),
	LdapSchemaDescription: LdapSchemaDescription{
		GroupMemberAttribute:         "memberUid",
		UserGroupIdentifierAttribute: "uid",
		LoginNameAttribute:           "uid",
		GroupNameAttribute:           "cn",
		UserObjectClass:              "posixAccount",
		GroupObjectClass:             "posixGroup",
		UidNumberAttribute:           "uidNumber",
		GidNumberAttribute:           "gidNumber",
	},
	EncryptConnection: false,
}

func testAccDefaultLdapServerConfig(ldap LdapServerSettingsBody) string {
	return fmt.Sprintf(`
 resource "qumulo_ldap_server" "some_ldap_server" {
   use_ldap = %v
   bind_uri = %q
   user = %q
   base_distinguished_names = %q
   ldap_schema = %q
   encrypt_connection = %v
 }`, ldap.UseLdap, ldap.BindUri, ldap.User, ldap.BaseDistinguishedNames, ldap.LdapSchema,
		ldap.EncryptConnection)
}

func testAccLdapServerConfig(ldap LdapServerSettingsBody) string {
	return fmt.Sprintf(`
 resource "qumulo_ldap_server" "some_ldap_server" {
   use_ldap = %v
   bind_uri = %q
   user = %q
   base_distinguished_names = %v
   ldap_schema = %q
   ldap_schema_description {
     group_member_attribute = %q
     user_group_identifier_attribute = %q
     login_name_attribute =  %q
     group_name_attribute = %q
     user_object_class = %q
     group_object_class = %q
     uid_number_attribute = %q
     gid_number_attribute = %q
   }
   encrypt_connection = %v
 }`, ldap.UseLdap, ldap.BindUri, ldap.User, strings.ReplaceAll(fmt.Sprintf("%+q", ldap.BaseDistinguishedNames), "\" \"", "\", \""), ldap.LdapSchema, ldap.LdapSchemaDescription.GroupMemberAttribute,
		ldap.LdapSchemaDescription.UserGroupIdentifierAttribute, ldap.LdapSchemaDescription.LoginNameAttribute, ldap.LdapSchemaDescription.GroupNameAttribute,
		ldap.LdapSchemaDescription.UserObjectClass, ldap.LdapSchemaDescription.GroupObjectClass, ldap.LdapSchemaDescription.UidNumberAttribute, ldap.LdapSchemaDescription.GidNumberAttribute,
		ldap.EncryptConnection)
}

func testAccCompareLdapServerSettings(ldap LdapServerSettingsBody) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("qumulo_ldap_server.some_ldap_server", "use_ldap",
			fmt.Sprintf("%v", ldap.UseLdap)),
		resource.TestCheckResourceAttr("qumulo_ldap_server.some_ldap_server", "bind_uri",
			fmt.Sprintf("%v", ldap.BindUri)),
		resource.TestCheckResourceAttr("qumulo_ldap_server.some_ldap_server", "user",
			fmt.Sprintf("%v", ldap.User)),
		resource.TestCheckResourceAttr("qumulo_ldap_server.some_ldap_server", "base_distinguished_names",
			fmt.Sprintf("%v", ldap.BaseDistinguishedNames)),
		resource.TestCheckResourceAttr("qumulo_ldap_server.some_ldap_server", "ldap_schema",
			fmt.Sprintf("%v", ldap.LdapSchema)),
		resource.TestCheckResourceAttr("qumulo_ldap_server.some_ldap_server", "ldap_schema_description.#",
			"1"),
		resource.TestCheckResourceAttr("qumulo_ldap_server.some_ldap_server", "ldap_schema_description.0.%",
			"8"),
		resource.TestCheckResourceAttr("qumulo_ldap_server.some_ldap_server", "ldap_schema_description.0.group_member_attribute",
			fmt.Sprintf("%v", ldap.LdapSchemaDescription.GroupMemberAttribute)),
		resource.TestCheckResourceAttr("qumulo_ldap_server.some_ldap_server", "ldap_schema_description.0.user_group_identifier_attribute",
			fmt.Sprintf("%v", ldap.LdapSchemaDescription.UserGroupIdentifierAttribute)),
		resource.TestCheckResourceAttr("qumulo_ldap_server.some_ldap_server", "ldap_schema_description.0.login_name_attribute",
			fmt.Sprintf("%v", ldap.LdapSchemaDescription.LoginNameAttribute)),
		resource.TestCheckResourceAttr("qumulo_ldap_server.some_ldap_server", "ldap_schema_description.0.group_name_attribute",
			fmt.Sprintf("%v", ldap.LdapSchemaDescription.GroupNameAttribute)),
		resource.TestCheckResourceAttr("qumulo_ldap_server.some_ldap_server", "ldap_schema_description.0.user_object_class",
			fmt.Sprintf("%v", ldap.LdapSchemaDescription.UserObjectClass)),
		resource.TestCheckResourceAttr("qumulo_ldap_server.some_ldap_server", "ldap_schema_description.0.group_object_class",
			fmt.Sprintf("%v", ldap.LdapSchemaDescription.GroupObjectClass)),
		resource.TestCheckResourceAttr("qumulo_ldap_server.some_ldap_server", "ldap_schema_description.0.uid_number_attribute",
			fmt.Sprintf("%v", ldap.LdapSchemaDescription.UidNumberAttribute)),
		resource.TestCheckResourceAttr("qumulo_ldap_server.some_ldap_server", "ldap_schema_description.0.gid_number_attribute",
			fmt.Sprintf("%v", ldap.LdapSchemaDescription.GidNumberAttribute)),
		resource.TestCheckResourceAttr("qumulo_ldap_server.some_ldap_server", "encrypt_connection",
			fmt.Sprintf("%v", ldap.EncryptConnection)),
	)
}

func testAccCheckLdapServerSettings(ldap LdapServerSettingsBody) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		c := testAccProvider.Meta().(*Client)
		ctx := context.Background()
		settings, err := DoRequest[LdapServerSettingsBody, LdapServerSettingsBody](ctx, c, GET, LdapServerEndpoint, nil)
		if err != nil {
			return err
		}
		if !reflect.DeepEqual(*settings, ldap) {
			return fmt.Errorf("SMB server settings mismatch: Expected %v, got %v", ldap, *settings)
		}
		return nil
	}
}
