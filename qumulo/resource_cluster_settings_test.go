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

func TestAccChangeClusterName(t *testing.T) {
	defaultName := "InigoMontoya"
	rName := "qfsd"
	rName2 := "Buttercup"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{ // Reset state to default
				Config: testAccClusterNameConf(defaultName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClusterName(defaultName),
				),
			},
			{
				Config: testAccClusterNameConf(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("qumulo_cluster_name.update_name", "cluster_name", rName),
					testAccCheckClusterName(rName),
				),
			},
			{
				Config: testAccClusterNameConf(rName2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("qumulo_cluster_name.update_name", "cluster_name", rName2),
					testAccCheckClusterName(rName2),
				),
			},
			{
				ResourceName:      "qumulo_cluster_name.update_name",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccClusterNameConf(name string) string {
	return fmt.Sprintf(`resource "qumulo_cluster_name" "update_name" {
	cluster_name = %q
}
`, name)
}

func testAccCheckClusterName(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		c := testAccProvider.Meta().(*openapi.APIClient)
		res, _, err := c.ClusterApi.V1ClusterSettingsGet(context.Background()).Execute()
		if err != nil {
			return err
		}
		if !reflect.DeepEqual(*res.ClusterName, name) {
			return fmt.Errorf("cluster name mismatch: Expected %s, got %s", name, *res.ClusterName)
		}
		return nil
	}
}
