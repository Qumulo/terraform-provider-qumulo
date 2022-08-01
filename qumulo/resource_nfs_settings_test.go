package qumulo

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"testing"
)

func TestAccChangeNfsSettings(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				//Reset to default
				Config: testAccNfsSettings(defaultNfsSettings),
				Check:  testAccCheckNfsSettings(defaultNfsSettings),
			},
			{
				Config: testAccNfsSettings(testingNfsSettings),
				Check: resource.ComposeTestCheckFunc(
					testAccCompareNfsSettings(testingNfsSettings),
					testAccCheckNfsSettings(testingNfsSettings)),
			},
		},
	})
}

var defaultNfsSettings = NfsSettings{
	V4Enabled:      true,
	Krb5Enabled:    true,
	AuthSysEnabled: true,
}

var testingNfsSettings = NfsSettings{
	V4Enabled:      false,
	Krb5Enabled:    true,
	AuthSysEnabled: false,
}

func testAccNfsSettings(ns NfsSettings) string {
	return fmt.Sprintf(`
resource "qumulo_nfs_settings" "new_nfs_settings" {
	v4_enabled = %v
	krb5_enabled = %v
	auth_sys_enabled = %v
}
`, ns.V4Enabled, ns.Krb5Enabled, ns.AuthSysEnabled)
}

func testAccCompareNfsSettings(ns NfsSettings) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("qumulo_nfs_settings.my_new_settings", "v4_enabled", fmt.Sprintf("%v", ns.V4Enabled)),
		resource.TestCheckResourceAttr("qumulo_nfs_settings.my_new_settings", "krb5_enabled", fmt.Sprintf("%v", ns.Krb5Enabled)),
		resource.TestCheckResourceAttr("qumulo_nfs_settings.my_new_settings", "auth_sys_enabled", fmt.Sprintf("%v", ns.AuthSysEnabled)))
}

func testAccCheckNfsSettings(ns NfsSettings) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		c := testAccProvider.Meta().(*Client)
		ctx := context.Background()
		settings, err := DoRequest[NfsSettings, NfsSettings](ctx, c, GET, NfsSettingsEndpoint, nil)
		if err != nil {
			return err
		}
		if *settings != ns {
			return fmt.Errorf("Nfs settings mismatch: Expected %v, got %v", ns, settings)
		}
		return nil
	}
}
