package qumulo

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccSetTimeConfiguration(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{ // Reset state to default
				Config: testAccTimeConfigurationConfig(defaultTimeConfiguration),
			},
			{
				Config: testAccTimeConfigurationConfig(testingTimeConfiguration),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTimeConfiguration(testingTimeConfiguration),
					testAccCompareTimeConfigurationSettings(testingTimeConfiguration),
				),
			},
			{ // Reset state to default
				Config: testAccTimeConfigurationConfig(defaultTimeConfiguration),
			},
		},
	})
}

var defaultTimeConfiguration = TimeConfigurationBody{
	UseAdForPrimary: false,
	NtpServers:      []string{"0.qumulo.pool.ntp.org", "1.qumulo.pool.ntp.org"},
}

var testingTimeConfiguration = TimeConfigurationBody{
	UseAdForPrimary: true,
	NtpServers:      []string{"0.qumulo.pool.ntp.org"},
}

func testAccTimeConfigurationConfig(req TimeConfigurationBody) string {
	return fmt.Sprintf(`
	resource "qumulo_time_configuration" "time_config" {
		use_ad_for_primary = %v
		ntp_servers = %v
	}
  `, req.UseAdForPrimary, strings.ReplaceAll(fmt.Sprintf("%+q", req.NtpServers), "\" \"", "\", \""))
}

func testAccCheckTimeConfiguration(timeConfig TimeConfigurationBody) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		c := testAccProvider.Meta().(*Client)
		ctx := context.Background()
		settings, err := DoRequest[TimeConfigurationBody, TimeConfigurationBody](ctx, c, GET, TimeConfigurationEndpoint, nil)
		if err != nil {
			return err
		}

		if !reflect.DeepEqual(*settings, timeConfig) {
			return fmt.Errorf("Time configuration settings mismatch: Expected %v, got %v", timeConfig, *settings)
		}
		return nil
	}
}

func testAccCompareTimeConfigurationSettings(timeConfig TimeConfigurationBody) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("qumulo_time_configuration.time_config", "use_ad_for_primary",
			fmt.Sprintf("%v", timeConfig.UseAdForPrimary)),
		resource.TestCheckResourceAttr("qumulo_time_configuration.time_config", "ntp_servers.#",
			fmt.Sprintf("%v", len(timeConfig.NtpServers))),
	)
}
