package qumulo

import (
	"context"
	"fmt"
	"terraform-provider-qumulo/openapi"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

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
				ForceNew: true,
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
	c := m.(*openapi.APIClient)

	quota := createOrUpdateDirectoryQuota(ctx, d, m)

	resp, r, err := c.FilesApi.V1FilesQuotasPost(context.Background()).
		V1FilesQuotasGet200ResponseQuotasInner(quota).Execute()
	if err != nil {
		tflog.Debug(ctx, fmt.Sprintf("Error when calling `FilesApi.V1FilesQuotasIdPost``: %v\n", err))
		tflog.Debug(ctx, fmt.Sprintf("Full HTTP response: %v\n", r))
		return diag.FromErr(err)
	}
	// response from `V1FilesQuotasIdPost`: V1FilesQuotasGet200ResponseQuotasInner
	tflog.Debug(ctx, fmt.Sprintf("Response from `FilesApi.V1FilesQuotasIdPost`: %v\n", resp))

	d.SetId(d.Get("directory_id").(string))

	return resourceDirectoryQuotaRead(ctx, d, m)
}

func resourceDirectoryQuotaRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*openapi.APIClient)

	var errs ErrorCollection

	id := d.Id() // string | Directory ID (uint64)

	resp, r, err := c.FilesApi.V1FilesQuotasIdGet(context.Background(), id).Execute()
	if err != nil {
		tflog.Debug(ctx, fmt.Sprintf("Error when calling `FilesApi.V1FilesQuotasIdGet``: %v\n", err))
		tflog.Debug(ctx, fmt.Sprintf("Full HTTP response: %v\n", r))
		return diag.FromErr(err)
	}
	// response from `V1FilesQuotasIdGet`: V1FilesQuotasGet200ResponseQuotasInner
	tflog.Debug(ctx, fmt.Sprintf("Response from `FilesApi.V1FilesQuotasIdGet`: %v\n", resp))

	errs.addMaybeError(d.Set("directory_id", resp.GetId()))
	errs.addMaybeError(d.Set("limit", resp.GetLimit()))

	return errs.diags
}

func resourceDirectoryQuotaUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*openapi.APIClient)

	quota := createOrUpdateDirectoryQuota(ctx, d, m)

	resp, r, err := c.FilesApi.V1FilesQuotasIdPut(context.Background(), d.Get("directory_id").(string)).
		V1FilesQuotasGet200ResponseQuotasInner(quota).Execute()
	if err != nil {
		tflog.Debug(ctx, fmt.Sprintf("Error when calling `FilesApi.V1FilesQuotasIdPut``: %v\n", err))
		tflog.Debug(ctx, fmt.Sprintf("Full HTTP response: %v\n", r))
		return diag.FromErr(err)
	}
	// response from `V1FilesQuotasIdPut`: V1FilesQuotasGet200ResponseQuotasInner
	tflog.Debug(ctx, fmt.Sprintf("Response from `FilesApi.V1FilesQuotasIdPut`: %v\n", resp))

	return resourceDirectoryQuotaRead(ctx, d, m)
}

func resourceDirectoryQuotaDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*openapi.APIClient)

	tflog.Info(ctx, fmt.Sprintf("Deleting directory quota with id %q", d.Id()))

	id := d.Id() // string | Directory ID (uint64)

	r, err := c.FilesApi.V1FilesQuotasIdDelete(context.Background(), id).Execute()
	if err != nil {
		tflog.Debug(ctx, fmt.Sprintf("Error when calling `FilesApi.V1FilesQuotasIdDelete``\n"))
		tflog.Debug(ctx, fmt.Sprintf("Full HTTP response: %v\n", r))
		return diag.FromErr(err)
	}

	return nil
}

func createOrUpdateDirectoryQuota(ctx context.Context, d *schema.ResourceData, m interface{}) openapi.V1FilesQuotasGet200ResponseQuotasInner {
	quota := *openapi.NewV1FilesQuotasGet200ResponseQuotasInner()

	directoryId := d.Get("directory_id").(string)

	quota.SetId(directoryId)
	quota.SetLimit(d.Get("limit").(string))

	tflog.Debug(ctx, fmt.Sprintf("Updating directory quota with id %q", directoryId))

	return quota
}
