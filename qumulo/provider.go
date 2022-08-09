package qumulo

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider -
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"host": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("QUMULO_HOST", nil),
			},
			"port": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("QUMULO_PORT", nil),
			},
			"username": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("QUMULO_USERNAME", nil),
			},
			"password": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("QUMULO_PASSWORD", nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"qumulo_cluster_name":            resourceClusterSettings(),
			"qumulo_ad_settings":             resourceActiveDirectory(),
			"qumulo_ldap_server":             resourceLdapServer(),
			"qumulo_ssl_cert":                resourceSsl(),
			"qumulo_ssl_ca":                  resourceSslCa(),
			"qumulo_monitoring":              resourceMonitoring(),
			"qumulo_nfs_export":              resourceNfsExport(),
			"qumulo_nfs_settings":            resourceNfsSettings(),
			"qumulo_smb_server":              resourceSmbServer(),
			"qumulo_smb_share":               resourceSmbShare(),
			"qumulo_role":                    resourceRole(),
			"qumulo_time_configuration":      resourceTimeConfiguration(),
			"qumulo_directory_quota":         resourceDirectoryQuota(),
			"qumulo_local_user":              resourceUser(),
			"qumulo_local_group":             resourceGroup(),
			"qumulo_web_ui":                  resourceWebUi(),
			"qumulo_file_system_settings":    resourceFileSystemSettings(),
			"qumulo_interface_configuration": resourceInterfaceConfiguration(),
			"qumulo_network_configuration":   resourceNetworkConfiguration(),
			"qumulo_ftp_server":              resourceFtpServer(),
			"qumulo_local_group_member":      resourceGroupMember(),
			"qumulo_syslog":                  resourceSyslog(),
			"qumulo_cloudwatch":              resourceCloudWatch(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {

	host := d.Get("host").(string)
	port := d.Get("port").(string)
	username := d.Get("username").(string)
	password := d.Get("password").(string)

	c, err := NewClient(ctx, &host, &port, &username, &password)
	if err != nil {
		return nil, diag.FromErr(err)
	}

	return c, nil
}
