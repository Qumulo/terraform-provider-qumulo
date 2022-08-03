package qumulo

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const RolesEndpoint = "/v1/auth/roles/"

type Role struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Privileges  []string `json:"privileges"`
}

func resourceRole() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRoleCreate,
		ReadContext:   resourceRoleRead,
		UpdateContext: resourceRoleUpdate,
		DeleteContext: resourceRoleDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"privileges": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice(validPrivileges, false)),
				},
			},
		},

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceRoleCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	role := setRoleSettings(ctx, d, m)

	_, err := DoRequest[Role, Role](ctx, c, POST, RolesEndpoint, &role)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(d.Get("name").(string))

	return resourceRoleRead(ctx, d, m)
}

func resourceRoleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	var errs ErrorCollection

	readRoleByNameUri := RolesEndpoint + d.Id()

	role, err := DoRequest[Role, Role](ctx, c, GET, readRoleByNameUri, nil)
	if err != nil {
		return diag.FromErr(err)
	}
	errs.addMaybeError(d.Set("name", d.Id()))
	errs.addMaybeError(d.Set("description", role.Description))
	errs.addMaybeError(d.Set("privileges", role.Privileges))

	return errs.diags
}

func resourceRoleUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	role := setRoleSettings(ctx, d, m)
	updateRoleByNameUri := RolesEndpoint + d.Get("name").(string)

	_, err := DoRequest[Role, Role](ctx, c, PUT, updateRoleByNameUri, &role)

	if err != nil {
		return diag.FromErr(err)
	}
	return resourceRoleRead(ctx, d, m)
}

func resourceRoleDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	deleteRoleByNameUri := RolesEndpoint + d.Get("name").(string)

	_, err := DoRequest[Role, Role](ctx, c, DELETE, deleteRoleByNameUri, nil)

	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func setRoleSettings(ctx context.Context, d *schema.ResourceData, m interface{}) Role {
	var privileges = []string{}

	if v, ok := d.Get("privileges").([]interface{}); ok {
		privileges = make([]string, len(v))
		for i, priv := range v {
			privileges[i] = priv.(string)
		}
	}

	RoleConfig := Role{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Privileges:  privileges,
	}

	tflog.Debug(ctx, "Updating or creating Role")

	return RoleConfig
}

var validPrivileges = []string{
	"PRIVILEGE_AD_READ",
	"PRIVILEGE_AD_READ",
	"PRIVILEGE_AD_USE",
	"PRIVILEGE_AD_WRITE",
	"PRIVILEGE_ANALYTICS_READ",
	"PRIVILEGE_AUDIT_READ",
	"PRIVILEGE_AUDIT_WRITE",
	"PRIVILEGE_AUTH_CACHE_READ",
	"PRIVILEGE_AUTH_CACHE_WRITE",
	"PRIVILEGE_CHECKSUMMING_READ",
	"PRIVILEGE_CLUSTER_READ",
	"PRIVILEGE_CLUSTER_WRITE",
	"PRIVILEGE_DEBUG",
	"PRIVILEGE_DNS_READ",
	"PRIVILEGE_DNS_USE",
	"PRIVILEGE_DNS_WRITE",
	"PRIVILEGE_ENCRYPTION_READ",
	"PRIVILEGE_ENCRYPTION_WRITE",
	"PRIVILEGE_FILE_FULL_ACCESS",
	"PRIVILEGE_FS_ATTRIBUTES_READ",
	"PRIVILEGE_FS_DELETE_TREE_READ",
	"PRIVILEGE_FS_DELETE_TREE_WRITE",
	"PRIVILEGE_FS_LOCK_READ",
	"PRIVILEGE_FS_LOCK_WRITE",
	"PRIVILEGE_FS_SETTINGS_READ",
	"PRIVILEGE_FS_SETTINGS_WRITE",
	"PRIVILEGE_FTP_READ",
	"PRIVILEGE_FTP_WRITE",
	"PRIVILEGE_IDENTITY_MAPPING_READ",
	"PRIVILEGE_IDENTITY_MAPPING_WRITE",
	"PRIVILEGE_IDENTITY_READ",
	"PRIVILEGE_IDENTITY_WRITE",
	"PRIVILEGE_KERBEROS_KEYTAB_READ",
	"PRIVILEGE_KERBEROS_KEYTAB_WRITE",
	"PRIVILEGE_KERBEROS_SETTINGS_READ",
	"PRIVILEGE_KERBEROS_SETTINGS_WRITE",
	"PRIVILEGE_KV_READ",
	"PRIVILEGE_LDAP_READ",
	"PRIVILEGE_LDAP_USE",
	"PRIVILEGE_LDAP_WRITE",
	"PRIVILEGE_LOCAL_GROUP_READ",
	"PRIVILEGE_LOCAL_GROUP_WRITE",
	"PRIVILEGE_LOCAL_USER_READ",
	"PRIVILEGE_LOCAL_USER_WRITE",
	"PRIVILEGE_METRICS_CONFIG_READ",
	"PRIVILEGE_METRICS_CONFIG_WRITE",
	"PRIVILEGE_METRICS_READ",
	"PRIVILEGE_NETWORK_IP_ALLOCATION_READ",
	"PRIVILEGE_NETWORK_READ",
	"PRIVILEGE_NETWORK_WRITE",
	"PRIVILEGE_NFS_EXPORT_READ",
	"PRIVILEGE_NFS_EXPORT_WRITE",
	"PRIVILEGE_NFS_SETTINGS_READ",
	"PRIVILEGE_NFS_SETTINGS_WRITE",
	"PRIVILEGE_POWER_CYCLE",
	"PRIVILEGE_QUOTA_READ",
	"PRIVILEGE_QUOTA_WRITE",
	"PRIVILEGE_REBOOT_READ",
	"PRIVILEGE_REBOOT_USE",
	"PRIVILEGE_RECONCILER_READ",
	"PRIVILEGE_REPLICATION_OBJECT_READ",
	"PRIVILEGE_REPLICATION_OBJECT_WRITE",
	"PRIVILEGE_REPLICATION_REVERSE_RELATIONSHIP",
	"PRIVILEGE_REPLICATION_SOURCE_READ",
	"PRIVILEGE_REPLICATION_SOURCE_WRITE",
	"PRIVILEGE_REPLICATION_TARGET_READ",
	"PRIVILEGE_REPLICATION_TARGET_WRITE",
	"PRIVILEGE_ROLE_READ",
	"PRIVILEGE_ROLE_WRITE",
	"PRIVILEGE_S3_CREDENTIALS_READ",
	"PRIVILEGE_S3_CREDENTIALS_WRITE",
	"PRIVILEGE_S3_SETTINGS_READ",
	"PRIVILEGE_S3_SETTINGS_WRITE",
	"PRIVILEGE_SERVICE_PUBLIC_KEYS_READ",
	"PRIVILEGE_SERVICE_PUBLIC_KEYS_WRITE",
	"PRIVILEGE_SMB_FILE_HANDLE_READ",
	"PRIVILEGE_SMB_FILE_HANDLE_WRITE",
	"PRIVILEGE_SMB_SESSION_READ",
	"PRIVILEGE_SMB_SESSION_WRITE",
	"PRIVILEGE_SMB_SHARE_READ",
	"PRIVILEGE_SMB_SHARE_WRITE",
	"PRIVILEGE_SNAPSHOT_CALCULATE_USED_CAPACITY_READ",
	"PRIVILEGE_SNAPSHOT_DIFFERENCE_READ",
	"PRIVILEGE_SNAPSHOT_POLICY_READ",
	"PRIVILEGE_SNAPSHOT_POLICY_WRITE",
	"PRIVILEGE_SNAPSHOT_READ",
	"PRIVILEGE_SNAPSHOT_WRITE",
	"PRIVILEGE_SUPPORT_READ",
	"PRIVILEGE_SUPPORT_WRITE",
	"PRIVILEGE_TEST_ONLY",
	"PRIVILEGE_TIME_READ",
	"PRIVILEGE_TIME_WRITE",
	"PRIVILEGE_UNCONFIGURED_NODE_READ",
	"PRIVILEGE_UPGRADE_READ",
	"PRIVILEGE_UPGRADE_WRITE",
	"PRIVILEGE_WEB_UI_SETTINGS_WRITE",
}
