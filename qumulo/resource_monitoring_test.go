package qumulo

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccChangeMonitoring(t *testing.T) {
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

var defaultMonitoringConfig = MonitoringSettings{
	Enabled:             false,
	MqHost:              "",
	MqPort:              0,
	MqProxyHost:         "",
	MqProxyPort:         0,
	S3ProxyHost:         "",
	S3ProxyPort:         0,
	S3ProxyDisableHttps: false,
	VpnEnabled:          false,
	VpnHost:             "",
	Period:              0,
}

var testingMonitoringConfig = MonitoringSettings{
	Enabled:             true,
	MqHost:              "missionq.qumulo.com",
	MqPort:              443,
	MqProxyHost:         "missionq.qumulo.com",
	MqProxyPort:         372,
	S3ProxyHost:         "monitor.qumulo.com",
	S3ProxyPort:         444,
	S3ProxyDisableHttps: true,
	VpnEnabled:          true,
	VpnHost:             "ep1.qumulo.com",
	Period:              60,
}

func testAccMonitoringConf(ms MonitoringSettings) string {
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
  `, ms.Enabled, ms.MqHost, ms.MqPort, ms.MqProxyHost, ms.MqProxyPort, ms.S3ProxyHost, ms.S3ProxyPort,
		ms.S3ProxyDisableHttps, ms.VpnEnabled, ms.VpnHost, ms.Period)
}

func testAccCompareMonitoringSetting(ms MonitoringSettings) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("qumulo_monitoring.update_monitoring", "enabled",
			fmt.Sprintf("%v", ms.Enabled)),
		resource.TestCheckResourceAttr("qumulo_monitoring.update_monitoring", "mq_host",
			ms.MqHost),
		resource.TestCheckResourceAttr("qumulo_monitoring.update_monitoring", "mq_port",
			fmt.Sprintf("%v", ms.MqPort)),
		resource.TestCheckResourceAttr("qumulo_monitoring.update_monitoring", "mq_proxy_host",
			ms.MqProxyHost),
		resource.TestCheckResourceAttr("qumulo_monitoring.update_monitoring", "mq_proxy_port",
			fmt.Sprintf("%v", ms.MqProxyPort)),
		resource.TestCheckResourceAttr("qumulo_monitoring.update_monitoring", "s3_proxy_host",
			ms.S3ProxyHost),
		resource.TestCheckResourceAttr("qumulo_monitoring.update_monitoring", "s3_proxy_port",
			fmt.Sprintf("%v", ms.S3ProxyPort)),
		resource.TestCheckResourceAttr("qumulo_monitoring.update_monitoring", "s3_proxy_disable_https",
			fmt.Sprintf("%v", ms.S3ProxyDisableHttps)),
		resource.TestCheckResourceAttr("qumulo_monitoring.update_monitoring", "vpn_enabled",
			fmt.Sprintf("%v", ms.VpnEnabled)),
		resource.TestCheckResourceAttr("qumulo_monitoring.update_monitoring", "vpn_host",
			ms.VpnHost),
		resource.TestCheckResourceAttr("qumulo_monitoring.update_monitoring", "period",
			fmt.Sprintf("%v", ms.Period)),
	)
}

func testAccCheckMonitoringSettings(ms MonitoringSettings) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		c := testAccProvider.Meta().(*Client)
		settings, err := DoRequest[MonitoringSettings, MonitoringSettings](context.Background(), c, GET, MonitoringEndpoint, nil)
		if err != nil {
			return err
		}
		if *settings != ms {
			return fmt.Errorf("Monitoring settings mismatch: Expected %v, got %v", ms, settings)
		}
		return nil
	}
}
