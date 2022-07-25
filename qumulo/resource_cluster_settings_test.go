package qumulo

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccChangeClusterName(t *testing.T) {
	fmt.Println("Hi")
	//rName := randSeq(10)
	rName := "InigoMontoya"
	fmt.Println(testAccClusterNameConf(rName))

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccClusterNameConf(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("qumulo_cluster_name.update_name", "name", rName),
				),
			},
		},
	})
}

func testAccClusterNameConf(rName string) string {
	return loginConf + fmt.Sprintf(`
resource "qumulo_cluster_name" "update_name" {
	name = %[1]q
}
`, rName)
}

// func testAccAPIDataSourceConfig_http(rName string) string {
// 	return fmt.Sprintf(`
// resource "aws_apigatewayv2_api" "test" {
//   description   = "test description"
//   name          = %[1]q
//   protocol_type = "HTTP"
//   version       = "v1"
//   cors_configuration {
//     allow_headers = ["Authorization"]
//     allow_methods = ["GET", "put"]
//     allow_origins = ["https://www.example.com"]
//   }
//   tags = {
//     Key1 = "Value1h"
//     Key2 = "Value2h"
//   }
// }
// data "aws_apigatewayv2_api" "test" {
//   api_id = aws_apigatewayv2_api.test.id
// }
// `, rName)
// }

// // TODO: Use this function to make sure API credentials are set
// // (after secure authentication via env vars is implemented)
// func testAccPreCheck(t *testing.T) {
// 	if v := os.Getenv("EXAMPLE_KEY"); v == "" {
// 		t.Fatal("EXAMPLE_KEY must be set for acceptance tests")
// 	}
// 	if v := os.Getenv("EXAMPLE_SECRET"); v == "" {
// 		t.Fatal("EXAMPLE_SECRET must be set for acceptance tests")
// 	}
// }

// func notatest(t *testing.T) {
// 	var widgetBefore, widgetAfter example.Widget
// 	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
// 	fmt.Println("HELLO THERE ERROR ERROR")

// 	resource.Test(t, resource.TestCase{
// 		PreCheck:     func() { testAccPreCheck(t) },
// 		Providers:    testAccProviders,
// 		CheckDestroy: testAccCheckExampleResourceDestroy,
// 		Steps: []resource.TestStep{
// 			{
// 				Config: testAccExampleResource(rName),
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheckExampleResourceExists("example_widget.foo", &widgetBefore),
// 				),
// 			},
// 			{
// 				Config: testAccExampleResource_removedPolicy(rName),
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheckExampleResourceExists("example_widget.foo", &widgetAfter),
// 				),
// 			},
// 		},
// 	})
// }
