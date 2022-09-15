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

func TestAccChangeSmbServer(t *testing.T) {
	defaultSmbServerConfig := getSmbServerConfig("NONE", []string{}, false, false, "DISABLED",
		false, false)

	testingSmbServerConfig := getSmbServerConfig("PREFERRED", []string{}, true, true, "VISIBLE",
		true, true)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{ // Reset state to default
				Config: testAccSmbServerConfig(defaultSmbServerConfig),

				Check: testAccCheckSmbServerSettings(defaultSmbServerConfig),
			},
			{
				Config: testAccSmbServerConfig(testingSmbServerConfig),
				Check: resource.ComposeTestCheckFunc(
					testAccCompareSmbServerSettings(testingSmbServerConfig),
					testAccCheckSmbServerSettings(testingSmbServerConfig),
				),
			},
		},
	})
}

func getSmbServerConfig(sessionEncryption string, supportedDialects []string, hideSharesFromUnauthorizedUsers bool, hideSharesFromUnauthorizedHosts bool,
	snapshotDirectoryMode string, bypassTraverseChecking bool, signingRequired bool) openapi.V1SmbSettingsGet200Response {
	return openapi.V1SmbSettingsGet200Response{
		SessionEncryption:               &sessionEncryption,
		SupportedDialects:               supportedDialects,
		HideSharesFromUnauthorizedUsers: &hideSharesFromUnauthorizedUsers,
		HideSharesFromUnauthorizedHosts: &hideSharesFromUnauthorizedHosts,
		SnapshotDirectoryMode:           &snapshotDirectoryMode,
		BypassTraverseChecking:          &bypassTraverseChecking,
		SigningRequired:                 &signingRequired,
	}
}
func testAccSmbServerConfig(smb openapi.V1SmbSettingsGet200Response) string {
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
  `, *smb.SessionEncryption, PrintTerraformListFromList(smb.SupportedDialects),
		*smb.HideSharesFromUnauthorizedUsers, *smb.HideSharesFromUnauthorizedHosts, *smb.SnapshotDirectoryMode,
		*smb.BypassTraverseChecking, *smb.SigningRequired)
}

func testAccCompareSmbServerSettings(smb openapi.V1SmbSettingsGet200Response) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("qumulo_smb_server.update_smb", "session_encryption",
			*smb.SessionEncryption),
		resource.TestCheckResourceAttr("qumulo_smb_server.update_smb", "supported_dialects.#",
			fmt.Sprintf("%v", len(smb.SupportedDialects))),
		resource.TestCheckResourceAttr("qumulo_smb_server.update_smb", "hide_shares_from_unauthorized_users",
			fmt.Sprintf("%v", *smb.HideSharesFromUnauthorizedUsers)),
		resource.TestCheckResourceAttr("qumulo_smb_server.update_smb", "hide_shares_from_unauthorized_hosts",
			fmt.Sprintf("%v", *smb.HideSharesFromUnauthorizedUsers)),
		resource.TestCheckResourceAttr("qumulo_smb_server.update_smb", "snapshot_directory_mode",
			*smb.SnapshotDirectoryMode),
		resource.TestCheckResourceAttr("qumulo_smb_server.update_smb", "bypass_traverse_checking",
			fmt.Sprintf("%v", *smb.BypassTraverseChecking)),
		resource.TestCheckResourceAttr("qumulo_smb_server.update_smb", "signing_required",
			fmt.Sprintf("%v", *smb.SigningRequired)),
	)
}

func testAccCheckSmbServerSettings(smb openapi.V1SmbSettingsGet200Response) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		c := testAccProvider.Meta().(*openapi.APIClient)
		settings, _, err := c.SmbApi.V1SmbSettingsGet(context.Background()).Execute()
		if err != nil {
			return err
		}
		if !reflect.DeepEqual(*settings, smb) {
			//lint:ignore ST1005 proper nouns should be capitalized
			return fmt.Errorf("SMB server settings mismatch: Expected %v, got %v", smb, *settings)
		}
		return nil
	}
}
