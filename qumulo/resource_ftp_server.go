package qumulo

import (
	"context"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strconv"
	"time"
)

const FtpServerEndpoint = "/v0/ftp/settings"

type FtpServerBody struct {
	Enabled                     bool                    `json:"enabled"`
	CheckRemoteHost             bool                    `json:"check_remote_host"`
	LogOperations               bool                    `json:"log_operations"`
	ChrootUsers                 bool                    `json:"chroot_users"`
	AllowUnencryptedConnections bool                    `json:"allow_unencrypted_connections"`
	ExpandWildcards             bool                    `json:"expand_wildcards"`
	AnonymousUser               *map[string]interface{} `json:"anonymous_user"`
	Greeting                    string                  `json:"greeting"`
}

func resourceFtpServer() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFtpServerCreate,
		ReadContext:   resourceFtpServerRead,
		UpdateContext: resourceFtpServerUpdate,
		DeleteContext: resourceFtpServerDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"enabled": &schema.Schema{
				Type:     schema.TypeBool,
				Required: true,
			},
			"check_remote_host": &schema.Schema{
				Type:     schema.TypeBool,
				Required: true,
			},
			"log_operations": &schema.Schema{
				Type:     schema.TypeBool,
				Required: true,
			},
			"chroot_users": &schema.Schema{
				Type:     schema.TypeBool,
				Required: true,
			},
			"allow_unencrypted_connections": &schema.Schema{
				Type:     schema.TypeBool,
				Required: true,
			},
			"expand_wildcards": &schema.Schema{
				Type:     schema.TypeBool,
				Required: true,
			},
			"anonymous_user": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Default: nil,
			},
			"greeting": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceFtpServerCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	err := modifyFtpServerSettings(ctx, d, m, PUT)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return resourceFtpServerRead(ctx, d, m)
}

func resourceFtpServerRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	var errs ErrorCollection

	ftpServer, err := DoRequest[FtpServerBody, FtpServerBody](ctx, c, GET, FtpServerEndpoint, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	errs.addMaybeError(d.Set("enabled", ftpServer.Enabled))
	errs.addMaybeError(d.Set("check_remote_host", ftpServer.CheckRemoteHost))
	errs.addMaybeError(d.Set("log_operations", ftpServer.LogOperations))
	errs.addMaybeError(d.Set("chroot_users", ftpServer.ChrootUsers))
	errs.addMaybeError(d.Set("allow_unencrypted_connections", ftpServer.AllowUnencryptedConnections))
	errs.addMaybeError(d.Set("expand_wildcards", ftpServer.ExpandWildcards))
	errs.addMaybeError(d.Set("anonymous_user", ftpServer.AnonymousUser))
	errs.addMaybeError(d.Set("greeting", ftpServer.Greeting))

	return errs.diags
}

func resourceFtpServerUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	err := modifyFtpServerSettings(ctx, d, m, PATCH)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceFtpServerRead(ctx, d, m)
}

func resourceFtpServerDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "Deleting Ftp server resource")
	return nil
}

func modifyFtpServerSettings(ctx context.Context, d *schema.ResourceData, m interface{}, method Method) error {
	c := m.(*Client)

	ftpServer := FtpServerBody{
		Enabled:                     d.Get("enabled").(bool),
		CheckRemoteHost:             d.Get("check_remote_host").(bool),
		LogOperations:               d.Get("log_operations").(bool),
		ChrootUsers:                 d.Get("chroot_users").(bool),
		AllowUnencryptedConnections: d.Get("allow_unencrypted_connections").(bool),
		ExpandWildcards:             d.Get("expand_wildcards").(bool),
		Greeting:                    d.Get("greeting").(string),
	}
	if v, ok := d.Get("anonymous_user").(map[string]interface{}); ok && len(v) > 0 {
		ftpServer.AnonymousUser = &v
	}
	tflog.Debug(ctx, "Modifying FTP server settings")
	_, err := DoRequest[FtpServerBody, FtpServerBody](ctx, c, method, FtpServerEndpoint, &ftpServer)
	return err
}
