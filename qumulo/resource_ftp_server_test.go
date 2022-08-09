package qumulo

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"reflect"
	"testing"
)

func TestAccFtpServer(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccFtpServer1(defaultFtpServer),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFtpServer(defaultFtpServer)),
			},
			{
				Config: testAccFtpServer1(testFtpServer),
				Check: resource.ComposeTestCheckFunc(
					testAccCompareFtpServers(testFtpServer),
					testAccCheckFtpServer(testFtpServer)),
			},
			{
				Config: testAccFtpServer2(testFtpServer2),
				Check: resource.ComposeTestCheckFunc(
					testAccCompareFtpServers(testFtpServer2),
					testAccCheckFtpServer(testFtpServer2)),
			},
		},
	})
}

var defaultFtpServer = FtpServerBody{
	Enabled:                     true,
	CheckRemoteHost:             true,
	LogOperations:               true,
	ChrootUsers:                 true,
	AllowUnencryptedConnections: true,
	ExpandWildcards:             false,
	AnonymousUser:               nil,
	Greeting:                    "Hello!",
}

var testFtpServer = FtpServerBody{
	Enabled:                     true,
	CheckRemoteHost:             false,
	LogOperations:               true,
	ChrootUsers:                 false,
	AllowUnencryptedConnections: true,
	ExpandWildcards:             false,
	AnonymousUser:               nil,
	Greeting:                    "Hello!!!",
}

var testFtpServer2 = FtpServerBody{
	Enabled:                     true,
	CheckRemoteHost:             false,
	LogOperations:               true,
	ChrootUsers:                 false,
	AllowUnencryptedConnections: true,
	ExpandWildcards:             false,
	AnonymousUser: &map[string]interface{}{
		"id_type":  "LOCAL_USER",
		"id_value": "admin",
	},
	Greeting: "Hello!!!!",
}

func testAccFtpServer1(fs FtpServerBody) string {
	return fmt.Sprintf(`
resource "qumulo_ftp_server" "some_ftp_server" {
  enabled = %v
  check_remote_host = %v
  log_operations = %v
  chroot_users = %v
  allow_unencrypted_connections = %v
  expand_wildcards = %v
  greeting = %q
}`, fs.Enabled, fs.CheckRemoteHost, fs.LogOperations, fs.ChrootUsers, fs.AllowUnencryptedConnections,
		fs.ExpandWildcards, fs.Greeting)
}

func testAccFtpServer2(fs FtpServerBody) string {
	return fmt.Sprintf(`
resource "qumulo_ftp_server" "some_ftp_server" {
  enabled = %v
  check_remote_host = %v
  log_operations = %v
  chroot_users = %v
  allow_unencrypted_connections = %v
  expand_wildcards = %v
  anonymous_user = {
	id_type = %q
	id_value = %q
	}
  greeting = %q
}`, fs.Enabled, fs.CheckRemoteHost, fs.LogOperations, fs.ChrootUsers, fs.AllowUnencryptedConnections,
		fs.ExpandWildcards, (*fs.AnonymousUser)["id_type"], (*fs.AnonymousUser)["id_value"], fs.Greeting)
}

func testAccCompareFtpServers(fs FtpServerBody) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("qumulo_ftp_server.some_ftp_server", "enabled",
			fmt.Sprintf("%v", fs.Enabled)),
		resource.TestCheckResourceAttr("qumulo_ftp_server.some_ftp_server", "check_remote_host",
			fmt.Sprintf("%v", fs.CheckRemoteHost)),
		resource.TestCheckResourceAttr("qumulo_ftp_server.some_ftp_server", "log_operations",
			fmt.Sprintf("%v", fs.LogOperations)),
		resource.TestCheckResourceAttr("qumulo_ftp_server.some_ftp_server", "chroot_users",
			fmt.Sprintf("%v", fs.ChrootUsers)),
		resource.TestCheckResourceAttr("qumulo_ftp_server.some_ftp_server", "allow_unencrypted_connections",
			fmt.Sprintf("%v", fs.AllowUnencryptedConnections)),
		resource.TestCheckResourceAttr("qumulo_ftp_server.some_ftp_server", "expand_wildcards",
			fmt.Sprintf("%v", fs.ExpandWildcards)),
		resource.TestCheckResourceAttr("qumulo_ftp_server.some_ftp_server", "greeting",
			fmt.Sprintf("%v", fs.Greeting)),
	)
}

func testAccCheckFtpServer(fs FtpServerBody) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		c := testAccProvider.Meta().(*Client)
		ctx := context.Background()
		ftpServer, err := DoRequest[FtpServerBody, FtpServerBody](ctx, c, GET, FtpServerEndpoint, nil)
		if err != nil {
			return err
		}

		if !reflect.DeepEqual(*ftpServer, fs) {
			return fmt.Errorf("ftp server mismatch: Expected %v, got %v", fs, *ftpServer)
		}
		return nil
	}
}
