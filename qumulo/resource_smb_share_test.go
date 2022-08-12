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
			{
				Config: smbShare1,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckShare(share1),
					testAccCompareShare(share1),
				),
			},
			{
				Config: smbShare1Updated,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckShare(share1Updated),
					testAccCompareShare(share1Updated),
				),
			},
		},
	})
}

var share1 = SmbShare{
	ShareName:   "ShareForTesting",
	FsPath:      "/",
	Description: "Surely not",
	Permissions: []SmbPermission{
		SmbPermission{
			Type: "ALLOWED",
			Trustee: SmbTrustee{
				Domain: "LOCAL",
				Name:   "user1",
			},
			Rights: []string{"READ", "WRITE", "CHANGE_PERMISSIONS"},
		},
		SmbPermission{
			Type: "DENIED",
			Trustee: SmbTrustee{
				Name: "user2",
			},
			Rights: []string{"WRITE"},
		},
	},
	NetworkPermissions: []SmbNetworkPermission{
		SmbNetworkPermission{
			Type:          "ALLOWED",
			AddressRanges: []string{},
			Rights:        []string{"READ", "WRITE", "CHANGE_PERMISSIONS"},
		},
	},
	AccessBasedEnumEnabled: false,
	RequireEncryption:      false,
}

var smbShare1 = fmt.Sprintf(`
resource "qumulo_local_user" "user1" {
	name = "user1"
	primary_group = 514
	password = "Test1234"
}

resource "qumulo_local_user" "user2" {
	name = "user2"
	primary_group = 514
	password = "Test1234"
}

resource "qumulo_smb_share" "share" {
		share_name = %q
		fs_path = %q
		description = %q
		permissions {
		  type = %q
		  trustee {
			name = %q
		  }
		  rights = %v
		}
		permissions {
			type = %q
			trustee {
				name = %q
			}
			rights = %v
		  }
		network_permissions {
		  type = %q
		  address_ranges = %v
		  rights = %v
		}
		access_based_enumeration_enabled = %v
		require_encryption = %v
	  }`, share1.ShareName, share1.FsPath, share1.Description, share1.Permissions[0].Type, share1.Permissions[0].Trustee.Name,
	PrintTerraformListFromList(share1.Permissions[0].Rights), share1.Permissions[1].Type, share1.Permissions[1].Trustee.Name,
	PrintTerraformListFromList(share1.Permissions[1].Rights), share1.NetworkPermissions[0].Type, "[]",
	PrintTerraformListFromList(share1.NetworkPermissions[0].Rights), share1.AccessBasedEnumEnabled, share1.RequireEncryption)

var share1Updated = SmbShare{
	ShareName:   "ShareForTesting",
	FsPath:      "/",
	Description: "Sharing is caring",
	Permissions: []SmbPermission{
		SmbPermission{
			Type: "ALLOWED",
			Trustee: SmbTrustee{
				Name: "user1",
			},
			Rights: []string{"READ", "WRITE"},
		},
		SmbPermission{
			Type: "DENIED",
			Trustee: SmbTrustee{
				Domain: "LOCAL",
				Name:   "user2",
			},
			Rights: []string{"WRITE"},
		},
	},
	NetworkPermissions: []SmbNetworkPermission{
		SmbNetworkPermission{
			Type:          "ALLOWED",
			AddressRanges: []string{},
			Rights:        []string{"READ", "WRITE", "CHANGE_PERMISSIONS"},
		},
	},
	AccessBasedEnumEnabled: true,
	RequireEncryption:      false,
}
var smbShare1Updated = fmt.Sprintf(`
resource "qumulo_local_user" "user1" {
	name = "user1"
	primary_group = 514
	password = "Test1234"
}

resource "qumulo_local_user" "user2" {
	name = "user2"
	primary_group = 514
	password = "Test1234"
}

resource "qumulo_smb_share" "share" {
		share_name = %q
		fs_path = %q
		description = %q
		permissions {
		  type = %q
		  trustee {
			name = %q
		  }
		  rights = %v
		}
		permissions {
			type = %q
			trustee {
				name = %q
			}
			rights = %v
		  }
		network_permissions {
		  type = %q
		  address_ranges = %v
		  rights = %v
		}
		access_based_enumeration_enabled = %v
		require_encryption = %v
	  }`, share1Updated.ShareName, share1Updated.FsPath, share1Updated.Description, share1Updated.Permissions[0].Type,
	share1Updated.Permissions[0].Trustee.Name, PrintTerraformListFromList(share1Updated.Permissions[0].Rights),
	share1Updated.Permissions[1].Type, share1Updated.Permissions[1].Trustee.Name,
	PrintTerraformListFromList(share1Updated.Permissions[1].Rights), share1Updated.NetworkPermissions[0].Type, "[]",
	PrintTerraformListFromList(share1Updated.NetworkPermissions[0].Rights), share1Updated.AccessBasedEnumEnabled,
	share1Updated.RequireEncryption)

func testAccCompareShare(share SmbShare) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("qumulo_smb_share.share", "share_name",
			share.ShareName),
		resource.TestCheckResourceAttr("qumulo_smb_share.share", "description",
			share.Description),
		resource.TestCheckResourceAttr("qumulo_smb_share.share", "fs_path",
			share.FsPath),
		resource.TestCheckResourceAttr("qumulo_smb_share.share", "permissions.#",
			fmt.Sprintf("%v", len(share.Permissions))),
		resource.TestCheckResourceAttr("qumulo_smb_share.share", "network_permissions.#",
			fmt.Sprintf("%v", len(share.NetworkPermissions))),
		resource.TestCheckResourceAttr("qumulo_smb_share.share", "access_based_enumeration_enabled",
			fmt.Sprintf("%v", share.AccessBasedEnumEnabled)),
		resource.TestCheckResourceAttr("qumulo_smb_share.share", "require_encryption",
			fmt.Sprintf("%v", share.RequireEncryption)),
	)
}

func testAccCheckShare(share SmbShare) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources["qumulo_smb_share.share"]
		if !ok {
			//lint:ignore ST1005 proper nouns should be capitalized
			return fmt.Errorf("SMB share not found, %v", rs)
		}

		if rs.Primary.ID == "" {
			//lint:ignore ST1005 proper nouns should be capitalized
			return fmt.Errorf("SMB share ID is not set")
		}

		readSmbShareByIdUri := SmbSharesEndpoint + rs.Primary.ID

		c := testAccProvider.Meta().(*Client)
		ctx := context.Background()
		sh, err := DoRequest[SmbShare, SmbShare](ctx, c, GET, readSmbShareByIdUri, nil)
		if err != nil {
			return err
		}

		if sh.ShareName != share.ShareName {
			//lint:ignore ST1005 proper nouns should be capitalized
			return fmt.Errorf("SMB share name mismatch: Expected %v, got %v", share.ShareName, sh.ShareName)
		}
		return nil
	}
}
