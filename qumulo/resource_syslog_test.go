package qumulo

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccCreateSyslogAuditLog(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{ // Reset state to default
				Config: testAccSyslogConfig(defaultSyslogConfig),
				Check: resource.ComposeTestCheckFunc(
					testAccCompareSyslogResource(defaultSyslogConfig),
					testAccCheckSyslogSettings(defaultSyslogConfig),
				),
			},
			{
				Config: testAccSyslogConfig(testSyslogConfig),
				Check: resource.ComposeTestCheckFunc(
					testAccCompareSyslogResource(testSyslogConfig),
					testAccCheckSyslogSettings(testSyslogConfig),
				),
			},
			{ // Reset state to default
				Config: testAccSyslogConfig(defaultSyslogConfig),
				Check: resource.ComposeTestCheckFunc(
					testAccCompareSyslogResource(defaultSyslogConfig),
					testAccCheckSyslogSettings(defaultSyslogConfig),
				),
			},
		},
	})
}

var defaultSyslogConfig = SyslogConfigBody{
	Enabled:       false,
	ServerAddress: "",
	ServerPort:    0,
}

var testSyslogConfig = SyslogConfigBody{
	Enabled:       true,
	ServerAddress: "127.0.0.1",
	ServerPort:    13337,
}

func testAccSyslogConfig(settings SyslogConfigBody) string {
	return fmt.Sprintf(`
resource "qumulo_syslog" "test_syslog_settings" {
	enabled = %v
	server_address = %q
	server_port = %v
}
  `, settings.Enabled, settings.ServerAddress, settings.ServerPort)
}

func testAccCompareSyslogResource(settings SyslogConfigBody) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("qumulo_syslog.test_syslog_settings", "enabled",
			fmt.Sprintf("%v", settings.Enabled)),
		resource.TestCheckResourceAttr("qumulo_syslog.test_syslog_settings", "server_address",
			settings.ServerAddress),
		resource.TestCheckResourceAttr("qumulo_syslog.test_syslog_settings", "server_port",
			fmt.Sprintf("%v", settings.ServerPort)),
	)
}

func testAccCheckSyslogSettings(settings SyslogConfigBody) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		c := testAccProvider.Meta().(*Client)
		ctx := context.Background()

		remoteSettings, err := DoRequest[SyslogConfigBody, SyslogConfigBody](ctx, c, GET, SyslogConfigEndpoint, nil)
		if err != nil {
			return err
		}

		if !(reflect.DeepEqual(*remoteSettings, settings)) {
			return fmt.Errorf("syslog configuration mismatch: Expected %v, got %v", settings, *remoteSettings)
		}

		return nil
	}
}
