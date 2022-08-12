package qumulo

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccSetWebUi(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWebUiConfig(testingWebUi),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWebUi(testingWebUi),
					testAccCompareWebUiSettings(testingWebUi),
				),
			},
		},
	})
}

var loginBanner = "SampleBanner"

var testingWebUi = WebUiBody{
	InactivityTimeout: WebUiTimeout{
		Nanoseconds: "900000000000",
	},
	LoginBanner: &loginBanner,
}

func testAccWebUiConfig(req WebUiBody) string {
	return fmt.Sprintf(`
	resource "qumulo_web_ui" "settings" {
		inactivity_timeout {
			nanoseconds = %v
		}
		login_banner = %q
	}
  `, req.InactivityTimeout.Nanoseconds, *testingWebUi.LoginBanner)
}

func testAccCheckWebUi(uiConfig WebUiBody) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		c := testAccProvider.Meta().(*Client)
		ctx := context.Background()
		settings, err := DoRequest[WebUiBody, WebUiBody](ctx, c, GET, WebUiEndpoint, nil)
		if err != nil {
			return err
		}

		if !reflect.DeepEqual(*settings, uiConfig) {
			return fmt.Errorf("web UI settings mismatch: Expected %v, got %v", uiConfig, *settings)
		}
		return nil
	}
}

func testAccCompareWebUiSettings(uiConfig WebUiBody) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("qumulo_web_ui.settings", "inactivity_timeout.0.nanoseconds",
			fmt.Sprintf("%v", uiConfig.InactivityTimeout.Nanoseconds)),
		resource.TestCheckResourceAttr("qumulo_web_ui.settings", "login_banner",
			fmt.Sprintf("%v", *uiConfig.LoginBanner)),
	)
}
