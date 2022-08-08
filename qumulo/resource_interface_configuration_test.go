package qumulo

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"reflect"
	"strconv"
	"testing"
)

func TestAccTestInterfaceConfiguration(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccInterfaceConfiguration(defaultInterfaceConfiguration),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInterfaceConfiguration(defaultInterfaceConfigurationResp, "qumulo_interface_configuration.interface_config")),
			},
			{
				Config: testAccInterfaceConfiguration(testInterfaceConfiguration),
				Check: resource.ComposeTestCheckFunc(
					testAccCompareInterfaceConfigurations(testInterfaceConfigurationResp),
					testAccCheckInterfaceConfiguration(testInterfaceConfigurationResp, "qumulo_interface_configuration.interface_config")),
			},
		},
	})
}

var defaultInterfaceConfiguration = InterfaceConfigurationRequest{
	Name:           "bond0",
	DefaultGateway: "10.220.0.1",
	BondingMode:    "IEEE_8023AD",
	Mtu:            1500,
	InterfaceId:    "1",
}

var defaultInterfaceConfigurationResp = InterfaceConfigurationResponse{
	Name:           "bond0",
	DefaultGateway: "10.220.0.1",
	BondingMode:    "IEEE_8023AD",
	Mtu:            1500,
}

var testInterfaceConfiguration = InterfaceConfigurationRequest{
	Name:           "bond0",
	DefaultGateway: "10.220.0.2",
	BondingMode:    "IEEE_8023AD",
	Mtu:            1700,
	InterfaceId:    "1",
}

var testInterfaceConfigurationResp = InterfaceConfigurationResponse{
	Name:           "bond0",
	DefaultGateway: "10.220.0.2",
	BondingMode:    "IEEE_8023AD",
	Mtu:            1700,
}

func testAccInterfaceConfiguration(ic InterfaceConfigurationRequest) string {
	return fmt.Sprintf(`
resource "qumulo_interface_configuration" "interface_config" {
  name = %q
  default_gateway = %q
  bonding_mode = %q
  mtu = %v
  interface_id = %q
}`, ic.Name, ic.DefaultGateway, ic.BondingMode, ic.Mtu, ic.InterfaceId)
}

func testAccCompareInterfaceConfigurations(ic InterfaceConfigurationResponse) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("qumulo_interface_configuration.interface_config", "name",
			fmt.Sprintf("%v", ic.Name)),
		resource.TestCheckResourceAttr("qumulo_interface_configuration.interface_config", "default_gateway",
			fmt.Sprintf("%v", ic.DefaultGateway)),
		resource.TestCheckResourceAttr("qumulo_interface_configuration.interface_config", "bonding_mode",
			fmt.Sprintf("%v", ic.BondingMode)),
		resource.TestCheckResourceAttr("qumulo_interface_configuration.interface_config", "mtu",
			fmt.Sprintf("%v", ic.Mtu)),
	)
}

func testAccCheckInterfaceConfiguration(ic InterfaceConfigurationResponse, resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		c := testAccProvider.Meta().(*Client)
		ctx := context.Background()

		res, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("terraform resource not found: %s", resourceName)
		}

		interfaceId := res.Primary.ID
		interfaceConfigUri := InterfaceConfigurationEndpoint + interfaceId

		interfaceConfig, err := DoRequest[InterfaceConfigurationRequest, InterfaceConfigurationResponse](ctx, c, GET, interfaceConfigUri, nil)
		if err != nil {
			return err
		}
		ic.Id, _ = strconv.Atoi(interfaceId)
		if !reflect.DeepEqual(*interfaceConfig, ic) {
			return fmt.Errorf("interface config mismatch: Expected %v, got %v", ic, *interfaceConfig)
		}
		return nil
	}
}
