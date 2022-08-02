package qumulo

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// TODO write test steps

func TestAccJoinActiveDirectory(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccActiveDirectoryConfigFull(defaultActiveDirectoryConfig),
				Check: resource.ComposeTestCheckFunc(
					testAccCompareActiveDirectorySettings(defaultActiveDirectoryConfig),
					testAccCheckActiveDirectorySettings(*defaultActiveDirectoryConfig.Settings),
					testAccCheckActiveDirectoryStatus(*defaultActiveDirectoryConfig.JoinSettings),
				),
			},
		},
	})
}

func TestAccChangeActiveDirectorySettings(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccActiveDirectoryConfigFull(defaultActiveDirectoryConfig),
				Check: resource.ComposeTestCheckFunc(
					testAccCompareActiveDirectorySettings(defaultActiveDirectoryConfig),
					testAccCheckActiveDirectorySettings(*defaultActiveDirectoryConfig.Settings),
					testAccCheckActiveDirectoryStatus(*defaultActiveDirectoryConfig.JoinSettings),
				),
			},
			{
				Config: testAccActiveDirectoryConfigFull(testingActiveDirectoryConfigSettings),
				Check: resource.ComposeTestCheckFunc(
					testAccCompareActiveDirectorySettings(testingActiveDirectoryConfigSettings),
					testAccCheckActiveDirectorySettings(*testingActiveDirectoryConfigSettings.Settings),
					testAccCheckActiveDirectoryStatus(*testingActiveDirectoryConfigSettings.JoinSettings),
				),
			},
		},
	})
}

// test update join settings (force new)
// test update usage settings (reconfigure)
// test empty settings
// test partially empty settings does nothing

// Active Directory configurations

// TODO create all relevant combinations for test

var defaultActiveDirectoryConfig = ActiveDirectoryRequest{
	Settings:     &defaultActiveDirectorySettingsConfig,
	JoinSettings: &testingActiveDirectoryJoinSettingsConfigFull,
}

var testingActiveDirectoryConfigSettings = ActiveDirectoryRequest{
	Settings:     &testingActiveDirectorySettingsConfigFull,
	JoinSettings: &testingActiveDirectoryJoinSettingsConfigFull,
}

// Active Directory Settings configurations

var defaultActiveDirectorySettingsConfig = ActiveDirectorySettingsBody{
	Signing: "WANT_SIGNING",
	Sealing: "WANT_SEALING",
	Crypto:  "WANT_AES",
}

var testingActiveDirectorySettingsConfigFull = ActiveDirectorySettingsBody{
	Signing: "REQUIRE_SIGNING",
	Sealing: "WANT_SEALING",
	Crypto:  "REQUIRE_AES",
}

var testingActiveDirectorySettingsConfigPartial = ActiveDirectorySettingsBody{
	Signing: "REQUIRE_SIGNING",
	Crypto:  "REQUIRE_AES",
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

		if adSettings.Sealing != adSettingsRequest.Sealing || adSettings.Signing != adSettingsRequest.Signing || adSettings.Crypto != adSettingsRequest.Crypto {
			return fmt.Errorf("Active Directory settings mismatch: Expected %v, got %v", adSettingsRequest, adSettings)
		}
		return nil
	}
}

func testAccCheckActiveDirectoryStatus(adJoinSettingsRequest ActiveDirectoryJoinRequest) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		c := testAccProvider.Meta().(*Client)
		ctx := context.Background()
		adStatus, err := DoRequest[ActiveDirectoryStatusBody, ActiveDirectoryStatusBody](ctx, c, GET, AdStatusEndpoint, nil)
		if err != nil {
			return err
		}

		if adStatus.Domain != adJoinSettingsRequest.Domain || adStatus.DomainNetBios != adJoinSettingsRequest.DomainNetBios || adStatus.Ou != adJoinSettingsRequest.Ou ||
			adStatus.UseAdPosixAttributes != adJoinSettingsRequest.UseAdPosixAttributes || adStatus.BaseDn != adJoinSettingsRequest.BaseDn {

			return fmt.Errorf("Active Directory status mismatch: Expected %v, got %v", adJoinSettingsRequest, adStatus)
		}
		return nil
	}
}
