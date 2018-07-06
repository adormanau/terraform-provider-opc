package opc

import (
	"fmt"

	"github.com/hashicorp/go-oracle-terraform/client"
	"github.com/hashicorp/go-oracle-terraform/compute"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func resourceOPCSecurityApplication() *schema.Resource {
	return &schema.Resource{
		Create: resourceOPCSecurityApplicationCreate,
		Read:   resourceOPCSecurityApplicationRead,
		Delete: resourceOPCSecurityApplicationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"protocol": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateIPProtocol,
				ForceNew:     true,
			},

			"dport": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"icmptype": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(compute.Echo),
					string(compute.Reply),
					string(compute.TTL),
					string(compute.TraceRoute),
					string(compute.Unreachable),
				}, true),
				ForceNew: true,
			},

			"icmpcode": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(compute.Admin),
					string(compute.Df),
					string(compute.Host),
					string(compute.Network),
					string(compute.Port),
					string(compute.Protocol),
				}, true),
				ForceNew: true,
			},
		},
	}
}

func resourceOPCSecurityApplicationCreate(d *schema.ResourceData, meta interface{}) error {
	name := d.Get("name").(string)
	protocol := d.Get("protocol").(string)
	dport := d.Get("dport").(string)
	icmptype := d.Get("icmptype").(string)
	icmpcode := d.Get("icmpcode").(string)
	description := d.Get("description").(string)

	computeClient, err := meta.(*Client).getComputeClient()
	if err != nil {
		return err
	}
	resClient := computeClient.SecurityApplications()

	input := compute.CreateSecurityApplicationInput{
		Name:        name,
		Description: description,
		Protocol:    compute.SecurityApplicationProtocol(protocol),
		DPort:       dport,
		ICMPCode:    compute.SecurityApplicationICMPCode(icmpcode),
		ICMPType:    compute.SecurityApplicationICMPType(icmptype),
	}
	info, err := resClient.CreateSecurityApplication(&input)
	if err != nil {
		return fmt.Errorf("Error creating security application %s: %s", name, err)
	}

	d.SetId(info.Name)

	return resourceOPCSecurityApplicationRead(d, meta)
}

func resourceOPCSecurityApplicationRead(d *schema.ResourceData, meta interface{}) error {
	computeClient, err := meta.(*Client).getComputeClient()
	if err != nil {
		return err
	}
	resClient := computeClient.SecurityApplications()
	name := d.Id()

	input := compute.GetSecurityApplicationInput{
		Name: name,
	}

	result, err := resClient.GetSecurityApplication(&input)
	if err != nil {
		if client.WasNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading security application %s: %s", name, err)
	}

	if result == nil {
		d.SetId("")
		return nil
	}

	d.Set("name", result.Name)
	d.Set("protocol", result.Protocol)
	d.Set("dport", result.DPort)
	d.Set("icmptype", result.ICMPType)
	d.Set("icmpcode", result.ICMPCode)
	d.Set("description", result.Description)

	return nil
}

func resourceOPCSecurityApplicationDelete(d *schema.ResourceData, meta interface{}) error {
	computeClient, err := meta.(*Client).getComputeClient()
	if err != nil {
		return err
	}
	resClient := computeClient.SecurityApplications()
	name := d.Id()

	input := compute.DeleteSecurityApplicationInput{
		Name: name,
	}
	if err := resClient.DeleteSecurityApplication(&input); err != nil {
		return fmt.Errorf("Error deleting security application '%s': %s", name, err)
	}
	return nil
}
