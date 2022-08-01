package qumulo

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccSetSslCa(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{ // Reset state to default
				Config: defaultSslCaConfig,
			},
			{
				Config: testAccSslCaConfig(testingSslCaCert),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSslCaCert(),
					resource.TestCheckResourceAttrSet("qumulo_ssl_ca.cert", "ca_certificate"),
				),
			},
		},
	})
}

func TestAccSetSslCa_ExpectError(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSslCaConfig(invalidCert),
				// Invalid certificate should raise an error
				ExpectError: regexp.MustCompile("ssl_certificate_invalid_error"),
			},
		},
	})
}

var defaultSslCaConfig = " "

var invalidCert = SslCaBody{
	CaCertificate: "Not a valid certificate",
}

var testingSslCaCert = SslCaBody{
	CaCertificate: `-----BEGIN CERTIFICATE-----
MIICIDCCAYmgAwIBAgIUZcdqCxZB1O4RD548ygFhGBXxQdQwDQYJKoZIhvcNAQEL
BQAwIjEPMA0GA1UEAwwGVGVzdENBMQ8wDQYDVQQKDAZRdW11bG8wHhcNMjIwNzIy
MTcwOTI4WhcNMzIwNzE5MTcwOTI4WjAiMQ8wDQYDVQQDDAZUZXN0Q0ExDzANBgNV
BAoMBlF1bXVsbzCBnzANBgkqhkiG9w0BAQEFAAOBjQAwgYkCgYEAv9Xupp43GfpI
0bVkB1BIa0ZBt5hpjxgee5PKwn3pbcg/M0M4qGhtX9/DR4utMqMib+X517hyo18E
Vd+gZa0plafaPfwzz8YkO2EovYEFIaBxgqYkTQ0YZVt40cWEMMCWuyPndX0bvOrW
1f5zvOcc0+dDXoiqbhUDKiXBfzK745UCAwEAAaNTMFEwHQYDVR0OBBYEFKYiYrFK
cZcR+gDTAqxV6u81B9htMB8GA1UdIwQYMBaAFKYiYrFKcZcR+gDTAqxV6u81B9ht
MA8GA1UdEwEB/wQFMAMBAf8wDQYJKoZIhvcNAQELBQADgYEAjPXNGT38WwyWu4Xe
Wngxmk0OIKZthsbZVDxSti3mse7KWadb6EkaRM/ZIO9CFPyB67zh3KAwhKiMbPVE
JH62qN5t5xoqdDzzuOUHw1SSF78lfMAWk84TplzXegdysXjYFVhxvqYV9DIEhsTw
HjX0jrbwN2tDfjTKNQwi7P7RPDY=
-----END CERTIFICATE-----`,
}

func testAccSslCaConfig(req SslCaBody) string {
	return fmt.Sprintf(`
resource "qumulo_ssl_ca" "cert" {
	ca_certificate = <<CERTDELIM
%v
CERTDELIM
}
  `, req.CaCertificate)
}

func testAccCheckSslCaCert() resource.TestCheckFunc {
	// Make sure there's a valid certificate through the API
	return func(s *terraform.State) error {
		c := testAccProvider.Meta().(*Client)
		ctx := context.Background()
		_, err := DoRequest[SslCaBody, SslCaBody](ctx, c, GET, SslCaEndpoint, nil)
		if err != nil {
			return err
		}
		return nil
	}
}
