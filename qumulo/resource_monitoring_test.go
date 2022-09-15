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

func TestAccChangeMonitoring(t *testing.T) {
	defaultMonitoringConfig := getMonitoringConfig(
		false, "", 0, "", 0, "", 0,
		false, false, "", 0)

	testingMonitoringConfig := getMonitoringConfig(true, "missionq.qumulo.com", 443, "missionq.qumulo.com",
		372,
		"monitor.qumulo.com",
		444,
		true,
		true,
		"ep1.qumulo.com",
		60)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{ // Reset state to default
				Config: testAccMonitoringConf(defaultMonitoringConfig),
				Check:  testAccCheckMonitoringSettings(defaultMonitoringConfig),
			},
			{
				Config: testAccMonitoringConf(testingMonitoringConfig),
				Check: resource.ComposeTestCheckFunc(
					testAccCompareMonitoringSetting(testingMonitoringConfig),
					testAccCheckMonitoringSettings(testingMonitoringConfig),
				),
			},
		},
	})
}

func getMonitoringConfig(enabled bool, mqHost string, mqPort float32, mqProxyHost string, mqProxyPort float32,
	s3ProxyHost string, s3ProxyPort float32, s3ProxyDisableHttps bool, vpnEnabled bool, vpnHost string, period float32) openapi.V1SupportSettingsGet200Response {
	return openapi.V1SupportSettingsGet200Response{
		Enabled:             &enabled,
		MqHost:              &mqHost,
		MqPort:              &mqPort,
		MqProxyHost:         &mqProxyHost,
		MqProxyPort:         &mqProxyPort,
		S3ProxyHost:         &s3ProxyHost,
		S3ProxyPort:         &s3ProxyPort,
		S3ProxyDisableHttps: &s3ProxyDisableHttps,
		VpnEnabled:          &vpnEnabled,
		VpnHost:             &vpnHost,
		Period:              &period,
	}
}

func testAccMonitoringConf(ms openapi.V1SupportSettingsGet200Response) string {
	return fmt.Sprintf(`
resource "qumulo_monitoring" "update_monitoring" {
	enabled = %v
	mq_host = %q
	mq_port = %v
	mq_proxy_host = %q
	mq_proxy_port = %v
	s3_proxy_host = %q
	s3_proxy_port = %v
	s3_proxy_disable_https = %v
	vpn_enabled = %v
	vpn_host = %q
	period = %v
  }
  `, *ms.Enabled, *ms.MqHost, *ms.MqPort, *ms.MqProxyHost, *ms.MqProxyPort, *ms.S3ProxyHost, *ms.S3ProxyPort,
		*ms.S3ProxyDisableHttps, *ms.VpnEnabled, *ms.VpnHost, *ms.Period)
}

func testAccCompareMonitoringSetting(ms openapi.V1SupportSettingsGet200Response) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("qumulo_monitoring.update_monitoring", "enabled",
			fmt.Sprintf("%v", *ms.Enabled)),
		resource.TestCheckResourceAttr("qumulo_monitoring.update_monitoring", "mq_host",
			*ms.MqHost),
		resource.TestCheckResourceAttr("qumulo_monitoring.update_monitoring", "mq_port",
			fmt.Sprintf("%v", *ms.MqPort)),
		resource.TestCheckResourceAttr("qumulo_monitoring.update_monitoring", "mq_proxy_host",
			*ms.MqProxyHost),
		resource.TestCheckResourceAttr("qumulo_monitoring.update_monitoring", "mq_proxy_port",
			fmt.Sprintf("%v", *ms.MqProxyPort)),
		resource.TestCheckResourceAttr("qumulo_monitoring.update_monitoring", "s3_proxy_host",
			*ms.S3ProxyHost),
		resource.TestCheckResourceAttr("qumulo_monitoring.update_monitoring", "s3_proxy_port",
			fmt.Sprintf("%v", *ms.S3ProxyPort)),
		resource.TestCheckResourceAttr("qumulo_monitoring.update_monitoring", "s3_proxy_disable_https",
			fmt.Sprintf("%v", *ms.S3ProxyDisableHttps)),
		resource.TestCheckResourceAttr("qumulo_monitoring.update_monitoring", "vpn_enabled",
			fmt.Sprintf("%v", *ms.VpnEnabled)),
		resource.TestCheckResourceAttr("qumulo_monitoring.update_monitoring", "vpn_host",
			*ms.VpnHost),
		resource.TestCheckResourceAttr("qumulo_monitoring.update_monitoring", "period",
			fmt.Sprintf("%v", *ms.Period)),
	)
}

func testAccCheckMonitoringSettings(ms openapi.V1SupportSettingsGet200Response) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		c := testAccProvider.Meta().(*openapi.APIClient)
		settings, _, err := c.SupportApi.V1SupportSettingsGet(context.Background()).Execute()
		if err != nil {
			return err
		}
		if !reflect.DeepEqual(*settings, ms) {
			return fmt.Errorf("monitoring settings mismatch: Expected %v, got %v", ms, *settings)
		}
		return nil
	}
}
