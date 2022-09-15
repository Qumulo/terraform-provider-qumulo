package qumulo

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"reflect"
	"terraform-provider-qumulo/openapi"
	"testing"
)

func TestAccChangeNfsSettings(t *testing.T) {
	defaultNfsSettings := getNfsSettings(false, true, true, true)
	testingNfsSettings := getNfsSettings(false, true, true, false)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{ // Reset state to default
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

func getNfsSettings(v4Enabled bool, krb5Enabled bool, krb5pEnabled bool, authSysEnabled bool) openapi.V2NfsSettingsGet200Response {
	return openapi.V2NfsSettingsGet200Response{
		V4Enabled:      &v4Enabled,
		Krb5Enabled:    &krb5Enabled,
		Krb5pEnabled:   &krb5pEnabled,
		AuthSysEnabled: &authSysEnabled,
	}
}

func testAccNfsSettings(ns openapi.V2NfsSettingsGet200Response) string {
	return fmt.Sprintf(`
resource "qumulo_nfs_settings" "new_nfs_settings" {
	v4_enabled = %v
	krb5_enabled = %v
	krb5p_enabled = %v
	auth_sys_enabled = %v
}
`, *ns.V4Enabled, *ns.Krb5Enabled, *ns.Krb5pEnabled, *ns.AuthSysEnabled)
}

func testAccCompareNfsSettings(ns openapi.V2NfsSettingsGet200Response) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("qumulo_nfs_settings.new_nfs_settings", "v4_enabled", fmt.Sprintf("%v", *ns.V4Enabled)),
		resource.TestCheckResourceAttr("qumulo_nfs_settings.new_nfs_settings", "krb5_enabled", fmt.Sprintf("%v", *ns.Krb5Enabled)),
		resource.TestCheckResourceAttr("qumulo_nfs_settings.new_nfs_settings", "krb5p_enabled", fmt.Sprintf("%v", *ns.Krb5pEnabled)),
		resource.TestCheckResourceAttr("qumulo_nfs_settings.new_nfs_settings", "auth_sys_enabled", fmt.Sprintf("%v", *ns.AuthSysEnabled)))
}

func testAccCheckNfsSettings(ns openapi.V2NfsSettingsGet200Response) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		c := testAccProvider.Meta().(*openapi.APIClient)
		settings, _, err := c.NfsApi.V2NfsSettingsGet(context.Background()).Execute()
		if err != nil {
			return err
		}
		if !reflect.DeepEqual(*settings, ns) {
			return fmt.Errorf("NFS settings mismatch: Expected %v, got %v", ns, *settings)
		}
		return nil
	}
}
