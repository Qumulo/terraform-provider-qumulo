package qumulo

import (
	"context"
	"fmt"
	"reflect"
	"terraform-provider-qumulo/openapi"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccSetTimeConfiguration(t *testing.T) {
	defaultTimeConfiguration := getTimeConfiguration(false, []string{"0.qumulo.pool.ntp.org", "1.qumulo.pool.ntp.org"})
	testingTimeConfiguration := getTimeConfiguration(true, []string{"0.qumulo.pool.ntp.org"})

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

func getTimeConfiguration(useAdForPrimary bool, ntpServers []string) openapi.V1TimeSettingsGet200Response {
	return openapi.V1TimeSettingsGet200Response{
		UseAdForPrimary: &useAdForPrimary,
		NtpServers:      ntpServers,
	}
}

func testAccTimeConfigurationConfig(req openapi.V1TimeSettingsGet200Response) string {
	return fmt.Sprintf(`
	resource "qumulo_time_configuration" "time_config" {
		use_ad_for_primary = %v
		ntp_servers = %v
	}
  `, *req.UseAdForPrimary, PrintTerraformListFromList(req.NtpServers))
}

func testAccCheckTimeConfiguration(timeConfig openapi.V1TimeSettingsGet200Response) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		c := testAccProvider.Meta().(*openapi.APIClient)
		settings, _, err := c.TimeApi.V1TimeSettingsGet(context.Background()).Execute()
		if err != nil {
			return err
		}

		if !reflect.DeepEqual(*settings, timeConfig) {
			return fmt.Errorf("time configuration settings mismatch: Expected %v, got %v", timeConfig, *settings)
		}
		return nil
	}
}

func testAccCompareTimeConfigurationSettings(timeConfig openapi.V1TimeSettingsGet200Response) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("qumulo_time_configuration.time_config", "use_ad_for_primary",
			fmt.Sprintf("%v", *timeConfig.UseAdForPrimary)),
		resource.TestCheckResourceAttr("qumulo_time_configuration.time_config", "ntp_servers.#",
			fmt.Sprintf("%v", len(timeConfig.NtpServers))),
	)
}
