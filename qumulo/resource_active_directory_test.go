package qumulo

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// TODO write test steps

func TestAccChangeAdSettings(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// { // Reset state to default
			// 	Config: testAccActiveDirectoryConfigFull(emptyActiveDirectoryConfig),
			// 	// Check: resource.ComposeTestCheckFunc(
			// 	// 	testAccCheckClusterName(defaultName),
			// 	// ),
			// },
			{
				Config: testAccActiveDirectoryConfigFull(defaultActiveDirectoryConfig),
				Check: resource.ComposeTestCheckFunc(
					testAccCompareActiveDirectorySettings(defaultActiveDirectoryConfig),
					testAccCheckActiveDirectorySettings(*defaultActiveDirectoryConfig.Settings),
				),
			},
			// {
			// 	Config: testAccClusterNameConf(rName2),
			// 	Check: resource.ComposeTestCheckFunc(
			// 		resource.TestCheckResourceAttr("qumulo_cluster_name.update_name", "name", rName2),
			// 		testAccCheckClusterName(rName2),
			// 	),
			// },
		},
	})
}

// Active Directory configurations

// TODO create all relevant combinations for test

// var emptyActiveDirectoryConfig = ActiveDirectoryRequest{
// 	Settings:      &defaultActiveDirectorySettingsConfig,
// 	JoinSettings:  &ActiveDirectoryJoinRequest{},
// 	UsageSettings: &ActiveDirectoryUsageSettingsRequest{},
// }

var defaultActiveDirectoryConfig = ActiveDirectoryRequest{
	Settings:      &defaultActiveDirectorySettingsConfig,
	JoinSettings:  &testingActiveDirectoryJoinSettingsConfigFull,
	UsageSettings: &ActiveDirectoryUsageSettingsRequest{},
}

// Active Directory Settings configurations

var defaultActiveDirectorySettingsConfig = ActiveDirectorySettingsBody{
	Signing: "WANT_SIGNING",
	Sealing: "WANT_SEALING",
	Crypto:  "WANT_AES",
}

var testingActiveDirectorySettingsConfigFull = ActiveDirectorySettingsBody{
	Signing: "NEED_SIGNING",
	Sealing: "WANT_SEALING",
	Crypto:  "NEED_AES",
}

var testingActiveDirectorySettingsConfigPartial = ActiveDirectorySettingsBody{
	Signing: "NEED_SIGNING",
	Crypto:  "NEED_AES",
}

// Active Directory Join Settings configurations

var testingActiveDirectoryJoinSettingsConfigFull = ActiveDirectoryJoinRequest{
	Domain:               "ad.eng.qumulo.com",
	DomainNetBios:        "AD",
	User:                 "Administrator",
	Password:             "a",
	Ou:                   "",
	UseAdPosixAttributes: false,
	BaseDn:               "CN=Users,DC=ad,DC=eng,DC=qumulo,DC=com",
}

var testingActiveDirectoryJoinSettingsConfigFullAlternate = ActiveDirectoryJoinRequest{
	Domain:               "ad.eng.qumulo.com",
	DomainNetBios:        "AD",
	User:                 "Administrator",
	Password:             "a",
	Ou:                   "",
	UseAdPosixAttributes: true,
	BaseDn:               "",
}

var testingActiveDirectoryJoinSettingsConfigPartial = ActiveDirectoryJoinRequest{
	Domain:   "ad.eng.qumulo.com",
	User:     "Administrator",
	Password: "a",
}

func testAccActiveDirectoryConfigFull(req ActiveDirectoryRequest) string {
	return fmt.Sprintf(`
	resource "qumulo_ad_settings" "ad_settings" {
		signing = %q
		sealing = %q
		crypto = %q
		domain = %q
		domain_netbios = %q
		ad_username = %q
		ad_password = %q
		ou = %q
		use_ad_posix_attributes = %v
		base_dn = %q
	}
	`, req.Settings.Signing, req.Settings.Sealing, req.Settings.Crypto, req.JoinSettings.Domain, req.JoinSettings.DomainNetBios, req.JoinSettings.User,
		req.JoinSettings.Password, req.JoinSettings.Ou, req.JoinSettings.UseAdPosixAttributes, req.JoinSettings.BaseDn)
}

// TODO write actual test

func testAccCompareActiveDirectorySettings(adRequest ActiveDirectoryRequest) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("qumulo_ad_settings.ad_settings", "signing",
			adRequest.Settings.Signing),
		resource.TestCheckResourceAttr("qumulo_ad_settings.ad_settings", "sealing",
			adRequest.Settings.Sealing),
		resource.TestCheckResourceAttr("qumulo_ad_settings.ad_settings", "crypto",
			adRequest.Settings.Crypto),
		resource.TestCheckResourceAttr("qumulo_ad_settings.ad_settings", "domain",
			adRequest.JoinSettings.Domain),
		resource.TestCheckResourceAttr("qumulo_ad_settings.ad_settings", "domain_netbios",
			adRequest.JoinSettings.DomainNetBios),
		resource.TestCheckResourceAttr("qumulo_ad_settings.ad_settings", "ad_username",
			adRequest.JoinSettings.User),
		resource.TestCheckResourceAttr("qumulo_ad_settings.ad_settings", "ad_password",
			adRequest.JoinSettings.Password),
		resource.TestCheckResourceAttr("qumulo_ad_settings.ad_settings", "ou",
			adRequest.JoinSettings.Ou),
		resource.TestCheckResourceAttr("qumulo_ad_settings.ad_settings", "use_ad_posix_attributes",
			fmt.Sprintf("%v", adRequest.JoinSettings.UseAdPosixAttributes)),
		resource.TestCheckResourceAttr("qumulo_ad_settings.ad_settings", "base_dn",
			adRequest.JoinSettings.BaseDn),
	)
}

func testAccCheckActiveDirectorySettings(adSettingsRequest ActiveDirectorySettingsBody) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		c := testAccProvider.Meta().(*Client)
		ctx := context.Background()
		adSettings, err := DoRequest[ActiveDirectorySettingsBody, ActiveDirectorySettingsBody](ctx, c, GET, AdSettingsEndpoint, nil)
		if err != nil {
			return err
		}

		if adSettings.Sealing != adSettingsRequest.Sealing {
			return fmt.Errorf("Active Directory settings mismatch: Expected %v, got %v", adSettingsRequest, adSettings)
		}
		return nil
	}
}
