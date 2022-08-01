package qumulo

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAddSmbShare(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{ // Reset state to default
				// Doesn't delete all existing shares, but makes sure any shares set through terraform
				// are gone
				Config: defaultSmbShareConfig,
			},
			{
				Config: smbShare2,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSmbShareExists(),
				),
			},
		},
	})
}

var defaultSmbShareConfig = " "

// var testingSmbShare = SmbShare{
// 	ShareName:              "Testing Share",
// 	FsPath:                 "/",
// 	Description:            "Sure",
// 	Permissions:            1,
// 	NetworkPermissions:     1,
// 	AccessBasedEnumEnabled: true,
// 	RequireEncryption:      false,
// }

func testAccSmbShare() string {
	return `resource "qumulo_smb_share" "share1" {
		share_name = "ShareForTesting"
		fs_path = "/"
		description = "Sure, sure"
		permissions {
		  type = "ALLOWED"
		  trustee {
			  domain = "LOCAL"
			  name = "admin"
		  }
		  rights = ["READ", "WRITE", "CHANGE_PERMISSIONS"]
		}
		permissions {
			type = "DENIED"
			trustee {
				domain = "LOCAL"
				uid = 65534
			}
			rights = ["WRITE"]
		  }
		network_permissions {
		  type = "ALLOWED"
		  address_ranges = []
		  rights = ["READ", "WRITE", "CHANGE_PERMISSIONS"]
		}
		access_based_enumeration_enabled = false
		require_encryption = false
	  }`
}

var smbShare2 = `resource "qumulo_smb_share" "share1" {
	share_name = "ShareForTestingInfinityPlusOne234557"
	fs_path = "/"
	description = "Sure, sure"
	permissions {
	  type = "ALLOWED"
	  trustee {
		name = "admin"
	  }
	  rights = ["READ", "WRITE", "CHANGE_PERMISSIONS"]
	}
	permissions {
	  type = "DENIED"
	  trustee {
		uid = 65534
	  }
	  rights = ["WRITE"]
	}
	network_permissions {
	  type = "ALLOWED"
	  address_ranges = []
	  rights = ["READ", "WRITE", "CHANGE_PERMISSIONS"]
	}
	access_based_enumeration_enabled = false
	require_encryption = false
	}`

func testAccCheckSmbShareExists() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources["qumulo_smb_share.share1"]
		if !ok {
			return fmt.Errorf("Not found: ")
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Share ID is not set")
		}
		// fmt.Printf("Share id is %v", rs.Primary.ID)
		readSmbShareByIdUri := SmbSharesEndpoint + rs.Primary.ID

		c := testAccProvider.Meta().(*Client)
		ctx := context.Background()
		_, err := DoRequest[SmbShare, SmbShare](ctx, c, GET, readSmbShareByIdUri, nil)
		if err != nil {
			return err
		}
		return nil
	}
}
