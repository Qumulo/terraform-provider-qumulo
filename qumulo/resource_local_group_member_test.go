package qumulo

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAddGroupMember(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMemberConfig(group1, user1),
				Check:  testAccValidateGroupMember(),
			},
		},
	})
}

var user1 = UserBody{
	Name:         "test_user4",
	PrimaryGroup: "514",
	Password:     "Test1234",
}

var group1 = CreateGroupRequest{
	Name: "test_group4",
	Gid:  "",
}

func testAccMemberConfig(group CreateGroupRequest, user UserBody) string {
	return fmt.Sprintf(`
resource "qumulo_local_group" "test_group" {
	name = %q
	gid = %q
}

resource "qumulo_local_user" "test_user" {
	name = %q
	primary_group = %v
	uid = "%v"
	home_directory = %q
	password = %q
}

resource "qumulo_local_group_member" "test_member" {
	member_id = qumulo_local_user.test_user.id
	group_id = qumulo_local_group.test_group.id
}
  `, group.Name, group.Gid, user.Name, user.PrimaryGroup, user.Uid, user.HomeDirectory, user.Password)
}

func testAccValidateGroupMember() resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttrSet("qumulo_local_group_member.test_member", "member_id"),
		resource.TestCheckResourceAttrSet("qumulo_local_group_member.test_member", "group_id"),
	)
}
