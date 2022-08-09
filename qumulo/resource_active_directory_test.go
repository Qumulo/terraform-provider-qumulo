package qumulo

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccJoinActiveDirectoryFull(t *testing.T) {

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

func TestAccJoinActiveDirectoryPartial(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccActiveDirectoryConfigFull(testingActiveDirectoryConfigPartialJoin),
				Check: resource.ComposeTestCheckFunc(
					testAccCompareActiveDirectorySettings(defaultActiveDirectoryConfig),
					testAccCheckActiveDirectorySettings(*defaultActiveDirectoryConfig.Settings),
					testAccCheckActiveDirectoryStatus(*defaultActiveDirectoryConfig.JoinSettings),
				),
				// This is treated as an update, which has a force-new update on the Ou field
				ExpectNonEmptyPlan: true,
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

func TestAccChangeActiveDirectoryStatusForceNew(t *testing.T) {

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
				// This config changes the DomainNetBios, which forces a new (leave and re-join) of the AD.
				Config: testAccActiveDirectoryConfigFull(testingActiveDirectoryConfigSettingsForceNew),
				Check: resource.ComposeTestCheckFunc(
					// XXX amanning32: Qumulo's AD server sets the DomainNetBios to "AD" on join, so verify against those settings
					// even though we changed the DomainNetBios to force the re-creation of the resource.
					testAccCompareActiveDirectorySettings(testingActiveDirectoryConfigSettings),
					testAccCheckActiveDirectorySettings(*testingActiveDirectoryConfigSettings.Settings),
					testAccCheckActiveDirectoryStatus(*testingActiveDirectoryConfigSettings.JoinSettings),
				),
				// Force new creates a non-empty plan, which is expected
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccChangeActiveDirectoryStatusReconfigure(t *testing.T) {

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
				Config: testAccActiveDirectoryConfigFull(testingActiveDirectoryConfigSettingsReconfigure),
				Check: resource.ComposeTestCheckFunc(
					// XXX amanning32: Qumulo's AD server sets the BaseDn on join, so verify against those settings
					// even though we changed the BaseDn to force the reconfiguration of the resource.
					testAccCompareActiveDirectorySettings(testingActiveDirectoryConfigSettings),
					testAccCheckActiveDirectorySettings(*testingActiveDirectoryConfigSettings.Settings),
					testAccCheckActiveDirectoryStatus(*testingActiveDirectoryConfigSettings.JoinSettings),
				),
				// Update creates a non-empty plan, which is expected
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccChangeActiveDirectorySettingsEmpty(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccActiveDirectoryConfigNoSettings(defaultActiveDirectoryConfig),
				Check: resource.ComposeTestCheckFunc(
					testAccCompareActiveDirectorySettings(defaultActiveDirectoryConfig),
					testAccCheckActiveDirectorySettings(*defaultActiveDirectoryConfig.Settings),
					testAccCheckActiveDirectoryStatus(*defaultActiveDirectoryConfig.JoinSettings),
				),
			},
		},
	})
}

func TestAccChangeActiveDirectorySettingsPartial(t *testing.T) {

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
				Config: testAccActiveDirectoryConfigPartialSettings(testingActiveDirectoryConfigSettings),
				Check: resource.ComposeTestCheckFunc(
					testAccCompareActiveDirectorySettings(testingActiveDirectoryConfigSettings),
					testAccCheckActiveDirectorySettings(*testingActiveDirectoryConfigSettings.Settings),
					testAccCheckActiveDirectoryStatus(*testingActiveDirectoryConfigSettings.JoinSettings),
				),
			},
		},
	})
}

func TestAccChangeActiveDirectorySettingsInvalid_ExpectError(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccActiveDirectoryConfigFull(testingActiveDirectoryConfigSettingsInvalid),
				ExpectError: regexp.MustCompile("expected (signing|sealing|crypto) to be one of .+, got WHAT_IS_(SIGNING|SEALING|AES)"),
			},
		},
	})
}

// Active Directory configurations

var defaultActiveDirectoryConfig = ActiveDirectoryRequest{
	Settings:     &defaultActiveDirectorySettingsConfig,
	JoinSettings: &testingActiveDirectoryJoinSettingsConfigFull,
}

var testingActiveDirectoryConfigPartialJoin = ActiveDirectoryRequest{
	Settings:     &defaultActiveDirectorySettingsConfig,
	JoinSettings: &testingActiveDirectoryJoinSettingsConfigPartial,
}

var testingActiveDirectoryConfigSettings = ActiveDirectoryRequest{
	Settings:     &testingActiveDirectorySettingsConfigFull,
	JoinSettings: &testingActiveDirectoryJoinSettingsConfigFull,
}

var testingActiveDirectoryConfigSettingsForceNew = ActiveDirectoryRequest{
	Settings:     &testingActiveDirectorySettingsConfigFull,
	JoinSettings: &testingActiveDirectoryJoinSettingsConfigFullForceNew,
}

var testingActiveDirectoryConfigSettingsReconfigure = ActiveDirectoryRequest{
	Settings:     &testingActiveDirectorySettingsConfigFull,
	JoinSettings: &testingActiveDirectoryJoinSettingsConfigFullReconfigure,
}

var testingActiveDirectoryConfigSettingsInvalid = ActiveDirectoryRequest{
	Settings:     &testingActiveDirectorySettingsConfigInvalid,
	JoinSettings: &testingActiveDirectoryJoinSettingsConfigFull,
}

// Active Directory Settings configurations

var defaultActiveDirectorySettingsConfig = ActiveDirectorySettingsBody{
	Signing: WantSigning.String(),
	Sealing: WantSealing.String(),
	Crypto:  WantCrypto.String(),
}

var testingActiveDirectorySettingsConfigFull = ActiveDirectorySettingsBody{
	Signing: RequireSigning.String(),
	Sealing: WantSealing.String(),
	Crypto:  RequireCrypto.String(),
}

var testingActiveDirectorySettingsConfigInvalid = ActiveDirectorySettingsBody{
	Signing: "WHAT_IS_SIGNING",
	Sealing: "WHAT_IS_SEALING",
	Crypto:  "WHAT_IS_AES",
}

// Active Directory Join Settings configurations

var testingActiveDirectoryJoinSettingsConfigFull = ActiveDirectoryJoinRequest{
	Domain:               "ad.eng.qumulo.com",
	DomainNetBios:        "AD",
	User:                 "Administrator",
	Password:             "a",
	UseAdPosixAttributes: false,
	BaseDn:               "CN=Users,DC=ad,DC=eng,DC=qumulo,DC=com",
}

// changes DomainNetBios
var testingActiveDirectoryJoinSettingsConfigFullForceNew = ActiveDirectoryJoinRequest{
	Domain:               "ad.eng.qumulo.com",
	User:                 "Administrator",
	Password:             "a",
	UseAdPosixAttributes: false,
	BaseDn:               "CN=Users,DC=ad,DC=eng,DC=qumulo,DC=com",
}

// changes BaseDn
var testingActiveDirectoryJoinSettingsConfigFullReconfigure = ActiveDirectoryJoinRequest{
	Domain:               "ad.eng.qumulo.com",
	DomainNetBios:        "AD",
	User:                 "Administrator",
	Password:             "a",
	UseAdPosixAttributes: false,
}

var testingActiveDirectoryJoinSettingsConfigPartial = ActiveDirectoryJoinRequest{
	Domain:   "ad.eng.qumulo.com",
	User:     "Administrator",
	Password: "a",
}

// Formatting functions

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

func testAccActiveDirectoryConfigPartialSettings(req ActiveDirectoryRequest) string {
	return fmt.Sprintf(`
	resource "qumulo_ad_settings" "ad_settings" {
		signing = %q
		crypto = %q
		domain = %q
		domain_netbios = %q
		ad_username = %q
		ad_password = %q
		ou = %q
		use_ad_posix_attributes = %v
		base_dn = %q
	}
	`, req.Settings.Signing, req.Settings.Crypto, req.JoinSettings.Domain, req.JoinSettings.DomainNetBios, req.JoinSettings.User,
		req.JoinSettings.Password, req.JoinSettings.Ou, req.JoinSettings.UseAdPosixAttributes, req.JoinSettings.BaseDn)
}

func testAccActiveDirectoryConfigNoSettings(req ActiveDirectoryRequest) string {
	return fmt.Sprintf(`
	resource "qumulo_ad_settings" "ad_settings" {
		domain = %q
		domain_netbios = %q
		ad_username = %q
		ad_password = %q
		ou = %q
		use_ad_posix_attributes = %v
		base_dn = %q
	}
	`, req.JoinSettings.Domain, req.JoinSettings.DomainNetBios, req.JoinSettings.User,
		req.JoinSettings.Password, req.JoinSettings.Ou, req.JoinSettings.UseAdPosixAttributes, req.JoinSettings.BaseDn)
}

// Test helper functions

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
			//lint:ignore ST1005 proper nouns should be capitalized
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
			//lint:ignore ST1005 proper nouns should be capitalized
			return fmt.Errorf("Active Directory status mismatch: Expected %v, got %v", adJoinSettingsRequest, adStatus)
		}
		return nil
	}
}
