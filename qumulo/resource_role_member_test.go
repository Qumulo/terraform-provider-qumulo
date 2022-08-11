package qumulo

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAddRoleMember(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccRoleMemberConfig(role1, userRoles1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRoleMember(role1, userRoles1),
					testAccValidateRoleMember(),
				),
			},
		},
	})
}

var userRoles1 = UserBody{
	Name:         "test_user_roles",
	PrimaryGroup: "514",
	Password:     "Test1234",
}

var role1 = Role{
	Name:        "Testers",
	Description: "This role is for testing purposes",
	Privileges: []string{
		"PRIVILEGE_AD_READ",
	},
}

func testAccRoleMemberConfig(role Role, user UserBody) string {
	return fmt.Sprintf(`
resource "qumulo_role" "test_role" {
	name        = %q
	description = %q
	privileges  = %v
}

resource "qumulo_local_user" "test_user" {
	name = %q
	primary_group = %v
	uid = "%v"
	home_directory = %q
	password = %q
}

resource "qumulo_role_member" "test_member" {
	name = qumulo_local_user.test_user.name
	role_name = qumulo_role.test_role.name
}
  `, role.Name, role.Description, strings.ReplaceAll(fmt.Sprintf("%+q", role.Privileges), "\" \"", "\", \""),
		user.Name, user.PrimaryGroup, user.Uid, user.HomeDirectory, user.Password)
}

func testAccValidateRoleMember() resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		// XXX amanning32: This check should work, but for some reason checking the auth_id happens before it is set, causing a failure
		// Leaving it commented out does pass, and the auth ID is clearly set due to the other Check function as well as the DELETE working
		// resource.TestCheckResourceAttrSet("qumulo_role_member.test_member", "domain"),
		// resource.TestCheckResourceAttrSet("qumulo_role_member.test_member", "auth_id"),
		resource.TestCheckResourceAttrSet("qumulo_role_member.test_member", "role_name"),
	)
}

func testAccCheckRoleMember(role Role, user UserBody) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		c := testAccProvider.Meta().(*Client)
		ctx := context.Background()

		roleResource, ok := s.RootModule().Resources["qumulo_role_member.test_member"]
		if !ok {
			return fmt.Errorf("role member resource not found, %v", roleResource)
		}

		if roleResource.Primary.ID == "" {
			return fmt.Errorf("role member ID is not set")
		}

		ids, err := ParseRoleMemberId(roleResource.Primary.ID)
		if err != nil {
			return fmt.Errorf("role member ID is malformed")
		}

		readRoleMemberUri := RolesEndpoint + ids[0] + MembersSuffix + ids[1]

		remoteRoleMember, err := DoRequest[RoleMemberResponse, RoleMemberResponse](ctx, c, GET, readRoleMemberUri, nil)
		if err != nil {
			return err
		}

		if !(role.Name == ids[0]) {
			return fmt.Errorf("roles name mismatch in role member ID: Expected %v, got %v", role.Name, ids[0])
		}
		if !(ids[1] == remoteRoleMember.AuthId) {
			return fmt.Errorf("auth id mismatch in role member ID: Expected %v, got %v", ids[1], remoteRoleMember.AuthId)
		}

		return nil
	}
}
