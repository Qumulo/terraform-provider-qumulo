package qumulo

//TODO testing after imports
import (
	"context"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"strconv"
	"time"
)

const NetworksEndpointSuffix = "/networks/"

var assignedByValues = []string{"DHCP", "STATIC"}

type NetworkConfiguration struct {
	Id               int      `json:"id"`
	Name             string   `json:"name"`
	AssignedBy       string   `json:"assigned_by"`
	FloatingIpRanges []string `json:"floating_ip_ranges"`
	DnsServers       []string `json:"dns_servers"`
	DnsSearchDomains []string `json:"dns_search_domains"`
	IpRanges         []string `json:"ip_ranges"`
	Netmask          string   `json:"netmask"`
	Mtu              int      `json:"mtu"`
	VlanId           int      `json:"vlan_id"`
}

func resourceNetworkConfiguration() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNetworkConfigurationCreate,
		ReadContext:   resourceNetworkConfigurationRead,
		UpdateContext: resourceNetworkConfigurationUpdate,
		DeleteContext: resourceNetworkConfigurationDelete,

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
				Required: true,
			},
			"assigned_by": &schema.Schema{
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice(assignedByValues, false)),
			},
			"floating_ip_ranges": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"dns_servers": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"dns_search_domains": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"ip_ranges": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"netmask": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"mtu": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"vlan_id": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"interface_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"network_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceNetworkConfigurationCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	interfaceId := d.Get("interface_id").(string)
	addNetworkConfigUri := InterfaceConfigurationEndpoint + interfaceId + NetworksEndpointSuffix
	res, err := addOrPatchNetworkConfiguration(ctx, d, m, POST, addNetworkConfigUri)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(strconv.Itoa(res.Id))
	return resourceNetworkConfigurationRead(ctx, d, m)
}

func resourceNetworkConfigurationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	var errs ErrorCollection

	interfaceId := d.Get("interface_id").(string)
	networkId := d.Id()
	readNetworkConfigUri := InterfaceConfigurationEndpoint + interfaceId + NetworksEndpointSuffix + networkId
	networkConfig, err := DoRequest[NetworkConfiguration, NetworkConfiguration](ctx, c, GET, readNetworkConfigUri, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	errs.addMaybeError(d.Set("id", strconv.Itoa(networkConfig.Id)))
	errs.addMaybeError(d.Set("name", networkConfig.Name))
	errs.addMaybeError(d.Set("assigned_by", networkConfig.AssignedBy))
	errs.addMaybeError(d.Set("floating_ip_ranges", networkConfig.FloatingIpRanges))
	errs.addMaybeError(d.Set("dns_servers", networkConfig.DnsServers))
	errs.addMaybeError(d.Set("dns_search_domains", networkConfig.DnsSearchDomains))
	errs.addMaybeError(d.Set("ip_ranges", networkConfig.IpRanges))
	errs.addMaybeError(d.Set("netmask", networkConfig.Netmask))
	errs.addMaybeError(d.Set("mtu", networkConfig.Mtu))
	errs.addMaybeError(d.Set("vlan_id", networkConfig.VlanId))

	return errs.diags
}

func resourceNetworkConfigurationUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	interfaceId := d.Get("interface_id").(string)
	networkId := d.Get("network_id").(string)
	updateNetworkConfigUri := InterfaceConfigurationEndpoint + interfaceId + NetworksEndpointSuffix + networkId
	_, err := addOrPatchNetworkConfiguration(ctx, d, m, PATCH, updateNetworkConfigUri)
	if err != nil {
		return diag.FromErr(err)
	}
	return resourceNetworkConfigurationRead(ctx, d, m)
}

func resourceNetworkConfigurationDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "Deleting network configuration", map[string]interface{}{
		"Id": d.Id(),
	})

	c := m.(*Client)

	var diags diag.Diagnostics
	interfaceId := d.Get("interface_id").(string)
	networkId := d.Id()
	deleteNetworkConfigUri := InterfaceConfigurationEndpoint + interfaceId + NetworksEndpointSuffix + networkId
	_, err := DoRequest[NetworkConfiguration, NetworkConfiguration](ctx, c, DELETE, deleteNetworkConfigUri, nil)
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func addOrPatchNetworkConfiguration(ctx context.Context, d *schema.ResourceData, m interface{}, method Method, uri string) (*NetworkConfiguration, error) {
	c := m.(*Client)

	networkConfig := NetworkConfiguration{
		Name:             d.Get("name").(string),
		AssignedBy:       d.Get("assigned_by").(string),
		FloatingIpRanges: InterfaceSliceToStringSlice(d.Get("floating_ip_ranges").([]interface{})),
		DnsServers:       InterfaceSliceToStringSlice(d.Get("dns_servers").([]interface{})),
		DnsSearchDomains: InterfaceSliceToStringSlice(d.Get("dns_search_domains").([]interface{})),
		IpRanges:         InterfaceSliceToStringSlice(d.Get("ip_ranges").([]interface{})),
		Netmask:          d.Get("netmask").(string),
		Mtu:              d.Get("mtu").(int),
		VlanId:           d.Get("vlan_id").(int),
	}

	tflog.Debug(ctx, "Adding/Patching network configuration")
	res, err := DoRequest[NetworkConfiguration, NetworkConfiguration](ctx, c, method, uri, &networkConfig)
	return res, err
}
