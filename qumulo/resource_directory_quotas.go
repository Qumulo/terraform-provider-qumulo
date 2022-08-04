package qumulo

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const DirectoryQuotaEndpoint = "/v1/files/quotas/"

// Create body, Read response, Update body
type DirectoryQuotaBody struct {
	Id    string `json:"id"`
	Limit string `json:"limit"`
}

// Read request, Delete body
type DirectoryQuotaEmptyBody struct{}

func resourceDirectoryQuota() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDirectoryQuotaCreate,
		ReadContext:   resourceDirectoryQuotaRead,
		UpdateContext: resourceDirectoryQuotaUpdate,
		DeleteContext: resourceDirectoryQuotaDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"directory_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"limit": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceDirectoryQuotaCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	err := createOrUpdateDirectoryQuota(ctx, d, m, POST, DirectoryQuotaEndpoint)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return resourceDirectoryQuotaRead(ctx, d, m)
}

func resourceDirectoryQuotaRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	var errs ErrorCollection

	quotaUrl := DirectoryQuotaEndpoint + d.Get("directory_id").(string)

	directoryQuota, err := DoRequest[DirectoryQuotaEmptyBody, DirectoryQuotaBody](ctx, c, GET, quotaUrl, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	errs.addMaybeError(d.Set("directory_id", directoryQuota.Id))
	errs.addMaybeError(d.Set("limit", directoryQuota.Limit))

	return errs.diags
}

func resourceDirectoryQuotaUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	directoryId := d.Get("directory_id").(string)
	quotaUrl := DirectoryQuotaEndpoint + directoryId

	err := createOrUpdateDirectoryQuota(ctx, d, m, PUT, quotaUrl)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceDirectoryQuotaRead(ctx, d, m)
}

func resourceDirectoryQuotaDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	var diags diag.Diagnostics

	quotaUrl := DirectoryQuotaEndpoint + d.Get("directory_id").(string)

	_, err := DoRequest[DirectoryQuotaEmptyBody, DirectoryQuotaEmptyBody](ctx, c, DELETE, quotaUrl, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func createOrUpdateDirectoryQuota(ctx context.Context, d *schema.ResourceData, m interface{}, method Method, url string) error {
	c := m.(*Client)

	directoryId := d.Get("directory_id").(string)

	directoryQuotaRequest := DirectoryQuotaBody{
		Id:    directoryId,
		Limit: d.Get("limit").(string),
	}

	tflog.Debug(ctx, fmt.Sprintf("Updating directory quota with id %q", directoryId))
	_, err := DoRequest[DirectoryQuotaBody, DirectoryQuotaBody](ctx, c, method, url, &directoryQuotaRequest)

	return err
}
