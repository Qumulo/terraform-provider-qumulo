package qumulo

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccChangeClusterName(t *testing.T) {
	defaultName := "qfsd"
	rName := "InigoMontoya"
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
					resource.TestCheckResourceAttr("qumulo_cluster_name.update_name", "name", rName),
					testAccCheckClusterName(rName),
				),
			},
			{
				Config: testAccClusterNameConf(rName2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("qumulo_cluster_name.update_name", "name", rName2),
					testAccCheckClusterName(rName2),
				),
			},
		},
	})
}

func testAccClusterNameConf(name string) string {
	return fmt.Sprintf(`
resource "qumulo_cluster_name" "update_name" {
	name = %q
}
`, name)
}

func testAccCheckClusterName(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		c := testAccProvider.Meta().(*Client)
		cs, err := DoRequest[ClusterSettingsBody, ClusterSettingsBody](context.Background(), c, GET, ClusterSettingsEndpoint, nil)
		if err != nil {
			return err
		}
		if cs.ClusterName != name {
			fmt.Println(cs.ClusterName)
			return fmt.Errorf("Cluster name is not updated: Expected %s, got %s", name, cs.ClusterName)
		}
		return nil
	}
}
