package vra

import (
	"fmt"
	"log"

	"github.com/vmware/vra-sdk-go/pkg/client/network_ip_range"
	"github.com/vmware/vra-sdk-go/pkg/models"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceExternalNetworkIPRange() *schema.Resource {
	return &schema.Resource{
		Read:   resourceExternalNetworkIPRangeRead,
		Update: resourceExternalNetworkIPRangeUpdate,
		Delete: resourceExternalNetworkIPRangeDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"end_ip_address": {
				Type:     schema.TypeString,
				Required: true,
				// Do we need to validate?
			},
			"external_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"fabric_network_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ip_version": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"IPv4", "IPv6"}, true),
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"org_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"start_ip_address": {
				Type:     schema.TypeString,
				Required: true,
				// Do we need to validate?
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags":  tagsSchema(),
			"links": linksSchema(),
		},
	}
}

func resourceExternalNetworkIPRangeRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("Reading the vra_network_profile resource with name %s", d.Get("name"))
	apiClient := m.(*Client).apiClient

	id := d.Id()
	resp, err := apiClient.NetworkIPRange.GetExternalNetworkIPRange(network_ip_range.NewGetExternalNetworkIPRangeParams().WithID(id))
	if err != nil {
		return err
	}

	networkIPRange := *resp.Payload
	//d.Set("custom_properties", networkIPRange.CustomProperties)
	d.Set("created_at", networkIPRange.CreatedAt)
	d.Set("description", networkIPRange.Description)
	d.Set("end_ip_address", networkIPRange.EndIPAddress)
	d.Set("ip_version", networkIPRange.IPVersion)
	d.Set("name", networkIPRange.Name)
	d.Set("org_id", networkIPRange.OrganizationID)
	d.Set("name", networkIPRange.Name)
	d.Set("owner", networkIPRange.Owner)
	d.Set("start_ip_address", networkIPRange.StartIPAddress)
	d.Set("updated_at", networkIPRange.UpdatedAt)

	if err := d.Set("tags", flattenTags(networkIPRange.Tags)); err != nil {
		return fmt.Errorf("error setting network ip range tags - error: %v", err)
	}

	if err := d.Set("links", flattenLinks(networkIPRange.Links)); err != nil {
		return fmt.Errorf("error setting network ip range links - error: %#v", err)
	}

	log.Printf("Finished reading the vra_network_ip_range resource with name %s", d.Get("name"))

	return nil
}

func resourceExternalNetworkIPRangeUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("Starting to Update vra_network_profile resource")
	apiClient := m.(*Client).apiClient

	id := d.Id()

	networkIPRangeSpecification := models.UpdateExternalNetworkIPRangeSpecification{
		FabricNetworkID: d.Get("fabric_network_id").(string),
	}

	log.Printf("[DEBUG] update network ip range: %#v", networkIPRangeSpecification)

	_, err := apiClient.NetworkIPRange.UpdateExternalNetworkIPRange(network_ip_range.NewUpdateExternalNetworkIPRangeParams().WithID(id).WithBody(&networkIPRangeSpecification))
	if err != nil {
		return err
	}
	log.Printf("finished Updating vra_network_profile resource")
	return resourceExternalNetworkIPRangeRead(d, m)

}

func resourceExternalNetworkIPRangeDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("Starting to delete the vra_network_ip_range resource with name %s", d.Get("name"))
	apiClient := m.(*Client).apiClient

	id := d.Id()

	networkIPRangeSpecification := models.UpdateExternalNetworkIPRangeSpecification{
		FabricNetworkID: "",
	}
	_, err := apiClient.NetworkIPRange.UpdateExternalNetworkIPRange(network_ip_range.NewUpdateExternalNetworkIPRangeParams().WithID(id).WithBody(&networkIPRangeSpecification))
	if err != nil {
		return err
	}

	d.SetId("")
	log.Printf("Finished deleting the vra_network_ip_range resource with name %s", d.Get("name"))
	return nil
}
