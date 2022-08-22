package qumulo

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"reflect"
	"testing"
)

func TestAccNetworkConfiguration(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				ImportStateId:    "1:1",
				ResourceName:     "qumulo_network_configuration.network_config",
				ImportState:      true,
				Config:           testAccNetworkConfiguration(defaultNetworkConfigurationReq),
				ImportStateCheck: testAccCheckNetworkConfiguration(defaultNetworkConfigurationResp),
			},
		},
	})
}

var defaultNetworkConfigurationReq = NetworkConfigurationRequest{
	Name:             "Default",
	AssignedBy:       "DHCP",
	FloatingIpRanges: []string{},
	DnsServers:       []string{},
	DnsSearchDomains: []string{},
	IpRanges:         []string{},
	Netmask:          "",
	Mtu:              1700,
	VlanId:           0,
	InterfaceId:      "1",
	NetworkId:        "1",
}

var defaultNetworkConfigurationResp = NetworkConfigurationResponse{
	Id:               1,
	Name:             "Default",
	AssignedBy:       "DHCP",
	FloatingIpRanges: []string{},
	DnsServers:       []string{},
	DnsSearchDomains: []string{},
	IpRanges:         []string{},
	Netmask:          "",
	Mtu:              1700,
	VlanId:           0,
}

func testAccNetworkConfiguration(nc NetworkConfigurationRequest) string {
	return fmt.Sprintf(`
resource "qumulo_network_configuration" "network_config" {
}`)
}

func testAccCheckNetworkConfiguration(nc NetworkConfigurationResponse) resource.ImportStateCheckFunc {
	return func(s []*terraform.InstanceState) error {
		c := testAccProvider.Meta().(*Client)
		ctx := context.Background()

		if len(s) != 1 {
			return fmt.Errorf("expected 1 state: %+v", s)
		}
		res := s[0]
		interfaceId := res.Attributes["interface_id"]
		networkId := res.ID
		networkConfigUri := InterfaceConfigurationEndpoint + interfaceId + NetworksEndpointSuffix + networkId

		networkConfig, err := DoRequest[NetworkConfigurationRequest, NetworkConfigurationResponse](ctx, c, GET, networkConfigUri, nil)
		if err != nil {
			return err
		}
		if !reflect.DeepEqual(*networkConfig, nc) {
			return fmt.Errorf("network config mismatch: Expected %v, got %v", nc, *networkConfig)
		}
		return nil
	}
}
