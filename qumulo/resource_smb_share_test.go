package qumulo

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// Note - this test assumes that a couple users exist. These should be created by default. These users are:
// User 1:
// name = admin
// User 2:
// uid = 65534

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
				Config: smbShare1,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSmbShareExists(),
				),
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

var share1Name = "ShareForTesting"
var share2Name = "ShareInfinity"

var smbShare1 = fmt.Sprintf(`resource "qumulo_smb_share" "share" {
		share_name = "%v"
		fs_path = "/"
		description = "Surely"
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
				domain = "POSIX_USER"
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
	  }`, share1Name)

var smbShare2 = fmt.Sprintf(`resource "qumulo_smb_share" "share" {
	share_name = "%v"
	fs_path = "/"
	description = "Sharing is caring"
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
		domain = "POSIX_USER"
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
	}`, share2Name)

func testAccCheckSmbShareExists() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources["qumulo_smb_share.share"]
		if !ok {
			return fmt.Errorf("Share not found, %v", rs)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Share ID is not set")
		}

		readSmbShareByIdUri := SmbSharesEndpoint + rs.Primary.ID

		c := testAccProvider.Meta().(*Client)
		ctx := context.Background()
		sh, err := DoRequest[SmbShare, SmbShare](ctx, c, GET, readSmbShareByIdUri, nil)
		if err != nil {
			return err
		}
		localShareName := rs.Primary.Attributes["share_name"]

		if sh.ShareName != localShareName {
			return fmt.Errorf("Share names do not match - Local: %v, Cluster: %v", localShareName, sh.ShareName)
		}
		return nil
	}
}
