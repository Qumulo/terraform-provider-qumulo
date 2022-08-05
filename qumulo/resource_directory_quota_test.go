package qumulo

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccChangeDirectoryQuota(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDirectoryQuotaConfig(defaultDirectoryQuota),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDirectoryQuota(defaultDirectoryQuota),
					testAccCompareDirectoryQuotaSettings(defaultDirectoryQuota),
				),
			},
			{
				Config: testAccDirectoryQuotaConfig(testingDirectoryQuota),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDirectoryQuota(testingDirectoryQuota),
					testAccCompareDirectoryQuotaSettings(testingDirectoryQuota),
				),
			},
		},
	})
}

var defaultDirectoryQuota = DirectoryQuotaBody{
	Id:    "2",
	Limit: "1000000000",
}

var testingDirectoryQuota = DirectoryQuotaBody{
	Id:    "2",
	Limit: "2000000000",
}

func testAccDirectoryQuotaConfig(req DirectoryQuotaBody) string {
	return fmt.Sprintf(`
	resource "qumulo_directory_quota" "test_quota" {
		directory_id = %v
		limit = %v
	}
  `, req.Id, req.Limit)
}

func testAccCheckDirectoryQuota(quota DirectoryQuotaBody) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		c := testAccProvider.Meta().(*Client)
		ctx := context.Background()
		quotaUrl := DirectoryQuotaEndpoint + quota.Id

		settings, err := DoRequest[DirectoryQuotaEmptyBody, DirectoryQuotaBody](ctx, c, GET, quotaUrl, nil)
		if err != nil {
			return err
		}

		if !reflect.DeepEqual(*settings, quota) {
			return fmt.Errorf("Directory quota settings mismatch: Expected %v, got %v", quota, *settings)
		}
		return nil
	}
}

func testAccCompareDirectoryQuotaSettings(quota DirectoryQuotaBody) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("qumulo_directory_quota.test_quota", "directory_id",
			fmt.Sprintf("%v", quota.Id)),
		resource.TestCheckResourceAttr("qumulo_directory_quota.test_quota", "limit",
			fmt.Sprintf("%v", quota.Limit)),
	)
}
