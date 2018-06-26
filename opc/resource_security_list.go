package opc

import (
	"fmt"

	"github.com/hashicorp/go-oracle-terraform/client"
	"github.com/hashicorp/go-oracle-terraform/compute"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func resourceOPCSecurityList() *schema.Resource {
	return &schema.Resource{
		Create: resourceOPCSecurityListCreate,
		Read:   resourceOPCSecurityListRead,
		Update: resourceOPCSecurityListUpdate,
		Delete: resourceOPCSecurityListDelete,
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
			},

			"policy": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "deny",
				ValidateFunc: validation.StringInSlice([]string{
					string(compute.SecurityListPolicyDeny),
					string(compute.SecurityListPolicyPermit),
					string(compute.SecurityListPolicyReject),
				}, true),
				DiffSuppressFunc: suppressCaseDifferences,
			},

			"outbound_cidr_policy": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "permit",
				ValidateFunc: validation.StringInSlice([]string{
					string(compute.SecurityListPolicyDeny),
					string(compute.SecurityListPolicyPermit),
					string(compute.SecurityListPolicyReject),
				}, true),
				DiffSuppressFunc: suppressCaseDifferences,
			},

			"fqdn": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceOPCSecurityListCreate(d *schema.ResourceData, meta interface{}) error {
	name := d.Get("name").(string)
	description := d.Get("description").(string)
	policy := d.Get("policy").(string)
	outboundCIDRPolicy := d.Get("outbound_cidr_policy").(string)

	client := meta.(*Client).computeClient.SecurityLists()
	input := compute.CreateSecurityListInput{
		Name:               name,
		Description:        description,
		Policy:             compute.SecurityListPolicy(policy),
		OutboundCIDRPolicy: compute.SecurityListPolicy(outboundCIDRPolicy),
	}
	info, err := client.CreateSecurityList(&input)
	if err != nil {
		return fmt.Errorf("Error creating security list %s: %s", name, err)
	}

	d.SetId(info.Name)

	return resourceOPCSecurityListRead(d, meta)
}

func resourceOPCSecurityListUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).computeClient.SecurityLists()

	name := d.Get("name").(string)
	description := d.Get("description").(string)
	policy := d.Get("policy").(string)
	outboundCIDRPolicy := d.Get("outbound_cidr_policy").(string)

	input := compute.UpdateSecurityListInput{
		Name:               name,
		Description:        description,
		Policy:             compute.SecurityListPolicy(policy),
		OutboundCIDRPolicy: compute.SecurityListPolicy(outboundCIDRPolicy),
	}
	_, err := client.UpdateSecurityList(&input)
	if err != nil {
		return fmt.Errorf("Error updating security list %s: %s", name, err)
	}

	return resourceOPCSecurityListRead(d, meta)
}

func resourceOPCSecurityListRead(d *schema.ResourceData, meta interface{}) error {
	computeClient := meta.(*Client).computeClient.SecurityLists()

	name := d.Id()

	input := compute.GetSecurityListInput{
		Name: name,
	}

	result, err := computeClient.GetSecurityList(&input)
	if err != nil {
		// Security List does not exist
		if client.WasNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading security list %s: %s", name, err)
	}

	if result == nil {
		d.SetId("")
		return nil
	}

	d.Set("name", result.Name)
	d.Set("description", result.Description)
	d.Set("policy", string(result.Policy))
	d.Set("outbound_cidr_policy", string(result.OutboundCIDRPolicy))
	d.Set("fqdn", result.FQDN)

	return nil
}

func resourceOPCSecurityListDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).computeClient.SecurityLists()

	name := d.Id()
	input := compute.DeleteSecurityListInput{
		Name: name,
	}
	if err := client.DeleteSecurityList(&input); err != nil {
		return fmt.Errorf("Error deleting security list %s: %s", name, err)
	}

	return nil
}
