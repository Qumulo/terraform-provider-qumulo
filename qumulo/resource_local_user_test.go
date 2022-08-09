package qumulo

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccCreateUser(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccUserConfig(userTest1),
				Check: resource.ComposeTestCheckFunc(
					testAccCompareUserResource(userTest1),
					testAccCheckUser(userTest1),
				),
			},
			{
				Config: testAccUserConfig(userTest2),
				Check: resource.ComposeTestCheckFunc(
					testAccCompareUserResource(userTest2),
					testAccCheckUser(userTest2),
				),
			},
		},
	})
}

var userTest1 = UserBody{
	Name:         "test_user1",
	PrimaryGroup: "514",
	Password:     "Test1234",
}

var userTest2 = UserBody{
	Name:          "test_user2",
	PrimaryGroup:  "513",
	Uid:           "123",
	HomeDirectory: "/",
	Password:      "Test1234",
}

func testAccUserConfig(user UserBody) string {
	return fmt.Sprintf(`
resource "qumulo_local_user" "test_user" {
	name = %q
	primary_group = %v
	uid = "%v"
	home_directory = %q
	password = %q
}
  `, user.Name, user.PrimaryGroup, user.Uid, user.HomeDirectory, user.Password)
}

func testAccCompareUserResource(user UserBody) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("qumulo_local_user.test_user", "name",
			user.Name),
		resource.TestCheckResourceAttr("qumulo_local_user.test_user", "primary_group",
			user.PrimaryGroup),
		resource.TestCheckResourceAttr("qumulo_local_user.test_user", "uid",
			user.Uid),
		resource.TestCheckResourceAttr("qumulo_local_user.test_user", "home_directory",
			user.HomeDirectory),
		resource.TestCheckResourceAttr("qumulo_local_user.test_user", "password",
			user.Password),
	)
}

func testAccCheckUser(user UserBody) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		c := testAccProvider.Meta().(*Client)
		ctx := context.Background()

		userResource, ok := s.RootModule().Resources["qumulo_local_user.test_user"]
		if !ok {
			return fmt.Errorf("user resource not found, %v", userResource)
		}

		if userResource.Primary.ID == "" {
			return fmt.Errorf("user ID is not set")
		}

		readUserByNameUri := UsersEndpoint + userResource.Primary.ID
		remoteUser, err := DoRequest[UserBody, UserBody](ctx, c, GET, readUserByNameUri, nil)
		if err != nil {
			return err
		}

		if !(user.Name == remoteUser.Name) {
			return fmt.Errorf("users name mismatch: Expected %v, got %v", user.Name, remoteUser.Name)
		}
		if !(user.PrimaryGroup == remoteUser.PrimaryGroup) {
			return fmt.Errorf("users name mismatch: Expected %v, got %v", user.PrimaryGroup, remoteUser.PrimaryGroup)
		}
		if !(user.Uid == remoteUser.Uid) {
			return fmt.Errorf("users name mismatch: Expected %v, got %v", user.Uid, remoteUser.Uid)
		}

		return nil
	}
}
