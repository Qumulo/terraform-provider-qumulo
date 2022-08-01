package qumulo

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccChangeSMBServer(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{ // Reset state to default
				Config: testAccSMBServerConfig(defaultSMBServerConfig),
				Check:  testAccCheckSMBServerSettings(defaultSMBServerConfig),
			},
			{
				Config: testAccSMBServerConfig(testingSMBServerConfig),
				Check: resource.ComposeTestCheckFunc(
					testAccCompareSMBServerSettings(testingSMBServerConfig),
					testAccCheckSMBServerSettings(testingSMBServerConfig),
				),
			},
		},
	})
}

var defaultSMBServerConfig = SmbServerBody{
	SessionEncryption:               "NONE",
	SupportedDialects:               nil,
	HideSharesFromUnauthorizedUsers: false,
	HideSharesFromUnauthorizedHosts: false,
	SnapshotDirectoryMode:           "DISABLED",
	BypassTraverseChecking:          false,
	SigningRequired:                 false,
}

var testingSMBServerConfig = SmbServerBody{
	SessionEncryption:               "PREFERRED",
	SupportedDialects:               nil,
	HideSharesFromUnauthorizedUsers: true,
	HideSharesFromUnauthorizedHosts: true,
	SnapshotDirectoryMode:           "VISIBLE",
	BypassTraverseChecking:          true,
	SigningRequired:                 true,
}

func testAccSMBServerConfig(smb SmbServerBody) string {
	return fmt.Sprintf(`
resource "qumulo_smb_server" "update_smb" {
	session_encryption = %q
	supported_dialects = %v
	hide_shares_from_unauthorized_users = %v
	hide_shares_from_unauthorized_hosts = %v
	snapshot_directory_mode = %q
	bypass_traverse_checking = %v
	signing_required = %v
}
  `, smb.SessionEncryption, strings.ReplaceAll(fmt.Sprintf("%+q", smb.SupportedDialects), "\" \"", "\", \""),
		smb.HideSharesFromUnauthorizedUsers, smb.HideSharesFromUnauthorizedHosts, smb.SnapshotDirectoryMode,
		smb.BypassTraverseChecking, smb.SigningRequired)
}

func testAccCompareSMBServerSettings(smb SmbServerBody) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("qumulo_smb_server.update_smb", "session_encryption",
			smb.SessionEncryption),
		resource.TestCheckResourceAttr("qumulo_smb_server.update_smb", "supported_dialects.#",
			fmt.Sprintf("%v", len(smb.SupportedDialects))),
		resource.TestCheckResourceAttr("qumulo_smb_server.update_smb", "hide_shares_from_unauthorized_users",
			fmt.Sprintf("%v", smb.HideSharesFromUnauthorizedUsers)),
		resource.TestCheckResourceAttr("qumulo_smb_server.update_smb", "hide_shares_from_unauthorized_hosts",
			fmt.Sprintf("%v", smb.HideSharesFromUnauthorizedUsers)),
		resource.TestCheckResourceAttr("qumulo_smb_server.update_smb", "snapshot_directory_mode",
			smb.SnapshotDirectoryMode),
		resource.TestCheckResourceAttr("qumulo_smb_server.update_smb", "bypass_traverse_checking",
			fmt.Sprintf("%v", smb.BypassTraverseChecking)),
		resource.TestCheckResourceAttr("qumulo_smb_server.update_smb", "signing_required",
			fmt.Sprintf("%v", smb.SigningRequired)),
	)
}

func testAccCheckSMBServerSettings(smb SmbServerBody) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		c := testAccProvider.Meta().(*Client)
		ctx := context.Background()
		settings, err := DoRequest[SmbServerBody, SmbServerBody](ctx, c, GET, SmbServerEndpoint, nil)
		if err != nil {
			return err
		}
		if !reflect.DeepEqual(settings.SessionEncryption, smb.SessionEncryption) {
			return fmt.Errorf("SMB server settings mismatch: Expected %v, got %v", smb, *settings)
		}
		return nil
	}
}
