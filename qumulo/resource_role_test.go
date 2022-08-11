package qumulo

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccCreateRole(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccRoleConfig(roleActors),
				Check: resource.ComposeTestCheckFunc(
					testAccCompareRoleResource(roleActors),
					testAccCheckRole(roleActors),
				),
			},
			{
				Config: testAccRoleConfig(roleActors2),
				Check: resource.ComposeTestCheckFunc(
					testAccCompareRoleResource(roleActors2),
					testAccCheckRole(roleActors2),
				),
			},
		},
	})
}

var roleActors = Role{
	Name:        "Actors",
	Description: "This role is for testing purposes",
	Privileges: []string{
		"PRIVILEGE_AD_READ",
		"PRIVILEGE_AD_USE",
		"PRIVILEGE_AD_WRITE",
		"PRIVILEGE_ANALYTICS_READ",
		"PRIVILEGE_AUDIT_READ",
		"PRIVILEGE_AUDIT_WRITE",
		"PRIVILEGE_AUTH_CACHE_READ",
		"PRIVILEGE_AUTH_CACHE_WRITE",
		"PRIVILEGE_CLUSTER_READ",
		"PRIVILEGE_CLUSTER_WRITE",
		"PRIVILEGE_DEBUG",
	},
}

var roleActors2 = Role{
	Name:        "Actors",
	Description: "This role is for testing purposes (part 2)",
	Privileges: []string{
		"PRIVILEGE_AD_READ",
	},
}

func testAccRoleConfig(role Role) string {
	return fmt.Sprintf(`
resource "qumulo_role" "test_role" {
	name        = %q
	description = %q
	privileges  = %v
}
  `, role.Name, role.Description, PrintTerraformListFromList(role.Privileges))
}

func testAccCompareRoleResource(role Role) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("qumulo_role.test_role", "name",
			role.Name),
		resource.TestCheckResourceAttr("qumulo_role.test_role", "description",
			role.Description),
		resource.TestCheckResourceAttr("qumulo_role.test_role", "privileges.#",
			fmt.Sprintf("%v", len(role.Privileges))),
	)
}

func testAccCheckRole(role Role) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		c := testAccProvider.Meta().(*Client)
		ctx := context.Background()

		roleResource, ok := s.RootModule().Resources["qumulo_role.test_role"]
		if !ok {
			return fmt.Errorf("role resource not found, %v", roleResource)
		}

		if roleResource.Primary.ID == "" {
			return fmt.Errorf("role ID is not set")
		}

		readRoleByNameUri := RolesEndpoint + roleResource.Primary.ID
		remoteRole, err := DoRequest[Role, Role](ctx, c, GET, readRoleByNameUri, nil)
		if err != nil {
			return err
		}

		if !(role.Description == remoteRole.Description) {
			return fmt.Errorf("roles descriptions mismatch: Expected %v, got %v", role.Description, remoteRole.Description)
		}
		if !reflect.DeepEqual(role.Privileges, remoteRole.Privileges) {
			return fmt.Errorf("roles privileges mismatch: Expected %v, got %v", role.Privileges, remoteRole.Privileges)
		}

		return nil
	}
}
