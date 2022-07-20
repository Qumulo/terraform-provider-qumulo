package qumulo

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const MonitoringEndpoint = "/v1/support/settings"

type MonitorRequest struct {
	Enabled             bool   `json:"enabled"`
	MQHost              string `json:"mq_host"`
	MQPort              int    `json:"mq_port"`
	MQProxyHost         string `json:"mq_proxy_host"`
	MQProxyPort         int    `json:"mq_proxy_port"`
	S3ProxyHost         string `json:"s3_proxy_host"`
	S3ProxyPort         int    `json:"s3_proxy_port"`
	S3ProxyDisableHTTPS bool   `json:"s3_proxy_disable_https"`
	VPNEnabled          bool   `json:"vpn_enabled"`
	VPNHost             string `json:"vpn_host"`
	Period              int    `json:"period"`
}

type MonitorResponse struct {
	MonitorUri string `json:"monitor_uri"`
}

func resourceMonitoring() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMonitoringCreate,
		ReadContext:   resourceMonitoringRead,
		UpdateContext: resourceMonitoringUpdate,
		DeleteContext: resourceMonitoringDelete,
		Schema: map[string]*schema.Schema{
			"enabled": &schema.Schema{
				Type:     schema.TypeBool,
				Required: true,
			},
			"mq_host": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"mq_port": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"mq_proxy_host": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"mq_proxy_port": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"s3_proxy_host": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"s3_proxy_port": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"s3_proxy_disable_https": &schema.Schema{
				Type:     schema.TypeBool,
				Required: true,
			},
			"vpn_enabled": &schema.Schema{
				Type:     schema.TypeBool,
				Required: true,
			},
			"vpn_host": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"period": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
		},
	}
}

func resourceMonitoringCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	MonitoringConfig := MonitorRequest{
		Enabled:             d.Get("enabled").(bool),
		MQHost:              d.Get("mq_host").(string),
		MQPort:              d.Get("mq_port").(int),
		MQProxyHost:         d.Get("mq_proxy_host").(string),
		MQProxyPort:         d.Get("mq_proxy_port").(int),
		S3ProxyHost:         d.Get("s3_proxy_host").(string),
		S3ProxyPort:         d.Get("s3_proxy_port").(int),
		S3ProxyDisableHTTPS: d.Get("s3_proxy_disable_https").(bool),
		VPNEnabled:          d.Get("vpn_enabled").(bool),
		VPNHost:             d.Get("vpn_host").(string),
		Period:              d.Get("period").(int),
	}

	_, err := c.UpdateMonitoring(MonitoringConfig)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}

func resourceMonitoringRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	return diags
}

func resourceMonitoringUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceMonitoringCreate(ctx, d, m)
}

func resourceMonitoringDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	return diags
}

func (c *Client) UpdateMonitoring(config MonitorRequest) (*MonitorResponse, error) {
	bearerToken := "Bearer " + c.Bearer_Token

	HostURL := c.HostURL

	rb, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/v1/support/settings", HostURL),
		strings.NewReader(string(rb)))
	req.Header.Set("Authorization", bearerToken)
	req.Header.Add("Content-Type", "application/json")
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	cr := MonitorResponse{}
	err = json.Unmarshal(body, &cr)
	if err != nil {
		return nil, err
	}

	return &cr, nil
}
