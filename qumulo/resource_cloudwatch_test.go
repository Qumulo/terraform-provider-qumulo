package qumulo

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccCreateCloudWatchAuditLog(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{ // Reset state to default
				Config: testAccCloudWatchConfig(defaultCloudWatchConfig),
				Check: resource.ComposeTestCheckFunc(
					testAccCompareCloudWatchResource(defaultCloudWatchConfig),
					testAccCheckCloudWatchSettings(defaultCloudWatchConfig),
				),
			},
			{
				Config: testAccCloudWatchConfig(testCloudWatchConfig),
				Check: resource.ComposeTestCheckFunc(
					testAccCompareCloudWatchResource(testCloudWatchConfig),
					testAccCheckCloudWatchSettings(testCloudWatchConfig),
				),
			},
			{ // Reset state to default
				Config: testAccCloudWatchConfig(defaultCloudWatchConfig),
				Check: resource.ComposeTestCheckFunc(
					testAccCompareCloudWatchResource(defaultCloudWatchConfig),
					testAccCheckCloudWatchSettings(defaultCloudWatchConfig),
				),
			},
		},
	})
}

var defaultCloudWatchConfig = CloudWatchConfigBody{
	Enabled:      false,
	LogGroupName: "",
	Region:       "",
}

var testCloudWatchConfig = CloudWatchConfigBody{
	Enabled:      true,
	LogGroupName: "test_group",
	Region:       "test_region",
}

func testAccCloudWatchConfig(settings CloudWatchConfigBody) string {
	return fmt.Sprintf(`
resource "qumulo_cloudwatch" "test_cloudwatch_settings" {
	enabled = %v
	log_group_name = %q
	region = %q
}
  `, settings.Enabled, settings.LogGroupName, settings.Region)
}

func testAccCompareCloudWatchResource(settings CloudWatchConfigBody) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("qumulo_cloudwatch.test_cloudwatch_settings", "enabled",
			fmt.Sprintf("%v", settings.Enabled)),
		resource.TestCheckResourceAttr("qumulo_cloudwatch.test_cloudwatch_settings", "log_group_name",
			settings.LogGroupName),
		resource.TestCheckResourceAttr("qumulo_cloudwatch.test_cloudwatch_settings", "region",
			settings.Region),
	)
}

func testAccCheckCloudWatchSettings(settings CloudWatchConfigBody) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		c := testAccProvider.Meta().(*Client)
		ctx := context.Background()

		remoteSettings, err := DoRequest[CloudWatchConfigBody, CloudWatchConfigBody](ctx, c, GET, CloudWatchConfigEndpoint, nil)
		if err != nil {
			return err
		}

		if !(reflect.DeepEqual(*remoteSettings, settings)) {
			return fmt.Errorf("CloudWatch configuration mismatch: Expected %v, got %v", settings, *remoteSettings)
		}

		return nil
	}
}
