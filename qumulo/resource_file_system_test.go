package qumulo

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccCreateFileSystemSettings(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccFileSystemSettingsConfig(defaultFileSystemSettings),
				Check: resource.ComposeTestCheckFunc(
					testAccCompareFileSystemSettingsResource(defaultFileSystemSettings),
					testAccCheckFileSystemPermissionsSettings(defaultPermissionsSettings),
					testAccCheckFileSystemAtimeSettings(defaultAtimeSettings),
				),
			},
			{
				Config: testAccFileSystemSettingsConfig(testFileSystemSettings),
				Check: resource.ComposeTestCheckFunc(
					testAccCompareFileSystemSettingsResource(testFileSystemSettings),
					testAccCheckFileSystemPermissionsSettings(testPermissionsSettings),
					testAccCheckFileSystemAtimeSettings(testAtimeSettings),
				),
			},
			{
				// reset to default state; don't want to leave weird settings enabled post-test
				Config: testAccFileSystemSettingsConfig(defaultFileSystemSettings),
				Check: resource.ComposeTestCheckFunc(
					testAccCompareFileSystemSettingsResource(defaultFileSystemSettings),
					testAccCheckFileSystemPermissionsSettings(defaultPermissionsSettings),
					testAccCheckFileSystemAtimeSettings(defaultAtimeSettings),
				),
			},
		},
	})
}

var defaultPermissionsSettings = FileSystemPermissionsSettingsBody{
	Mode: CrossProtocol.String(),
}

var defaultAtimeSettings = FileSystemAtimeSettingsBody{
	Enabled:     false,
	Granularity: Hour.String(),
}

var defaultFileSystemSettings = FileSystemSettingsBody{
	Permissions:   &defaultPermissionsSettings,
	AtimeSettings: &defaultAtimeSettings,
}

var testPermissionsSettings = FileSystemPermissionsSettingsBody{
	Mode: Native.String(),
}

var testAtimeSettings = FileSystemAtimeSettingsBody{
	Enabled:     true,
	Granularity: Day.String(),
}

var testFileSystemSettings = FileSystemSettingsBody{
	Permissions:   &testPermissionsSettings,
	AtimeSettings: &testAtimeSettings,
}

func testAccFileSystemSettingsConfig(settings FileSystemSettingsBody) string {
	return fmt.Sprintf(`
resource "qumulo_file_system_settings" "test_fs_settings" {
	permissions_mode = %q
	atime_enabled = %v
	atime_granularity = %q
}
  `, settings.Permissions.Mode, settings.AtimeSettings.Enabled, settings.AtimeSettings.Granularity)
}

func testAccCompareFileSystemSettingsResource(settings FileSystemSettingsBody) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("qumulo_file_system_settings.test_fs_settings", "permissions_mode",
			settings.Permissions.Mode),
		resource.TestCheckResourceAttr("qumulo_file_system_settings.test_fs_settings", "atime_enabled",
			fmt.Sprintf("%v", settings.AtimeSettings.Enabled)),
		resource.TestCheckResourceAttr("qumulo_file_system_settings.test_fs_settings", "atime_granularity",
			settings.AtimeSettings.Granularity),
	)
}

func testAccCheckFileSystemPermissionsSettings(permissions FileSystemPermissionsSettingsBody) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		c := testAccProvider.Meta().(*Client)
		ctx := context.Background()

		remotePermissions, err := DoRequest[FileSystemPermissionsSettingsBody, FileSystemPermissionsSettingsBody](ctx,
			c, GET, FileSystemPermissionsEndpoint, nil)
		if err != nil {
			return err
		}

		if !(permissions.Mode == remotePermissions.Mode) {
			return fmt.Errorf("file system permissions mode mismatch: Expected %v, got %v", permissions.Mode, remotePermissions.Mode)
		}

		return nil
	}
}

func testAccCheckFileSystemAtimeSettings(atime FileSystemAtimeSettingsBody) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		c := testAccProvider.Meta().(*Client)
		ctx := context.Background()

		remoteAtime, err := DoRequest[FileSystemAtimeSettingsBody, FileSystemAtimeSettingsBody](ctx,
			c, GET, FileSystemAtimeEndpoint, nil)
		if err != nil {
			return err
		}

		if !reflect.DeepEqual(*remoteAtime, atime) {
			return fmt.Errorf("file system atime settings mismatch: Expected %v, got %v", atime, remoteAtime)
		}

		return nil
	}
}
