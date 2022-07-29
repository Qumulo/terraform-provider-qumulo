package qumulo

import (
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

var defaultMonitoringConfig = MonitorSettings{
	Enabled:             false,
	MQHost:              "",
	MQPort:              0,
	MQProxyHost:         "",
	MQProxyPort:         0,
	S3ProxyHost:         "",
	S3ProxyPort:         0,
	S3ProxyDisableHTTPS: false,
	VPNEnabled:          false,
	VPNHost:             "",
	Period:              0,
}

var testingMonitoringConfig = MonitorSettings{
	Enabled:             true,
	MQHost:              "missionq.qumulo.com",
	MQPort:              443,
	MQProxyHost:         "missionq.qumulo.com",
	MQProxyPort:         372,
	S3ProxyHost:         "monitor.qumulo.com",
	S3ProxyPort:         444,
	S3ProxyDisableHTTPS: true,
	VPNEnabled:          true,
	VPNHost:             "ep1.qumulo.com",
	Period:              60,
}

func testAccMonitoringConf(ms MonitorSettings) string {
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
  `, ms.Enabled, ms.MQHost, ms.MQPort, ms.MQProxyHost, ms.MQProxyPort, ms.S3ProxyHost, ms.S3ProxyPort,
		ms.S3ProxyDisableHTTPS, ms.VPNEnabled, ms.VPNHost, ms.Period)
}

func testAccCompareMonitoringSetting(ms MonitorSettings) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("qumulo_monitoring.update_monitoring", "enabled",
			fmt.Sprintf("%v", ms.Enabled)),
		resource.TestCheckResourceAttr("qumulo_monitoring.update_monitoring", "mq_host",
			ms.MQHost),
		resource.TestCheckResourceAttr("qumulo_monitoring.update_monitoring", "mq_port",
			fmt.Sprintf("%v", ms.MQPort)),
		resource.TestCheckResourceAttr("qumulo_monitoring.update_monitoring", "mq_proxy_host",
			ms.MQProxyHost),
		resource.TestCheckResourceAttr("qumulo_monitoring.update_monitoring", "mq_proxy_port",
			fmt.Sprintf("%v", ms.MQProxyPort)),
		resource.TestCheckResourceAttr("qumulo_monitoring.update_monitoring", "s3_proxy_host",
			ms.S3ProxyHost),
		resource.TestCheckResourceAttr("qumulo_monitoring.update_monitoring", "s3_proxy_port",
			fmt.Sprintf("%v", ms.S3ProxyPort)),
		resource.TestCheckResourceAttr("qumulo_monitoring.update_monitoring", "s3_proxy_disable_https",
			fmt.Sprintf("%v", ms.S3ProxyDisableHTTPS)),
		resource.TestCheckResourceAttr("qumulo_monitoring.update_monitoring", "vpn_enabled",
			fmt.Sprintf("%v", ms.VPNEnabled)),
		resource.TestCheckResourceAttr("qumulo_monitoring.update_monitoring", "vpn_host",
			ms.VPNHost),
		resource.TestCheckResourceAttr("qumulo_monitoring.update_monitoring", "period",
			fmt.Sprintf("%v", ms.Period)),
	)
}

func testAccCheckMonitoringSettings(ms MonitorSettings) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		c := testAccProvider.Meta().(*Client)
		settings, err := DoRequest[MonitorSettings, MonitorSettings](c, GET, MonitorEndpoint, nil)
		if err != nil {
			return err
		}
		if *settings != ms {
			return fmt.Errorf("Monitoring settings mismatch: Expected %v, got %v", ms, settings)
		}
		return nil
	}
}
