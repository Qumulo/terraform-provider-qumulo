package qumulo

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"strconv"
	"time"
)

const InterfaceConfigurationEndpoint = "/v2/network/interfaces/"

type InterfaceConfigurationResponse struct {
	Id                 int    `json:"id"`
	Name               string `json:"name"`
	DefaultGateway     string `json:"default_gateway"`
	DefaultGatewayIpv6 string `json:"default_gateway_ipv6"`
	BondingMode        string `json:"bonding_mode"`
	Mtu                int    `json:"mtu"`
}

type InterfaceConfigurationRequest struct {
	Id                 int    `json:"id"`
	Name               string `json:"name"`
	DefaultGateway     string `json:"default_gateway"`
	DefaultGatewayIpv6 string `json:"default_gateway_ipv6"`
	BondingMode        string `json:"bonding_mode"`
	Mtu                int    `json:"mtu"`
	InterfaceId        string `json:"interface_id"`
}
type InterfaceConfigBondingMode int

var InterfaceConfigBondingModes = []string{"ACTIVE_BACKUP", "IEEE_8023AD", "UNSPECIFIED"}

const (
	ActiveBackup InterfaceConfigBondingMode = iota + 1
	Ieee8023Ad
	Unspecified
)

func (e InterfaceConfigBondingMode) String() string {
	return InterfaceConfigBondingModes[e-1]
}

func resourceInterfaceConfiguration() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceInterfaceConfigurationCreate,
		ReadContext:   resourceInterfaceConfigurationRead,
		UpdateContext: resourceInterfaceConfigurationUpdate,
		DeleteContext: resourceInterfaceConfigurationDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"default_gateway": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"default_gateway_ipv6": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"bonding_mode": &schema.Schema{
				Type:             schema.TypeString,
				Optional:         true,
				Default:          Unspecified,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice(InterfaceConfigBondingModes, false)),
			},
			"mtu": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"interface_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceInterfaceConfigurationCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	err := setOrPatchInterfaceConfiguration(ctx, d, m, PUT)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(d.Get("interface_id").(string))

	return resourceInterfaceConfigurationRead(ctx, d, m)
}

func resourceInterfaceConfigurationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	var errs ErrorCollection
	interfaceId := d.Id()

	interfaceConfigUri := InterfaceConfigurationEndpoint + interfaceId
	interfaceConfig, err := DoRequest[InterfaceConfigurationRequest, InterfaceConfigurationResponse](ctx, c, GET, interfaceConfigUri, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	errs.addMaybeError(d.Set("id", strconv.Itoa(interfaceConfig.Id)))
	errs.addMaybeError(d.Set("name", interfaceConfig.Name))
	errs.addMaybeError(d.Set("default_gateway", interfaceConfig.DefaultGateway))
	errs.addMaybeError(d.Set("default_gateway_ipv6", interfaceConfig.DefaultGatewayIpv6))
	errs.addMaybeError(d.Set("bonding_mode", interfaceConfig.BondingMode))
	errs.addMaybeError(d.Set("mtu", interfaceConfig.Mtu))

	return errs.diags
}

func resourceInterfaceConfigurationUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	err := setOrPatchInterfaceConfiguration(ctx, d, m, PATCH)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceInterfaceConfigurationRead(ctx, d, m)
}

func resourceInterfaceConfigurationDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, fmt.Sprintf("Deleting interface configuration resource with id: %q", d.Id()))
	return nil
}

func setOrPatchInterfaceConfiguration(ctx context.Context, d *schema.ResourceData, m interface{}, method Method) error {
	c := m.(*Client)

	interfaceId := d.Get("interface_id").(string)
	interfaceConfigUri := InterfaceConfigurationEndpoint + interfaceId

	//ID has to be set to the interface ID passed in the URI as per API validation
	id, _ := strconv.Atoi(interfaceId)

	interfaceConfig := InterfaceConfigurationRequest{
		Id:                 id,
		Name:               d.Get("name").(string),
		DefaultGateway:     d.Get("default_gateway").(string),
		DefaultGatewayIpv6: d.Get("default_gateway_ipv6").(string),
		BondingMode:        d.Get("bonding_mode").(string),
		Mtu:                d.Get("mtu").(int),
		InterfaceId:        interfaceId,
	}

	tflog.Debug(ctx, "Setting/Patching interface configuration")
	_, err := DoRequest[InterfaceConfigurationRequest, InterfaceConfigurationResponse](ctx, c, method, interfaceConfigUri, &interfaceConfig)
	return err
}
