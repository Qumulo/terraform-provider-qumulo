package qumulo

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccCreateGroup(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccGroupConfig(groupTest1),
				Check: resource.ComposeTestCheckFunc(
					testAccCompareGroupResource(groupTest1),
					testAccCheckGroup(groupTest1),
				),
			},
			{
				Config: testAccGroupConfig(groupTest2),
				Check: resource.ComposeTestCheckFunc(
					testAccCompareGroupResource(groupTest2),
					testAccCheckGroup(groupTest2),
				),
			},
		},
	})
}

var groupTest1 = CreateGroupRequest{
	Name: "test_group",
	Gid:  "",
}

var groupTest2 = CreateGroupRequest{
	Name: "test_group2",
	Gid:  "",
}

func testAccGroupConfig(group CreateGroupRequest) string {
	return fmt.Sprintf(`
resource "qumulo_local_group" "test_group" {
	name = %q
	gid = %q
}
  `, group.Name, group.Gid)
}

func testAccCompareGroupResource(group CreateGroupRequest) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("qumulo_local_group.test_group", "name",
			group.Name),
		resource.TestCheckResourceAttr("qumulo_local_group.test_group", "gid",
			group.Gid),
	)
}

func testAccCheckGroup(group CreateGroupRequest) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		c := testAccProvider.Meta().(*Client)
		ctx := context.Background()

		groupResource, ok := s.RootModule().Resources["qumulo_local_group.test_group"]
		if !ok {
			return fmt.Errorf("Group resource not found, %v", groupResource)
		}

		if groupResource.Primary.ID == "" {
			return fmt.Errorf("Group ID is not set")
		}

		readGroupByNameUri := GroupsEndpoint + groupResource.Primary.ID
		remoteGroup, err := DoRequest[GroupResponse, GroupResponse](ctx, c, GET, readGroupByNameUri, nil)
		if err != nil {
			return err
		}

		if group.Name != remoteGroup.Name {
			return fmt.Errorf("Group name mismatch: Expected %v, got %v", group.Name, remoteGroup.Name)
		}
		if group.Gid != remoteGroup.Gid {
			return fmt.Errorf("Group NFS GID mismatch: Expected %v, got %v", group.Gid, remoteGroup.Gid)
		}

		return nil
	}
}
