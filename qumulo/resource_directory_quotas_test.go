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

func TestAccChangeDirectoryQuota(t *testing.T) {
	defaultDirectoryQuota := getAccDirectoryQuota("2", "1000000000")
	testingDirectoryQuota := getAccDirectoryQuota("2", "2000000000")

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

func getAccDirectoryQuota(id string, limit string) openapi.V1FilesQuotasGet200ResponseQuotasInner {
	var reqBody = openapi.V1FilesQuotasGet200ResponseQuotasInner{
		Id:    &id,
		Limit: &limit,
	}
	return reqBody
}

func testAccDirectoryQuotaConfig(req openapi.V1FilesQuotasGet200ResponseQuotasInner) string {
	return fmt.Sprintf(`
	resource "qumulo_directory_quota" "test_quota" {
		directory_id = %v
		limit = %v
	}
  `, *req.Id, *req.Limit)
}

func testAccCheckDirectoryQuota(quota openapi.V1FilesQuotasGet200ResponseQuotasInner) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		c := testAccProvider.Meta().(*openapi.APIClient)

		resp, _, err := c.FilesApi.V1FilesQuotasIdGet(context.Background(), *quota.Id).Execute()
		if err != nil {
			return err
		}

		if !reflect.DeepEqual(*resp, quota) {
			return fmt.Errorf("directory quota settings mismatch: Expected %v, got %v", quota, *resp)
		}
		return nil
	}
}

func testAccCompareDirectoryQuotaSettings(quota openapi.V1FilesQuotasGet200ResponseQuotasInner) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("qumulo_directory_quota.test_quota", "directory_id",
			fmt.Sprintf("%v", *quota.Id)),
		resource.TestCheckResourceAttr("qumulo_directory_quota.test_quota", "limit",
			fmt.Sprintf("%v", *quota.Limit)),
	)
}
